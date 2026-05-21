# PPA Installation

Configures Personal Package Archives (PPA) using the modern DEB822 (`.sources`) format and installs specified packages.

## Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | The PPA identifier (e.g., `git-core/ppa`) or a descriptive name. |
| `key` | string | No | Raw GPG public key block. |
| `key_server` | string | No | The keyserver to fetch from (e.g., `keyserver.ubuntu.com`). |
| `key_id` | string | No | The GPG key fingerprint/ID to fetch. |
| `uri` | string | No | The full URI of the repository. Supports `:arch` placeholder. |
| `suites` | string | No | Ubuntu suite (e.g., `jammy`). Defaults to current OS suite. |
| `components` | string | No | Archive components (e.g., `main`). |
| `pkgs` | array | Yes | List of packages to install from this PPA. |

## Examples

### Standard Launchpad PPA
```yaml
groups:
  - name: git-ppa
    ppa:
      - name: "git-core/ppa"
        pkgs: ["git"]
```

### Manual Repository with Keyserver
```yaml
groups:
  - name: gh-cli
    ppa:
      - name: "gh-cli"
        uri: "https://cli.github.com/packages"
        key_server: "keyserver.ubuntu.com"
        key_id: "23F3D4AAAD4B6B4C"
        pkgs: ["gh"]
```

### Dynamic Architecture using `:arch`
```yaml
groups:
  - name: custom-repo
    ppa:
      - name: "myrepo"
        uri: "https://repo.example.com/:arch/debian"
        pkgs: ["my-package"]
```
