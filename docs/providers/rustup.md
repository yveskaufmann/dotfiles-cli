# RustUp Provider

Manages Rust toolchains, components, and targets via `rustup`.

## Overview

The RustUp provider automatically:
1. Installs `rustup` if not already present (via the official `sh.rustup.rs` installer)
2. Installs additional toolchains and sets a default toolchain
3. Adds components (for example `clippy`, `rustfmt`, `rust-analyzer`)
4. Adds cross-compilation targets
5. Skips anything already installed to avoid redundant work

## Configuration

The `rustup` key accepts a list of spec objects. All fields are optional; multiple entries are merged before installation.

| Key                | Type   | Required | Description                                                    |
| :----------------- | :----- | :------- | :------------------------------------------------------------- |
| `default_toolchain`| string | No       | Toolchain to set as the active default (for example `stable`). |
| `toolchains`       | array  | No       | Additional toolchains to install (for example `nightly`).      |
| `components`       | array  | No       | Components to add to the active toolchain.                     |
| `targets`          | array  | No       | Cross-compilation targets to add.                              |

## Behavior

### Setup

If `rustup` is not on `$PATH`, the provider runs the official installer:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
```

### Deduplication

Entries across multiple spec objects in the same group are merged. Duplicate component, target, and toolchain names are installed only once.

### Already-installed items

Before installing each component, target, or toolchain, the provider checks whether it is already present. Items that are already installed are reported as **skipped** (not an error) and the run continues.

### Error handling

A failed `rustup component add`, `rustup target add`, or `rustup toolchain install` halts the install immediately and reports `StatusFailed`.

## Examples

### Minimal — stable toolchain only

```yaml
groups:
  - name: Rust
    rustup:
      - default_toolchain: stable
```

### Common development setup

```yaml
groups:
  - name: Rust
    rustup:
      - default_toolchain: stable
        components:
          - rustfmt
          - clippy
          - rust-analyzer
```

### Multiple toolchains and a cross-compilation target

```yaml
groups:
  - name: Rust
    rustup:
      - default_toolchain: stable
        toolchains:
          - nightly
        components:
          - rustfmt
          - clippy
        targets:
          - aarch64-unknown-linux-gnu
```

### Split across entries (merged at install time)

```yaml
groups:
  - name: Rust base
    rustup:
      - default_toolchain: stable
        components:
          - rustfmt
          - clippy

  - name: Rust cross-compile
    rustup:
      - targets:
          - aarch64-unknown-linux-gnu
          - x86_64-unknown-linux-musl
```

## Task Reporting

The provider emits one task result per item:

- **`component <name>`**: component installation (success/skipped/failed)
- **`target <name>`**: target installation (success/skipped/failed)
- **`toolchain <name>`**: toolchain installation (success/skipped/failed)

## Notes

- `rustup` must be available on `$PATH` after setup. The provider sources `$HOME/.cargo/env` automatically during installation.
- The `default_toolchain` field only controls which toolchain is set as default via the existing `rustup` state; it does not install a toolchain on its own. Add it to `toolchains` as well if it may not be present.
- Targets and components apply to the currently active toolchain. Use `rustup override` in your shell configuration to pin directories to a specific toolchain.
