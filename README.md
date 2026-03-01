# My personal dotfiles

Dotfiles manager and bootstrap tool for configuring development environments across multiple machines.

## Supported Platforms

- Linux (Ubuntu, Debian, Arch, etc.)
- macOS (Intel & Apple Silicon)
- Windows (WSL2)

## Quick Install

Install from your private repository (requires GitHub CLI for authentication):

```bash
  cd ~/.dotfiles
  git clone git@github.com:yveskaufmann/.dotfiles.git . 
./scripts/install.sh --github-token "$(gh auth token)"
```

Install the dotfiles CLI tool from public repository (no authentication required):
# TODO: I need a location where I can publish the install script without authentication, or make it available in a public repository. 
# For now, I will need to clone the repository first and run the install script locally.

```bash
curl -fsSL https://raw.githubusercontent.com/yveskaufmann/.dotfiles/main/scripts/install.sh | sh
```

Then bootstrap your dotfiles:

```bash
dotfiles bootstrap
```

That's it! The tool will:

1. Clone your dotfiles repository (prompts for URL if not configured)
2. Install tools and packages from your init scripts
3. Create symlinks from your dotfiles to your home directory

## Detailed Installation

See [INSTALL.md](docs/INSTALL.md) for detailed installation instructions, including:

- Custom installation directories
- System-wide installation
- Manual installation
- Troubleshooting

## Usage

### Bootstrap Everything

```bash
# Full bootstrap (clone/pull, install tools, create symlinks)
dotfiles bootstrap

# Bootstrap with custom repository
dotfiles bootstrap --repository git@github.com:yourusername/dotfiles.git

# Bootstrap for specific profile
dotfiles bootstrap --profile work
```

### Individual Commands

```bash
# Install tools only
dotfiles install --profile default

# Create symlinks only
dotfiles link

# Dry-run to see what would be linked
dotfiles link --dry-run
```

### Configuration

Configuration is stored in `~/.config/.dotfiles/config.yaml`:

```yaml
dotfiles:
  repository:
    url: git@github.com:yveskaufmann/.dotfiles.git
    type: ssh
```

## Project Structure

```text
.dotfiles/
├── cmd/              # CLI application entry point
├── internal/         # Internal packages
│   ├── cli/         # Command implementations
│   ├── engine/      # Core bootstrap/link/install logic
│   ├── config/      # Configuration management
│   └── provider/    # Package manager providers
├── init/            # Tool installation definitions
├── link/            # Files to symlink to home
├── scripts/         # Utility scripts
└── docs/            # Documentation
```

## Development

### Building from Source

```bash
make build
./bin/dotfiles version
```

### Running Tests

```bash
make test
```

## License

This project is licensed under the terms specified in the LICENSE file.
