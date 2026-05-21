# SdkMan Provider

Installs JVM ecosystem tools via SDKMAN.

## Configuration

Each item supports one candidate and one or more versions:

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `candidate` | string | No | SDKMAN candidate name (for example `java`, `maven`). |
| `version` | string | No | Default version to install/use. |
| `versions` | array | No | Additional versions to install. |

## Example

```yaml
groups:
  - name: jvm
    sdkman:
      - candidate: java
        version: "21.0.2-open"
      - candidate: maven
        version: "3.9.10"
```
