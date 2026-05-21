# Snap Provider

Installs snap packages.

## Configuration

`snap` entries support two forms:

1. String form: `"kubectl"`
2. Object form: `{ name: "kubectl" }`

Object keys:

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Snap package name. |

## Example

```yaml
groups:
  - name: cloud-tools
    snap:
      - "kubectl"
      - name: "helm"
```
