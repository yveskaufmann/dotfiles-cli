# GitHub Release Installation

Downloads and installs binaries directly from GitHub Releases.

## Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | string | Yes | The tool name (used for idempotency check via `type`). |
| `repo` | string | Yes | The repository in `owner/repo` format. |
| `version` | string | No | Specific tag name. Defaults to `latest`. |
| `asset_pattern`| string | Yes | Regex pattern to match the asset filename (e.g., `linux_amd64.tar.gz$`). |
| `binaries` | array | No | List of files to extract and install from the archive. |
| `binary_path` | string | No | Relative path inside the archive to the single binary. |
| `install_path` | string | No | Destination directory. Defaults to `/usr/local/bin` (root) or `~/bin` (user). |

## Examples

### Single Binary Release
```yaml
groups:
  - name: yq
    github_release:
      - name: yq
        repo: mikefarah/yq
        asset_pattern: "yq_linux_amd64$"
```

### Multi-Binary Archive Extraction
Extracts multiple tools from a single download.
```yaml
groups:
  - name: go-containerregistry
    github_release:
      - name: crane
        repo: google/go-containerregistry
        asset_pattern: "go-containerregistry_Linux_x86_64.tar.gz"
        binaries: ["crane", "gcrane", "krane"]
```

### Binary inside a Subdirectory
```yaml
groups:
  - name: helm
    github_release:
      - name: helm
        repo: helm/helm
        asset_pattern: "linux-amd64.tar.gz$"
        binary_path: "linux-amd64/helm"
```
