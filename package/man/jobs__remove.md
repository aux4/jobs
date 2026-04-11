#### Description

The `remove` command deletes a finished job and its output files from the jobs storage directory. Use this to clean up after a job has completed, failed, or been killed.

By default, removing a running job is refused to prevent accidental data loss. Pass `--force true` to kill the job and remove it in one step.

#### Usage

```bash
aux4 jobs remove <id> [--force <true|false>] [--path <dir>]
```

--force  Kill the job first if still running (default: `false`)
--path   Custom jobs storage directory (default: `.jobs`)

#### Example

```bash
aux4 jobs remove 3
```

```text
job 3 removed
```

Force-remove a running job:

```bash
aux4 jobs remove 7 --force true
```
