package executor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
	"yv35.com/dotfiles/internal/theme"
	"yv35.com/dotfiles/internal/util/pathutil"
	"yv35.com/dotfiles/internal/util/stringutils"
)

type FileLinkerOptions struct {
	DryRun                    bool
	DefaultConflictResolution LinkConflictResolution
}

type FileLinker struct {
	sourceDir        string
	targetDir        string
	ignoredFileNames []string
	// When defined, all conflicts are resolved in the same way
	defaultConflictResolution LinkConflictResolution
	DryRun                    bool
}

const DOTFILES_DEEP_LINK = ".dotfiles-deep-link"
const GIT_DIR = ".git"
const TEMP_LINK_DIR = ".dotfiles-tmp-link"

type LinkOperationConflict int

var (
	ConflictNone          LinkOperationConflict = 0x0
	ConflictSymlinkExists LinkOperationConflict = 0x2
	ConflictFileExists    LinkOperationConflict = 0x4
)

type LinkConflictResolution int

var (
	ResolutionNone    LinkConflictResolution = 0x0
	ResolutionSkip    LinkConflictResolution = 0x2
	ResolutionReplace LinkConflictResolution = 0x4
	ResolutionBackup  LinkConflictResolution = 0x8
)

type linkOperation struct {
	target             string
	location           string
	conflict           LinkOperationConflict
	conflictResolution LinkConflictResolution
}

type operationCounters struct {
	created       int
	skipped       int
	alreadyLinked int
	backedUp      int
	replaced      int
}

func NewFileLinker(sourceDir, targetDir string, opts FileLinkerOptions) *FileLinker {
	return &FileLinker{
		sourceDir:                 sourceDir,
		targetDir:                 targetDir,
		ignoredFileNames:          []string{".DS_Store", GIT_DIR, DOTFILES_DEEP_LINK},
		defaultConflictResolution: opts.DefaultConflictResolution,
		DryRun:                    opts.DryRun,
	}
}

func (l *FileLinker) Execute() error {

	deepLinkDirs := []string{}
	isDirectoryDeepLinked := func(path string) bool {
		for _, dir := range deepLinkDirs {
			if strings.HasPrefix(path, dir) {
				return true
			}
		}
		return false
	}

	linkOperations := make([]linkOperation, 0)
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		absPath := filepath.Join(l.sourceDir, path)

		if stringutils.ListContains(l.ignoredFileNames, d.Name()) {
			return nil
		}

		if d.IsDir() && path == "." {
			return nil
		}

		if d.IsDir() {
			if path == "." {
				return nil
			}

			deeplinkSignaturePath := filepath.Join(absPath, DOTFILES_DEEP_LINK)
			if _, err := os.Stat(deeplinkSignaturePath); err == nil {
				deepLinkDirs = append(deepLinkDirs, path)
				return nil
			}

			if isDirectoryDeepLinked(path) {
				return nil
			}

			return fs.SkipDir
		}

		linkOperation := linkOperation{
			target:   absPath,
			location: filepath.Join(l.targetDir, stringutils.StripPrefixDirPath(path, l.sourceDir)),
		}
		linkOperations = append(linkOperations, linkOperation)

		return nil
	}

	if err := fs.WalkDir(os.DirFS(l.sourceDir), ".", walkDirFunc); err != nil {
		return fmt.Errorf("failed to scan link directory %s: %w", l.sourceDir, err)
	}

	counters := &operationCounters{}
	for _, op := range linkOperations {
		if err := l.doLinkCreationProcess(&op, counters); err != nil {
			return fmt.Errorf("failed to create link from %s to %s: %w", op.target, op.location, err)
		}
	}

	l.printSummary(counters)
	return nil
}

