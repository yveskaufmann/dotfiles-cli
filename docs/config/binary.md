# Binary Installation

Downloads and installs binaries from a direct URL.

## Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | The tool name. |
| `url` | string | Yes | The download URL. Supports `:version` placeholder. |
| `version` | string | No | The version string to substitute into the URL. |
| `binaries` | array | No | List of files to extract if the URL is an archive. |
| `binary_path`| string | No | Path to binary inside archive. |
| `install_path`| string | No | Destination directory. |

## Examples

### Direct Binary URL
```yaml
groups:
  - name: kubectl
    binary:
      - name: kubectl
        url: "https://dl.k8s.io/release/:version/bin/linux/amd64/kubectl"
        version: "v1.35.0"
```

### Direct Archive URL
```yaml
groups:
  - name: go
    binary:
      - name: go
        url: "https://go.dev/dl/go1.25.6.linux-amd64.tar.gz"
        install_path: "/opt"
        binary_path: "go"
```
