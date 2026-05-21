# Pipx Provider

Installs Python CLI tools with pipx.

## Configuration

`pipx` entries support two forms:

1. String form: `"uv"`
2. Object form: `{ name: "uv" }`

Object keys:

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | pipx package name. |

## Example

```yaml
groups:
  - name: python-tools
    pipx:
      - "black"
      - name: "uv"
      - name: "yt-dlp"
```
