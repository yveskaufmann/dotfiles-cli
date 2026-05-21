# NVM Provider

Installs and manages Node.js versions using NVM (Node Version Manager).

## Overview

The NVM provider automatically:
1. Installs NVM if not already present (from the latest GitHub release)
2. Installs specified Node.js versions
3. Sets a default Node.js version
4. Verifies existing installations to avoid redundant downloads
5. Automatically deduplicates versions (default is always included)

## Configuration

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `default` | string | Yes | The default Node.js version to use and install. |
| `versions` | array | No | Additional Node.js versions to install. The default version is automatically included. |

### Version Specifiers

NVM supports various version formats:
- **LTS Codenames**: `lts/iron`, `lts/hydrogen`, `lts/gallium`
- **Latest**: `latest` (installs the most recent Node.js version)
- **Major Version**: `18`, `20`, `21` (installs the latest patch of that major version)
- **Specific Version**: `20.11.0`, `18.19.1` (installs exact version)

## Behavior

### Installation Process

1. **NVM Installation**: Checks if `~/.nvm/nvm.sh` exists. If not, downloads and runs the official NVM installer from the latest GitHub release tag.

2. **Version Deduplication**: Automatically includes the `default` version in the installation list, so you don't need to specify it in `versions`.

3. **Version Verification**: Before installing each version, checks if it's already installed using `nvm ls`. Skips installation if the version is present.

4. **Default Version**: After all versions are installed, sets the specified `default` version using `nvm alias default`.

### Error Handling

- Installation **stops immediately** if any Node.js version fails to install
- Failed installations report detailed error messages
- NVM installation failures halt the entire provider

### Shell Integration

The provider **does not** modify shell configuration files (.bashrc, .zshrc, etc.). It expects NVM to be sourced in your shell environment, which is typically handled by the NVM installer itself or your dotfiles' shell initialization scripts.

### NVM Directory

The provider respects the `NVM_DIR` environment variable. If not set, it defaults to `~/.nvm`.

## Examples

### Minimal Configuration

Install only the LTS Iron version and set it as default:

```yaml
groups:
  - name: "Node.js"
    nvm:
      - default: "lts/iron"
```

### Multiple Versions

Install LTS Iron (default) and the latest version:

```yaml
groups:
  - name: "Node.js"
    nvm:
      - default: "lts/iron"
        versions:
          - "latest"
```

### Development Setup

Install multiple LTS versions and a specific version:

```yaml
groups:
  - name: "Node.js Development"
    nvm:
      - default: "lts/iron"
        versions:
          - "lts/hydrogen"
          - "lts/gallium"
          - "20.11.0"
          - "latest"
```

### Version Number Default

Use a specific major version as default:

```yaml
groups:
  - name: "Node.js"
    nvm:
      - default: "20"
        versions:
          - "18"
          - "latest"
```

## Task Reporting

The provider reports progress for each operation:

- **`nvm`**: NVM installation (success/failed/skipped)
- **`node@<version>`**: Node.js version installation (success/failed/skipped)
- **`node@<version> (set as default)`**: Setting default version (success/failed/skipped)

Example output:
```
✓ nvm (installed)
✓ node@lts/iron (installed)
⊘ node@latest (already installed)
✓ node@lts/iron (set as default)
```

## Notes

- NVM must be sourced in your shell to use Node.js. This is typically done in `.bashrc` or `.zshrc` via your dotfiles.
- The provider uses `bash -lc` to execute NVM commands, ensuring login shell initialization.
- NVM updates are not handled by this provider. To update NVM, you must do so manually or remove `~/.nvm` and re-run.
- Version deduplication ensures the `default` version appears only once in the installation list, even if explicitly specified in `versions`.

## Troubleshooting

### NVM not found in shell

If `nvm` command is not available after installation:
1. Ensure your shell initialization files (`.bashrc`, `.zshrc`) source NVM:
   ```bash
   export NVM_DIR="$HOME/.nvm"
   [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
   ```
2. Restart your shell or source the configuration: `source ~/.bashrc`

### Version installation fails

- Check your internet connection
- Verify the version specifier is valid (run `nvm ls-remote` to see available versions)
- Ensure you have sufficient disk space

### Default not switching

- Verify the version was successfully installed first
- Check `nvm alias default` manually to see current setting
- Restart your shell after setting the default
