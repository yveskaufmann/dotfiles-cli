# dotfiles-cli

`dotfiles-cli` is a standalone Go CLI for keeping developer tools consistent and manageable
across multiple machines and operating systems through declarative configuration.
It keeps executable logic separate from personal shell/config content.

## Motivation

- Keep developer tooling reproducible across laptops, workstations, and fresh installs.
- Manage tool setup declaratively (as YAML) so the same configuration applies across different operating systems without manual duplication.
- Version control your tool setup in your own dotfiles repository, giving you a complete audit trail and easy rollback of configuration changes.
- Reduce drift between machines by maintaining a single source of truth for tool versions and setup procedures.
- Keep personal dotfiles and executable bootstrap logic in separate repositories for clean separation of concerns.

## What are Dotfiles?

Dotfiles are configuration files for Unix tools (like `.bashrc`, `.gitconfig`, `.zshrc`). They typically live in your home directory and start with a dot (hence "dotfiles"). A dotfiles repository is a version-controlled collection of these configurations plus setup scripts. By storing them in Git, you get:

- **Consistency**: Replicate your exact environment across all machines
- **History**: Track changes to your configuration over time
- **Reproducibility**: Spin up new machines with tested, known-good settings
- **Collaboration**: Share configurations across your team or keep personal variations

## What It Does

- `bootstrap`: clones or updates the dotfiles repository into `$HOME/.dotfiles`
- `install`: installs tools from declarative `init/*.yaml` groups
- `link`: creates symlinks from `link/` into the home directory

## Configuration Concepts

Your dotfiles repository contains YAML files in `init/` that declaratively define which tools to install. Each file contains **groups** of tools grouped by install provider:

```yaml
groups:
  - name: Essential Tools
    profile: default
    systems: [linux, macos]
    apt:
      - git
      - curl
    npm:
      - eslint
      - prettier

  - name: Java Development
    profile: java
    sdkman:
      - java: "21"
      - maven: "3.9.0"
```

**Groups** are logical collections of tools that can be installed together. Each group can specify:
- `profile`: which command profiles include this group (e.g., `dotfiles install --profile java`)
- `systems`: which operating systems to apply to (linux, macos, etc.)
- **Providers**: sections like `apt`, `npm`, `sdkman` that handle the actual installation

**Providers** are the install mechanisms (package managers, version managers, GitHub releases). See [docs/providers/index.md](docs/providers/index.md) for complete provider list and configuration options.

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
- [Provider Configuration](docs/providers/index.md)
- [Installation Guide](docs/INSTALL.md)

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
- [RustUp](docs/providers/rustup.md)
- [Script](docs/providers/script.md)
- [Custom](docs/providers/custom.md)
