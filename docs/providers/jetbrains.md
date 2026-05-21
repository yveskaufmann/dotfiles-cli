# JetBrains Toolbox Provider

Installs JetBrains IDEs through the JetBrains provider.

## Configuration

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `ide` | string | Yes | IDE identifier used by the provider. |
| `version` | string | No | Optional version selector. |
| `name` | string | No | Optional display label. |

## Example

```yaml
groups:
  - name: jetbrains
    jetbrains:
      - ide: goland
      - ide: intellij-idea-ultimate
        version: "2025.1"
```
