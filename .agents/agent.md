# Dotfiles Declarative Bootstrapping Agent

This agent manages the transition from imperative Bash-based tool installation to a declarative Go-based system.

## Architecture Overview

The system is designed to be idempotent, sequential, and interactive. It replaces the logic previously held in `init/*.sh` files with `*.yaml` configurations.

### Core Components

1.  **Configuration Schema (`internal/init/setup-tools/types.go`)**:
    -   Defines the `Config` and `Group` structures.
    -   Supports multiple package managers: `apt`, `snap`, `pipx`, `npm`, `ppa`, `github_release`.
    -   Includes `custom` and `script` for legacy or complex tasks.
    -   Filters by `profile` (default, work) and `systems` (linux, darwin).

2.  **Sequential Loader (`internal/init/setup-tools/loader.go`)**:
    -   Reads all `*.yaml` files from the `init/` directory.
    -   Sorts them alphanumerically (e.g., `01-tools.yaml` before `02-develop.yaml`).
    -   Filters groups based on runtime environment (OS, profile).

3.  **Execution Engine (`internal/init/setup-tools/executor.go`)**:
    -   Executes tool installations one by one.
    -   Stops immediately on error.
    -   Ensures `os.Stdin`, `os.Stdout`, and `os.Stderr` are connected to allow for `sudo` prompts and interactive scripts.

4.  **Tool Managers (`internal/tools/`)**:
    -   Individual wrappers for `apt`, `snap`, etc.
    -   **Idempotency Logic**:
        -   `apt`: Checks `dpkg-query -W`.
        -   `snap`: Checks `snap list`.
        -   `pipx`: Checks `pipx list`.
        -   `github_release`: Uses regex patterns to find assets and verifies via `type`.
        -   `ppa`: 2-phase check (Check `.sources` or `.list` file, then check package).

## Operational Protocol

### Task Management
Before starting any task, the agent MUST:
1.  Check the `/tasks/backlog.md` file for open tasks and technical debt.
2.  Prioritize tasks based on the current sprint goals.
3.  Update the task status to "in-progress" in `/tasks/backlog.md` (or the respective active task file).
4.  Upon completion, mark the task as "completed" and append any relevant notes or results.

### Configuration Schema
All YAML configurations in `init/*.yaml` MUST be validated against the JSON schema located at `/.schema/config.schema.json`. When introducing new fields to `types.go`, the schema must be updated immediately.

## Idempotency Policy

Unless explicitly overridden by a `check` or `installCheck` property in the YAML, the system follows these default checks:
-   **Package Managers**: Query the manager's database of installed packages.
-   **Binaries/Releases**: Verify via `type <name>`.
-   **Custom**: Always run unless `check` command returns success (exit 0).

## Execution Environment

The CLI is expected to be run in an interactive terminal. `sudo` prompts will be handled by the underlying package manager or shell command by inheriting the parent process's file descriptors.
