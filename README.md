# dotfiles-cli

`dotfiles-cli` is a standalone Go CLI for keeping developer tools consistent and manageable
across multiple machines and operating systems through declarative configuration.
It keeps executable logic separate from personal shell/config content.

## Motivation

- Keep developer tooling reproducible across laptops, workstations, and fresh installs.
- Manage tool setup declaratively so the same config can be applied on different OSes.
- Reduce drift between machines by reusing one installation definition (`init/*.yaml`).
- Keep personal dotfiles and executable bootstrap logic in separate repositories.

## What It Does

- `bootstrap`: clones or updates the dotfiles repository into `$HOME/.dotfiles`
- `install`: installs tools from declarative `init/*.yaml` groups
- `link`: creates symlinks from `link/` into the home directory

## V1 Scope

- Local workspace path is fixed to `$HOME/.dotfiles`.
- Repository source is configurable via `--repository`.
- Custom local paths are intentionally deferred to a later version.

## Installation

Install the latest release binary:

```bash
curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh
```

More details and uninstall instructions: [docs/INSTALL.md](docs/INSTALL.md)

## Usage

```bash
dotfiles bootstrap
dotfiles bootstrap --repository git@github.com:user/dotfiles.git
dotfiles install --profile default
dotfiles link --dry-run
```

## Build From Source

```bash
make build
./bin/dotfiles version
```

## Documentation

- [Architecture](architecture.md)
- [Agent Documentation](docs/agent.md)
- [Provider Configuration](docs/providers/index.md)
- [Task Backlog](tasks/backlog.md)

Provider pages:

- [APT](docs/providers/apt.md)
- [PPA](docs/providers/ppa.md)
- [Brew](docs/providers/brew.md)
- [GitHub Release](docs/providers/github.md)
- [Binary](docs/providers/binary.md)
- [NVM](docs/providers/nvm.md)
- [SdkMan](docs/providers/sdkman.md)
- [NPM](docs/providers/npm.md)
- [Pipx](docs/providers/pipx.md)
- [Snap](docs/providers/snap.md)
- [JetBrains Toolbox](docs/providers/jetbrains.md)
- [Script](docs/providers/script.md)
- [Custom](docs/providers/custom.md)

## Repository Layout

```text
cmd/            CLI entry point
internal/       bootstrap, install, link, config, providers, and utilities
docs/           human-facing documentation
.agents/        planning, tasks, and learnings for agent workflows
tasks/          backlog notes
```
