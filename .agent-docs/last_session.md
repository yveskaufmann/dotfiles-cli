# Session Handoff: dotfiles-cli Extraction and Migration

## Objectives

- Extract the Go bootstrap tool into a dedicated private repository (`yveskaufmann/dotfiles-cli`)
- Keep `.dotfiles` as configuration/content repository
- V1 scope decision: keep local runtime convention fixed to `$HOME/.dotfiles`

## Completed Todos

### dotfiles-cli repository

- [x] Created local extracted repository with preserved history
- [x] Added missing repository scaffolding (`.devcontainer`, `docs`, `tasks`, `architecture.md`)
- [x] Created first extraction commit
- [x] Created private GitHub repository and pushed main branch
- [x] Renamed module/import identity to `yv35.com/dotfiles-cli`
- [x] Updated release/install identity from `.dotfiles` to `dotfiles-cli`
- [x] Updated docs for V1 scope (fixed path, generic `--repository`)
- [x] Ran private release smoke test via tag `v0.2.0-rc.1`
- [x] Verified release artifacts for linux/darwin/windows amd64/arm64 + checksums
- [x] Ran installer smoke test in isolated target directory
- [x] Fixed installer for prerelease-only repos (`releases/latest` fallback)
- [x] Fixed installer archive name mismatch (`dotfiles-cli_*` assets)
- [x] Verified installed binary runs from smoke-test path

### .dotfiles repository

- [x] Updated `README.md` Quick Install section with private `dotfiles-cli` installation instructions

## Open Todos

### .dotfiles repository (pre-cleanup docs)

- [x] Align `docs/INSTALL.md` to external `dotfiles-cli` flow (currently still references `.dotfiles` release installer)
- [x] Remove stale README sections that still imply CLI source lives in `.dotfiles` (project structure/development sections)

### .dotfiles repository (final cleanup step)

- [x] Remove tool-owned files and directories now moved to `dotfiles-cli`:
  - `cmd/`
  - `internal/`
  - `go.mod`, `go.sum`
  - `Makefile`
  - `Dockerfile`, `Dockerfile.dev`
  - `.goreleaser.yaml`
  - `.schema/`
  - `.vscode/launch.json`
  - `.github/workflows/release.yml`
  - `scripts/install.sh` (replace with wrapper that installs dotfiles-cli)
- [x] Commit and push final config-only shape of `.dotfiles`

## Plan For Next Agent Session

1. **Iteration 1 complete**
  - dotfiles-cli extraction and split finished
  - config-only `.dotfiles` shape committed and pushed
  - installer docs aligned to public raw GitHub URL strategy

2. **Optional next iteration hardening**
  - add checksum verification to installer release downloads
  - run a full isolated bootstrap smoke test (`dotfiles bootstrap`) in disposable HOME

## Important Notes

- Keep V1 path behavior unchanged (`$HOME/.dotfiles`)
- Do not introduce path configurability in this phase
- Current tested installer target path from smoke test:
  - `/tmp/dotfiles-cli-smoke/bin/dotfiles`
- Latest relevant commits in `dotfiles-cli`:
  - `75e87c5` fix: support prerelease-only installs in private repo
  - `5753987` docs: align V1 scope and release metadata
  - `373c1e2` refactor: decouple CLI identity from personal dotfiles repo
