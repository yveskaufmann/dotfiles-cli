# Product Backlog

## Dotfiles-CLI Migration (Current)

- [x] Extract tool to private `yveskaufmann/dotfiles-cli` repository
- [x] Validate private release pipeline and installer smoke test
- [x] Add README quick-install instructions for external bootstrap CLI
- [x] Align `.dotfiles/docs/INSTALL.md` with external `dotfiles-cli` install flow
- [x] Remove tool/runtime files from `.dotfiles` (final cleanup step)
- [x] Commit and push final config-only shape of `.dotfiles`

## Governance & Infrastructure

- [x] Update `agent.md` with task-checking protocols <!-- id: 1 -->
- [x] Add `$schema` field to `internal/init/setup-tools/types.go` <!-- id: 2 -->
- [x] Create `/.schema/config.schema.json` with detailed field descriptions <!-- id: 3 -->
- [x] Specify `$schema` in all `init/*.yaml` files <!-- id: 7 -->
- [x] Create isolated Docker test environment <!-- id: 8 -->

## Documentation

- [x] Create `architecture.md` with Mermaid diagrams <!-- id: 4 -->
- [x] Create `learnings.md` <!-- id: 5 -->
- [x] Create configuration guides in `/docs/providers/` <!-- id: 6 -->
  - [x] `index.md`
  - [x] `apt.md`
  - [x] `ppa.md`
  - [x] `github.md`
  - [x] `binary.md`
  - [x] `snap_pipx_npm.md`
  - [x] `scripts.md`

## Completed

- [x] Create `/tasks` folder and initial backlog

## Public Readiness (Current)

- [x] Consolidate planning docs under `.agents/`
- [x] Add `agents.md` baseline for all agents
- [x] Add `claude.md` linked from `agents.md`
- [x] Rewrite `architecture.md` for humans and agents
- [x] Polish human-facing `README.md`
- [x] Split provider docs into one page per provider
- [x] Verify release publisher restrictions (only authorized maintainers can publish)

Verification notes (2026-05-21):
- Direct collaborator permissions show only `yveskaufmann` with admin/maintain/push.
- Release workflow is tag-triggered (`v*.*.*`), so publish requires tag push permission.
- Actions permissions are enabled with `allowed_actions: all`.
- Rulesets/branch protection API endpoints returned 403 due current repo plan limits.
