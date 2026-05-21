# Snap, Pipx, & NPM Installation

Install packages using specialized language or container managers.

## Format

These tools support two formats for list items:
1.  **String**: Short-form package name.
2.  **Object**: Struct with `name` (and potentially other keys in the future).

## Examples

### Snap Packages (Ubuntu only)
```yaml
groups:
  - name: cloud-tools
    snap:
      - "aws-cli"
      - name: "google-cloud-sdk"
```

### Pipx Packages
```yaml
groups:
  - name: python-tools
    pipx:
      - "black"
      - name: "yt-dlp"
```

### NPM Packages
```yaml
groups:
  - name: web-tools
    npm:
      - "typescript"
      - name: "@angular/cli"
```
