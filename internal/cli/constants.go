package cli

import (
	"os"
	"path"
	"time"
)

var (
	DOTFILES_PATH     = "$HOME/.dotfiles"
	CACHE_PATH        = path.Join(DOTFILES_PATH, ".caches")
	BACKUP_PATH       = path.Join(DOTFILES_PATH, "backups", time.Now().Format("2006_01_02-15_04_05"))
	INIT_SCRIPTS_PATH = path.Join(DOTFILES_PATH, "init")
)

func init() {
	DOTFILES_PATH = os.ExpandEnv(DOTFILES_PATH)
	BACKUP_PATH = os.ExpandEnv(BACKUP_PATH)
	CACHE_PATH = os.ExpandEnv(CACHE_PATH)
	INIT_SCRIPTS_PATH = os.ExpandEnv(INIT_SCRIPTS_PATH)
}
