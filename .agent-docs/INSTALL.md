# Installation Guide

## Quick Install

Install the latest `dotfiles` binary from the private
`yveskaufmann/dotfiles-cli` repository:

```bash
GITHUB_TOKEN="$(gh auth token)"
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
   -H "Accept: application/vnd.github.raw" \
   -fsSL "https://api.github.com/repos/yveskaufmann/dotfiles-cli/contents/scripts/install.sh?ref=main" \
   | sh -s -- --github-token "$GITHUB_TOKEN"
```

## Installation Options

### Default Installation (User-local)

Installs to `~/.local/bin/dotfiles` (no sudo required):

```bash
GITHUB_TOKEN="$(gh auth token)"
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
   -H "Accept: application/vnd.github.raw" \
   -fsSL "https://api.github.com/repos/yveskaufmann/dotfiles-cli/contents/scripts/install.sh?ref=main" \
   | sh -s -- --github-token "$GITHUB_TOKEN"
```

### Force Home Directory Installation

Explicitly install to `~/.local/bin`:

```bash
GITHUB_TOKEN="$(gh auth token)"
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
   -H "Accept: application/vnd.github.raw" \
   -fsSL "https://api.github.com/repos/yveskaufmann/dotfiles-cli/contents/scripts/install.sh?ref=main" \
   | sh -s -- --github-token "$GITHUB_TOKEN" --home
```

### System-wide Installation

Install to `/usr/local/bin` (may require sudo):

```bash
GITHUB_TOKEN="$(gh auth token)"
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
   -H "Accept: application/vnd.github.raw" \
   -fsSL "https://api.github.com/repos/yveskaufmann/dotfiles-cli/contents/scripts/install.sh?ref=main" \
   | sudo sh -s -- --github-token "$GITHUB_TOKEN" --system
```

### Custom Directory Installation

Install to a specific directory:

```bash
GITHUB_TOKEN="$(gh auth token)"
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
   -H "Accept: application/vnd.github.raw" \
   -fsSL "https://api.github.com/repos/yveskaufmann/dotfiles-cli/contents/scripts/install.sh?ref=main" \
   | sh -s -- --github-token "$GITHUB_TOKEN" --dir ~/bin
```

## What the Script Does

1. **Auto-detects** your OS (Linux, macOS, Windows) and architecture (amd64, arm64)
2. **Fetches** the latest release version from `yveskaufmann/dotfiles-cli`
3. **Downloads** the appropriate binary for your system
4. **Installs** to the specified location (default: `~/.local/bin`)
5. **Verifies** the installation by running `dotfiles version`
6. **Provides** PATH instructions if needed

## After Installation

Once installed, bootstrap your dotfiles:

```bash
# Bootstrap everything (clone repository, install tools, create symlinks)
dotfiles bootstrap

# Or step by step:
dotfiles bootstrap --repository git@github.com:yourusername/dotfiles.git
dotfiles install --profile work
dotfiles link
```

## Supported Platforms

| OS      | Architecture | Status      |
|---------|-------------|-------------|
| Linux   | amd64       | ✅ Supported |
| Linux   | arm64       | ✅ Supported |
| macOS   | amd64       | ✅ Supported |
| macOS   | arm64       | ✅ Supported |
| Windows | amd64       | ✅ Supported |
| Windows | arm64       | ✅ Supported |

## Requirements

- `curl` or `wget` (for downloading)
- `tar` (for extracting archives on Linux/macOS)
- `unzip` (for extracting archives on Windows)

## Troubleshooting

### Binary not in PATH

If the installer shows a PATH warning, add the install directory to your PATH:

**For bash:**
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**For zsh:**
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Permission denied

If installing to `/usr/local/bin` fails with permission errors, either:
- Use `sudo`: `curl ... | sudo sh -s -- --system`
- Or install to home directory: `curl ... | sh -s -- --home`

### Download fails

If the download fails:
1. Check your internet connection
2. Verify your GitHub token can access the private repository
3. Verify the repository has releases: https://github.com/yveskaufmann/dotfiles-cli/releases
4. Try using `wget` instead of `curl` (the script auto-detects both)

### Unsupported OS/Architecture

The installer supports Linux, macOS, and Windows on amd64 and arm64. If you're using a different platform, you'll need to:
1. Download the binary manually from [releases](https://github.com/yveskaufmann/dotfiles-cli/releases)
2. Extract and move it to a directory in your PATH
3. Make it executable: `chmod +x dotfiles`

## Manual Installation

If you prefer not to use the install script:

1. Download the appropriate binary from [releases](https://github.com/yveskaufmann/dotfiles-cli/releases)
2. Extract the archive:
   ```bash
   tar -xzf dotfiles-cli_v1.0.0_linux_amd64.tar.gz
   ```
3. Move to a directory in your PATH:
   ```bash
   mv dotfiles ~/.local/bin/
   chmod +x ~/.local/bin/dotfiles
   ```
4. Verify installation:
   ```bash
   dotfiles version
   ```

## Uninstallation

To remove the dotfiles binary:

```bash
# If installed to ~/.local/bin
rm ~/.local/bin/dotfiles

# If installed to /usr/local/bin
sudo rm /usr/local/bin/dotfiles
```

To also remove configuration and dotfiles:

```bash
rm -rf ~/.dotfiles
rm -rf ~/.config/.dotfiles
```
