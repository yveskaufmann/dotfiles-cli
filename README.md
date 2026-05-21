# dotfiles-cli

`dotfiles-cli` is a Go-based bootstrap tool for installing developer tooling and linking dotfiles from a separate configuration repository.

## V1 Scope

- The local dotfiles workspace is fixed to `$HOME/.dotfiles`.
- Repository source is generic and can be provided with `--repository`.
- Configurable local paths are intentionally deferred to a later version.

## What It Does

- clones or updates a dotfiles repository
- installs tools from declarative `init/*.yaml` files
- creates symlinks from `link/` into your home directory
- keeps the bootstrap logic separate from your personal config

## Install

The installer script lives in this repository at `scripts/install.sh` and
downloads the matching `dotfiles` binary from GitHub Releases.

Installer command:

```bash
curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh
```

## Build

```bash
make build
./bin/dotfiles version
```

## Usage

```bash
dotfiles bootstrap
dotfiles bootstrap --repository git@github.com:user/dotfiles.git
dotfiles install --profile default
dotfiles link --dry-run
```

## Documentation

- [Installation](docs/INSTALL.md)
- [Architecture](architecture.md)
- [Task Backlog](tasks/backlog.md)

## Repository Layout

```text
cmd/        CLI entry point
internal/    bootstrap, install, link, providers, and utilities
docs/        CLI-focused documentation
tasks/       tracking and backlog notes
.devcontainer/  development container setup
```
