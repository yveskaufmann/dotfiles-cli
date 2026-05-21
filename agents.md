# Agent Baseline

This file defines the common baseline for automated contributors in this repository.

## Project Goal

`dotfiles-cli` is the executable bootstrap layer for dotfiles workflows.

Keep this repository focused on:

- command behavior (`bootstrap`, `install`, `link`)
- configuration parsing and validation
- provider orchestration and reporting

Keep user-specific shell and dotfile content in a separate dotfiles repository.

## Read First

- [docs/agent.md](docs/agent.md)
- [architecture.md](architecture.md)
- [docs/providers/index.md](docs/providers/index.md)
- [.agents/tasks.md](.agents/tasks.md)
- [.agents/learnings.md](.agents/learnings.md)

## Documentation Expectations

When behavior changes, update:

1. [README.md](README.md) for human usage impact
2. [docs/providers](docs/providers/index.md) pages for config changes
3. [docs/agent.md](docs/agent.md) for workflow impact

## Validation Sequence

1. Build
2. Run command help/version smoke checks
3. Run tests
4. Verify release workflow assumptions

## Model-Specific Notes

Claude-specific execution notes are in [claude.md](claude.md).
