# Brew Provider

Installs packages and casks via Homebrew.

## Configuration

`brew` items support object form:

| Key    | Type   | Required | Description                      |
| :----- | :----- | :------- | :------------------------------- |
| `name` | string | Yes      | Formula or cask name.            |
| `cask` | bool   | No       | Set to `true` for cask installs. |

`brew_taps` is also supported to add taps and optional package sets:

| Key    | Type   | Required | Description                          |
| :----- | :----- | :------- | :----------------------------------- |
| `name` | string | Yes      | Tap name.                            |
| `url`  | string | No       | Optional custom tap URL.             |
| `pkgs` | array  | No       | Optional package list from this tap. |

## Example

```yaml
groups:
  - name: mac-tools
    systems: darwin
    brew_taps:
      - name: homebrew/cask-fonts
    brew:
      - name: wget
      - name: iterm2
        cask: true
```
