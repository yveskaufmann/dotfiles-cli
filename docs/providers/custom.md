# Custom Provider

Runs custom shell commands for install/update flows.

## Configuration

| Key | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Task name shown in output. |
| `install` | string | Yes | Install command. |
| `update` | string | No | Optional update command. |
| `installCheck` | string | No | Command used for idempotency check. Exit code `0` means already installed. |
| `systems` | string | No | Optional OS-specific constraint. |

## Example

```yaml
groups:
  - name: sdkman-bootstrap
    custom:
      - name: sdkman
        installCheck: "test -s $HOME/.sdkman/bin/sdkman-init.sh"
        install: "curl -fsSL https://get.sdkman.io | bash"
```
