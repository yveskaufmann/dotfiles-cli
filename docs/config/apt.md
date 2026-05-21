# APT Installation

Installs standard packages via `apt-get`.

## Parameters

- `apt`: An array of packagenames (strings).

## Examples

### Standard Packages
```yaml
groups:
  - name: base-tools
    apt:
      - curl
      - git
      - vim
```

### Profile-specific Tools
```yaml
groups:
  - name: work-tools
    profile: work
    apt:
      - jq
      - tree
```
