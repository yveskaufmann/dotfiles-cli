# Dotfiles CLI Agent

This repository holds the executable side of the dotfiles bootstrapper.

## Purpose

The CLI handles orchestration only:

- repository bootstrap
- tool installation from YAML definitions
- symlink creation
- logging and provider execution

## Operating Rules

1. Keep executable behavior in this repository.
2. Keep personal configuration and shell content in the separate dotfiles repository.
3. Validate changes with build and execution checks before reshaping repository boundaries.
4. Use the task backlog to track remaining extraction work.

## Architecture Summary

The code is organized around:

- `cmd/dotfiles` for the binary entry point
- `internal/cli` for commands and flags
- `internal/engine` for bootstrap/install/link orchestration
- `internal/config` for YAML loading and validation
- `internal/provider` for package manager integrations

## Validation Priority

When working in this repository, validate in this order:

1. build
2. run
3. tests
4. only then adjust repository structure or distribution