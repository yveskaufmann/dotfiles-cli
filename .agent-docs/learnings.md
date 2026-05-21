# Technical Learnings

This document tracks key technical decisions and lessons learned during the development of the dotfile bootstrapper.

## 1. Declarative vs Imperative

**Problem**: The original system used Bash scripts (`init/*.sh`) that were difficult to audit, lacked consistent idempotency, and were hard to scale across profiles.
**Solution**: Migrated to a Go-based system using YAML (`init/*.yaml`).
**Learning**: Declarative configuration makes it easier to track what *should* be installed, while Go provides a robust standard library for cross-platform file operations and process management.

## 2. Shared interactive TTY

**Problem**: Running `sudo apt install` via `exec.Command` in Go often hangs if it requires a password but has no TTY.
**Learning**: By setting `cmd.Stdin = os.Stdin`, `cmd.Stdout = os.Stdout`, and `cmd.Stderr = os.Stderr`, the sub-process can directly leverage the parent's TTY for password prompts and progress bars.

## 3. PPA Modernization (DEB822)

**Problem**: Traditional `add-apt-repository` creates `.list` files which are being deprecated in favor of the more structured DEB822 `.sources` format.
**Learning**: We implemented a manual PPA manager that:

- Downloads GPG keys directly via GPG keyservers.
- Writes `.sources` files with explicit architecture (`Signed-By`).
- Uses `:arch` placeholders to support dynamic environment detection (`amd64` vs `arm64`).

## 4. Multi-Binary Archive Extraction

**Problem**: Tools like `crane` release several binaries (`crane`, `gcrane`, `krane`) in a single archive. Traditionally, this would trigger 3 separate downloads.
**Learning**: We implemented a `binaries` list property. If present, the system:

1. Downloads the archive once.
2. Extracts it to a temporary location.
3. Recursively searches for each binary in the extraction tree.
4. Installs them individually.

## 5. Idempotency Checks

**Learning**: Relying purely on file existence is dangerous for package managers.

- **APT**: We use `dpkg-query -W` for reliable package status.
- **Pipx/Snap**: We wrap their respective list commands.
- **Binaries**: We use `type <name>` to check the system PATH, which handles tools installed in `/opt`, `/usr/bin`, or `~/bin` regardless of the specific installation method.

## 6. Provider Setup/TearDown Pattern

**Problem**: Some providers need one-time initialization before processing any installations (e.g., NVM needs to be installed before installing Node versions).
**Solution**: Implemented optional `Setupable` and `TearDownable` interfaces.
**Learning**: 
- The executor calls `Setup()` on all providers before any `Install()` calls
- This allows providers to perform expensive one-time operations (like installing NVM)
- Keeps the `Install()` method focused on the actual artifact installation
- Improves code organization and follows single responsibility principle

## 7. Version Resolution and Synthetic Versions

**Problem**: Users want to specify "latest LTS" without hard-coding version numbers that become outdated.
**Solution**: Implemented version resolution in the NVM provider with `lts/latest` synthetic version.
**Learning**:
- Synthetic versions (like `lts/latest`) should follow existing naming conventions (NVM uses `lts/codename`)
- Resolution must happen before deduplication to ensure proper version matching
- Error handling is critical - if resolution fails, installation should stop immediately
- Version resolution requires the underlying tool (NVM) to be installed first, hence the need for `Setup()`

## 8. Shell Command Execution Context

**Problem**: Running shell commands that source scripts (like NVM) requires proper shell initialization.
**Solution**: Use `bash -lc` (login shell) for commands that need environment setup.
**Learning**:
- Different shell contexts produce different behaviors:
  - `bash -c`: Non-interactive, non-login shell
  - `bash -lc`: Login shell, sources profile files
- Tools like NVM that modify the shell environment need login shell context
- Exit code 127 typically means "command not found" - often indicates the command isn't in the expected shell environment

## 9. Provider Dependency Management

**Problem**: Version managers (SdkMan, NVM) require system tools (zip, unzip, curl) to be installed first.
**Solution**: Implemented provider priority system with automatic dependency installation.
**Learning**:
- Provider execution order matters - system package managers must run before version managers
- Priority-based execution (lower = first): apt/brew (10) → sdkman/nvm (50) → apps (100)
- Registry sorts providers by priority automatically
- Version managers should auto-install their own system dependencies rather than requiring manual config
- OS-specific dependency installation:
  - Linux: Use `apt-get` for zip/unzip
  - macOS: Usually preinstalled, use `brew` as fallback
  - Check if dependencies exist before attempting installation

## 10. Flexible Configuration Parsing

**Problem**: Users want both simple and complex configuration options without separate config structures.
**Solution**: Implemented `UnmarshalYAML` for flexible string/object parsing in SdkMan provider.
**Learning**:
- YAML supports both scalar and object nodes for the same field
- `UnmarshalYAML` allows parsing "candidate:version" strings into structured data
- Provides ergonomic simple syntax while maintaining power-user flexibility
- Pattern: Simple for single versions, object format for multiple versions or advanced options
- Example:
  ```yaml
  sdkman:
    - "java:17"              # Simple string
    - candidate: "gradle"    # Object with options
      version: "8.5"
      versions: ["7.6"]
  ```

## 11. Runtime OS Detection and Cross-Platform Support

**Problem**: Different operating systems require different installation methods for system dependencies.
**Solution**: Use `runtime.GOOS` for OS detection in provider code.
**Learning**:
- Go's `runtime.GOOS` provides reliable OS detection ("linux", "darwin", "windows")
- Cross-platform providers should handle OS-specific logic internally
- Use switch statements on `runtime.GOOS` for clean OS-specific behavior
- Fail gracefully on unsupported platforms with clear error messages
- Don't assume tools are preinstalled - verify their presence first
