# Provider Configuration Reference

The installer is configured through YAML files under `init/`.

## Configuration Concept

Each config file defines installation `groups`. A group can include one or more provider sections.

Top-level keys:

| Key       | Type   | Description                               |
| :-------- | :----- | :---------------------------------------- |
| `$schema` | string | Optional schema path for YAML validation. |
| `groups`  | array  | Ordered installation groups.              |

Group keys:

| Key       | Type   | Required | Description                                     |
| :-------- | :----- | :------- | :---------------------------------------------- |
| `name`    | string | Yes      | Group identifier for logging and reporting.     |
| `profile` | string | No       | Profile filter (for example `default`, `work`). |
| `systems` | string | No       | OS filter (for example `linux`, `darwin`).      |

Example:

```yaml
$schema: "./.schema/config.schema.json"
groups:
	- name: base-tools
		profile: default
		systems: linux
		apt:
			- git
			- curl
		github_release:
			- name: yq
				repo: mikefarah/yq
				asset_pattern: "yq_linux_amd64$"
```

## Provider Pages

- [APT](apt.md)
- [PPA](ppa.md)
- [Brew](brew.md)
- [GitHub Release](github.md)
- [Binary](binary.md)
- [NVM](nvm.md)
- [SdkMan](sdkman.md)
- [NPM](npm.md)
- [Pipx](pipx.md)
- [Snap](snap.md)
- [JetBrains Toolbox](jetbrains.md)
- [RustUp](rustup.md)
- [Script](script.md)
- [Custom](custom.md)

## Notes

- Providers execute by priority in the runtime registry.
- System packages run before version managers and application-level tools.
- Unsupported or unregistered providers should not be used in production configs.
