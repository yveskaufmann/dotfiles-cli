# Script Provider

Runs external scripts from your repository.

## Configuration

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Task label used in output. |
| `script` | string | Yes | Script path relative to repository root. |

## Example

```yaml
groups:
  - name: legacy-hooks
    script:
      - name: university-setup
        script: ./init/99-university.sh
```
