# Configuration Reference

The dotfiles bootstrapper is configured via YAML files in the `init/` directory.

## Core Schema

Each file must follow the `Config` structure, which contains a list of `groups`.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `$schema` | string | (Optional) Path to the JSON schema for validation. |
| `groups` | array | A list of installation groups. |

### Group Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Unique name for the group. |
| `profile` | string | No | Filters by profile (`default`, `work`, `university`). |
| `systems` | string | No | Filters by OS (`linux`, `darwin`). |

## Installation Tools

Each group can contain one or more of the following installation tool definitions:

- [APT](apt.md) - Standard Debian/Ubuntu packages.
- [PPA](ppa.md) - Personal Package Archives (DEB822).
- [GitHub Releases](github.md) - Binaries from GitHub.
- [Binary Downloads](binary.md) - Direct URLs to binaries or archives.
- [Snap, Pipx, & NPM](snap_pipx_npm.md) - Language and container-specific managers.
- [NVM](../providers/nvm.md) - Node Version Manager for managing multiple Node.js versions.
- [Scripts & Custom](scripts.md) - Legacy bash scripts and inline commands.
