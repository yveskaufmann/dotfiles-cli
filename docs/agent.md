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
4. Record all todos in [`.agents/tasks.md`](../.agents/tasks.md). Update task status (todo → in-progress → done) as you work.
5. Record technical insights, patterns, and discoveries in [`.agents/learnings.md`](../.agents/learnings.md) for continuity across sessions.

## Provider Changes

When modifying provider interfaces, adding new providers, or removing providers:

1. **Update the schema** in `internal/config/types.go` to reflect provider struct changes
2. **Update the provider registry** in the provider implementation to register/unregister providers
3. **Update documentation**:
   - If adding/removing providers: update [docs/providers/index.md](providers/index.md) provider list
   - If changing provider config format: update the corresponding [docs/providers/*.md](providers/) page
   - If adding a new provider: create new `docs/providers/<provider>.md` page
4. **Validate consistency**: Run `make build` and tests to ensure schema and implementation remain synchronized

## Validation Order

1. Build (`make build`)
2. Command run checks (`dotfiles version`, `dotfiles <cmd> --help`)
3. Tests
4. Release and distribution workflows

## Repository Layout

```text
cmd/            CLI entry point
internal/       bootstrap, install, link, config, providers, and utilities
docs/           human-facing documentation
.agents/        planning, tasks, and learnings for agent workflows
tasks/          backlog notes
```