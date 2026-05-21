# Scripts & Custom Installation

Execute inline shell commands or reference external scripts for complex tasks.

## Custom Installation

User-defined shell commands with custom idempotency checks.

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Name of the task. |
| `install` | string | Yes | The shell command to execute for installation. |
| `installCheck` | string | No | A command to verify if already installed. Exit code `0` skips installation. |

## External Scripts

Executes an existing script file.

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | Name of the script. |
| `script` | string | Yes | Relative path to the script from the root (e.g., `./init/99-work.sh`). |

## Examples

### Custom Inline Logic
```yaml
groups:
  - name: dummy-tool
    custom:
      - name: "Dummy Tool"
        install: "touch /tmp/dummy"
        installCheck: "ls /tmp/dummy"
```

### Reference Legacy Scripts
```yaml
groups:
  - name: university-tools
    script:
      - name: "99-university.sh"
        script: "./init/99-university.sh"
```
