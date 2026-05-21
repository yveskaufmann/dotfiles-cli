# Dotfiles CLI Agent Guide

This page is the shared baseline for agents working in this repository.

## Project Basics

- Project: standalone bootstrap CLI for dotfiles workflows
- Goal: keep executable behavior here and personal shell/config in a separate dotfiles repository
- Core commands: `bootstrap`, `install`, `link`

## Source of Truth

- Human-facing project summary: [README.md](../README.md)
- Architecture overview: [architecture.md](../architecture.md)
- Provider configuration docs: [docs/providers/index.md](providers/index.md)
- Agent coordination docs: [agents.md](../agents.md)
- Claude-specific operating notes: [claude.md](../claude.md)
- Active planning and learnings: [.agents/tasks.md](../.agents/tasks.md)

## Repository Boundaries

1. Keep command behavior and provider execution logic in this repository.
2. Keep user-specific dotfiles content in the separate dotfiles repository.
3. Prefer documentation updates alongside behavior changes.

## Working Protocol

1. Read [agents.md](../agents.md) before implementing changes.
2. Confirm architecture impact in [architecture.md](../architecture.md).
3. Update provider docs when config behavior changes.
4. Track progress and notes in `.agents/` files.

## Validation Order

1. Build (`make build`)
2. Command run checks (`dotfiles version`, `dotfiles <cmd> --help`)
3. Tests
4. Release and distribution workflows