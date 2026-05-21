# Installation Guide

## Quick Install

Install the latest `dotfiles` binary via the repository installer:

```bash
curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh
```

The installer script resolves the latest release, downloads the matching binary
for your OS and architecture, and installs it locally.

## Build From Source

```bash
make build
./bin/dotfiles version
```

## Running the CLI

```bash
dotfiles bootstrap
dotfiles bootstrap --repository git@github.com:user/dotfiles.git
dotfiles install --profile default
dotfiles link --dry-run
```

## Uninstall

`dotfiles` is a single binary. Remove it with:

```bash
rm "$(where dotfiles)"
```

## V1 Path Convention

For V1, the CLI expects dotfiles content at `$HOME/.dotfiles`.
Path customization is intentionally deferred to a later version.

## Requirements

- Go 1.24+
- Git
- A separate dotfiles configuration repository for bootstrap and linking