func (l *FileLinker) doLinkCreationProcess(linkOp *linkOperation, counters *operationCounters) error {

	err := linkOp.CheckConflicts()
	if err != nil {
		return fmt.Errorf("failed to check for conflicts for link %s: %w", pathutil.MinimizePath(linkOp.location), err)
	}

	// Check if already skipped (e.g., symlink already pointing to correct target)
	if linkOp.conflictResolution == ResolutionSkip && linkOp.conflict == ConflictNone {
		counters.alreadyLinked++
		return nil
	}

	if linkOp.conflict != ConflictNone {
		if err := l.resolveConflict(linkOp); err != nil {
			return err
		}

		if linkOp.conflictResolution == ResolutionSkip {
			counters.skipped++
			return nil
		}

		if err := l.applyResolution(linkOp, counters); err != nil {
			return err
		}
	}

	if !l.DryRun {
		err = os.MkdirAll(filepath.Dir(linkOp.location), 0700)
		if err != nil {
			return fmt.Errorf("failed to create parent directory for link %s: %w", pathutil.MinimizePath(linkOp.location), err)
		}
	}

	if l.DryRun {
		fmt.Printf("%s[DRY RUN] ✅ %s → %s%s\n", theme.Colorize(theme.ColorGreen), pathutil.MinimizePath(linkOp.location), pathutil.MinimizePath(linkOp.target), theme.Colorize(theme.ColorReset))
	} else {
		err = os.Symlink(linkOp.target, linkOp.location)
		if err != nil {
			return fmt.Errorf("failed to create symlink from %s → %s: %w", pathutil.MinimizePath(linkOp.location), pathutil.MinimizePath(linkOp.target), err)
		}
		fmt.Printf("%s✅ %s → %s%s\n", theme.Colorize(theme.ColorGreen), pathutil.MinimizePath(linkOp.location), pathutil.MinimizePath(linkOp.target), theme.Colorize(theme.ColorReset))
	}

	counters.created++

	return nil
}

func (l *FileLinker) resolveConflict(linkOp *linkOperation) error {
	if l.defaultConflictResolution != ResolutionNone {
		linkOp.conflictResolution = l.defaultConflictResolution
		return nil
	}

	message := l.formatConflictMessage(linkOp.conflict, linkOp.location, linkOp.target)
	action := promptForNextAction(message)

	resolutions := map[string]LinkConflictResolution{
		"o": ResolutionReplace, "O": ResolutionReplace,
		"b": ResolutionBackup, "B": ResolutionBackup,
		"s": ResolutionSkip, "S": ResolutionSkip,
	}

	if res, ok := resolutions[action]; ok {
		linkOp.conflictResolution = res
		if action == strings.ToUpper(action) {
			l.defaultConflictResolution = res
		}
	} else {
		linkOp.conflictResolution = ResolutionSkip
	}

	return nil
}

func (l *FileLinker) formatConflictMessage(conflict LinkOperationConflict, location, target string) string {
	minLoc := pathutil.MinimizePath(location)
	minTarget := pathutil.MinimizePath(target)
	switch conflict {
	case ConflictSymlinkExists:
		existingTarget, _ := os.Readlink(location)
		return fmt.Sprintf("⚠️  Symlink already exists at %s\n  points to     → %s\n  want to link to → %s", minLoc, pathutil.MinimizePath(existingTarget), minTarget)
	case ConflictFileExists:
		return fmt.Sprintf("⚠️  File/directory already exists at %s\n  want to create symlink to → %s", minLoc, minTarget)
	default:
		return fmt.Sprintf("⚠️  Conflict at %s", minLoc)
	}
}

