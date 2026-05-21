# dotfiles-cli

`dotfiles-cli` is a Go-based bootstrap tool for installing developer tooling and linking dotfiles from a separate configuration repository.

## What It Does

- clones or updates a dotfiles repository
- installs tools from declarative `init/*.yaml` files
- creates symlinks from `link/` into your home directory
- keeps the bootstrap logic separate from your personal config

## Build

```bash
make build
./bin/dotfiles version
```

## Usage

```bash
dotfiles bootstrap
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
