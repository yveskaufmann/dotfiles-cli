# Agent Baseline

This file defines the common baseline for automated contributors in this repository.

## Project Basics

- Project: standalone bootstrap CLI for dotfiles workflows
- Goal: keep executable behavior here and personal shell/config in a separate dotfiles repository
- Core commands: `bootstrap`, `install`, `link`

Keep this repository focused on:

- command behavior (`bootstrap`, `install`, `link`)
- configuration parsing and validation
- provider orchestration and reporting

## Read First

- [README.md](README.md)
- [architecture.md](architecture.md)
- [docs/providers/index.md](docs/providers/index.md)
- [.agents/tasks.md](.agents/tasks.md)
- [.agents/learnings.md](.agents/learnings.md)
- [claude.md](claude.md)

## Repository Boundaries

1. Keep command behavior and provider execution logic in this repository.
2. Keep user-specific dotfiles content in the separate dotfiles repository.
3. Prefer documentation updates alongside behavior changes.

## Working Protocol

1. Read this file before implementing changes.
2. Confirm architecture impact in [architecture.md](architecture.md).
3. Update provider docs when config behavior changes.
4. Record all todos in [.agents/tasks.md](.agents/tasks.md). Update task status (todo -> in-progress -> done) as you work.
5. Record technical insights, patterns, and discoveries in [.agents/learnings.md](.agents/learnings.md) for continuity across sessions.

## Provider Changes

When modifying provider interfaces, adding new providers, or removing providers:

1. Update the schema in `internal/config/types.go` to reflect provider struct changes.
2. Update the provider registry in the provider implementation to register/unregister providers.
3. Update documentation:
	- If adding/removing providers: update [docs/providers/index.md](docs/providers/index.md) provider list.
	- If changing provider config format: update the corresponding `docs/providers/*.md` page.
	- If adding a new provider: create a new `docs/providers/<provider>.md` page.
4. Validate consistency: run `make build` and tests to ensure schema and implementation remain synchronized.

## Documentation Expectations

When behavior changes, update:

1. [README.md](README.md) for human usage impact.
2. [docs/providers/index.md](docs/providers/index.md) and related provider pages for config changes.
3. [architecture.md](architecture.md) for workflow/architecture impact.

## Validation Sequence

1. Build (`make build`)
2. Run command help/version smoke checks (`dotfiles version`, `dotfiles <cmd> --help`)
3. Run tests
4. Verify release workflow assumptions

## Repository Layout

```text
cmd/            CLI entry point
internal/       bootstrap, install, link, config, providers, and utilities
docs/           human-facing documentation
.agents/        planning, tasks, and learnings for agent workflows
```

## Model-Specific Notes

Claude-specific execution notes are in [claude.md](claude.md).
