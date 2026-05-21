# Installation Guide

## Current Status

This repository is in the private-first extraction phase. The CLI already builds and runs locally, and the remaining work is to publish it to a private remote when migration is complete.

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

## V1 Path Convention

For V1, the CLI expects dotfiles content at `$HOME/.dotfiles`.
Path customization is intentionally deferred to a later version.

## After Publication

Once releases are published in `yveskaufmann/dotfiles-cli`, the intended install flow is:

1. Install with `scripts/install.sh` from this repository.
2. Run `dotfiles version` to confirm installation.
3. Run `dotfiles bootstrap` to clone or update your dotfiles repository.

## Requirements

- Go 1.24+
- Git
- A separate dotfiles configuration repository for bootstrap and linking