func (l *FileLinker) applyResolution(linkOp *linkOperation, counters *operationCounters) error {
	switch linkOp.conflictResolution {
	case ResolutionBackup:
		backupLocation := linkOp.location + ".backup"
		if l.DryRun {
			fmt.Printf("%s[DRY RUN] 💾 %s → %s%s\n", theme.Colorize(theme.ColorCyan), pathutil.MinimizePath(linkOp.location), pathutil.MinimizePath(backupLocation), theme.Colorize(theme.ColorReset))
		} else {
			if err := os.Rename(linkOp.location, backupLocation); err != nil {
				return fmt.Errorf("failed to backup existing file at %s to %s: %w", linkOp.location, backupLocation, err)
			}
			fmt.Printf("%s💾 %s → %s%s\n", theme.Colorize(theme.ColorCyan), pathutil.MinimizePath(linkOp.location), pathutil.MinimizePath(backupLocation), theme.Colorize(theme.ColorReset))
		}
		counters.backedUp++
	case ResolutionReplace:
		if l.DryRun {
			fmt.Printf("%s[DRY RUN] 🗑️  %s%s\n", theme.Colorize(theme.ColorYellow), pathutil.MinimizePath(linkOp.location), theme.Colorize(theme.ColorReset))
		} else {
			if err := os.RemoveAll(linkOp.location); err != nil {
				return fmt.Errorf("failed to remove existing file at %s: %w", linkOp.location, err)
			}
			fmt.Printf("%s🗑️  %s%s\n", theme.Colorize(theme.ColorYellow), pathutil.MinimizePath(linkOp.location), theme.Colorize(theme.ColorReset))
		}
		counters.replaced++
	}
	return nil
}

func (lp *linkOperation) CheckConflicts() error {
	stat, err := os.Lstat(lp.location)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lp.conflict = ConflictNone
			return nil
		}
		return fmt.Errorf("failed to lstat link location %s: %w", lp.location, err)
	}

	if stat.Mode()&os.ModeSymlink != 0 {
		return lp.checkSymlinkConflict()
	}

	lp.conflict = ConflictFileExists
	return nil
}

func (lp *linkOperation) checkSymlinkConflict() error {
	linkTarget, err := os.Readlink(lp.location)
	if err != nil {
		return fmt.Errorf("failed to readlink existing symlink at %s: %w", lp.location, err)
	}

	if linkTarget == lp.target {
		// Symlink already points to correct target - no conflict, skip silently
		lp.conflict = ConflictNone
		lp.conflictResolution = ResolutionSkip
	} else {
		// Symlink points to different target - conflict
		lp.conflict = ConflictSymlinkExists
	}

	return nil
}

func promptForNextAction(message string) string {
	fmt.Printf("%s%s%s\n", theme.Colorize(theme.ColorYellow), message, theme.Colorize(theme.ColorReset))
	fmt.Printf("Action? [%ss%skip [%sS%skip all | %so%sverwrite [%sO%sverwrite all | %sb%sackup [%sB%sackup all] (default: s) ",
		theme.Colorize(theme.ColorGreen), theme.Colorize(theme.ColorReset),
		theme.Colorize(theme.ColorGreen), theme.Colorize(theme.ColorReset),
		theme.Colorize(theme.ColorYellow), theme.Colorize(theme.ColorReset),
		theme.Colorize(theme.ColorYellow), theme.Colorize(theme.ColorReset),
		theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset),
		theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset))

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "s"
	}

	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("\n\n")
	}()

	// Read a single byte
	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return "s"
	}

	return string(buf[0])

}

func (l *FileLinker) printSummary(counters *operationCounters) {
	prefix := ""
	if l.DryRun {
		prefix = "[DRY RUN] "
	}

	fmt.Printf("\n%s─────────────────────────────────────────%s\n", theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset))
	fmt.Printf("%s%sSummary:%s\n", prefix, theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset))
	fmt.Printf("  %s✅ %d created%s\n", theme.Colorize(theme.ColorGreen), counters.created, theme.Colorize(theme.ColorReset))
	if counters.alreadyLinked > 0 {
		fmt.Printf("  ✓ %d already linked\n", counters.alreadyLinked)
	}
	if counters.skipped > 0 {
		fmt.Printf("  ↪️  %d skipped\n", counters.skipped)
	}
	if counters.backedUp > 0 {
		fmt.Printf("  %s💾 %d backed up%s\n", theme.Colorize(theme.ColorCyan), counters.backedUp, theme.Colorize(theme.ColorReset))
	}
	if counters.replaced > 0 {
		fmt.Printf("  %s🗑️  %d replaced%s\n", theme.Colorize(theme.ColorYellow), counters.replaced, theme.Colorize(theme.ColorReset))
	}
	fmt.Printf("%s─────────────────────────────────────────%s\n", theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset))
}
