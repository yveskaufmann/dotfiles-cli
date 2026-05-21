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
dotfiles install --profile default
dotfiles link --dry-run
```

## After Publication

Once the private repository is created and releases are published, the intended install flow is:

1. Download the binary from the release assets.
2. Verify the checksum.
3. Run `dotfiles version` to confirm installation.

## Requirements

- Go 1.24+
- Git
- A separate dotfiles configuration repository for bootstrap and linking