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
- [x] Create configuration guides in `/docs/config/` <!-- id: 6 -->
  - [x] `index.md`
  - [x] `apt.md`
  - [x] `ppa.md`
  - [x] `github.md`
  - [x] `binary.md`
  - [x] `snap_pipx_npm.md`
  - [x] `scripts.md`

## Completed

- [x] Create `/tasks` folder and initial backlog
