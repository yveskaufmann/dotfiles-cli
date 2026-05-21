# NVM Provider

Installs and manages Node.js versions using NVM.

## Behavior

- installs NVM when missing
- installs the requested Node.js versions
- sets the default version
- skips versions that are already present

## Notes

- the provider expects an interactive shell environment
- it uses `bash -lc` so shell initialization is available
- it respects `NVM_DIR` when provided

## Example

```yaml
groups:
  - name: "Node.js"
    nvm:
      - default: "lts/iron"
```