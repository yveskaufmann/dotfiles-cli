# Configuration Reference

The CLI loads YAML files from the configured `init/` directory and validates them against the schema in `/.schema/config.schema.json`.

## Scope

This repository owns the executable logic and schema, while the user-specific dotfiles repository owns the actual configuration files.

## What Lives Here

- schema definitions
- loader and provider behavior
- execution rules for install groups

## What Lives In The Dotfiles Repo

- concrete `init/*.yaml` configuration files
- shell dotfiles and `link/` content