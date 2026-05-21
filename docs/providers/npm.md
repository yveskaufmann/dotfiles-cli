# NPM Provider

Installs global Node.js packages with npm.

## Configuration

`npm` entries support two forms:

1. String form: `"typescript"`
2. Object form: `{ name: "typescript" }`

Object keys:

| Key    | Type   | Required | Description              |
| :----- | :----- | :------- | :----------------------- |
| `name` | string | Yes      | Global npm package name. |

## Example

```yaml
groups:
  - name: web-tools
    npm:
      - "typescript"
      - name: "eslint"
      - name: "@angular/cli"
```
