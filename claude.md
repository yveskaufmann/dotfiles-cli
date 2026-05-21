# Claude Notes

This file extends [agents.md](agents.md) with model-specific guidance.

## Working Style

- Keep changes minimal and scoped.
- Prefer updating existing docs over duplicating content.
- Keep provider examples aligned with actual YAML keys in `internal/config/types.go`.

## Required Cross-Checks

Before finishing documentation work:

1. Ensure README links resolve.
2. Ensure `docs/providers/index.md` lists all documented providers.
3. Ensure [agents.md](agents.md) references architecture and provider docs.
4. Ensure `.agents/` contains current task and learning notes.

## Release Safety Notes

When preparing public release docs, include explicit verification steps for:

- who can create/push release tags
- who can trigger release workflow runs
- who has admin/maintain access to the repository
