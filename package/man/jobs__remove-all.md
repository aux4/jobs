#### Description

The `remove-all` command deletes all finished jobs from the storage directory. Running jobs are skipped — use `kill` or `killall` first to stop them.

You can scope the removal by `--state` or `--source` to clean up only matching jobs. This is useful when multiple agents share the same jobs directory and one wants to clean up only its own completed jobs.

#### Usage

```bash
aux4 jobs remove-all [--state <COMPLETED|FAILED|KILLED>] [--source <tag>] [--path <dir>]
```

--state   Only remove jobs in this terminal state
--source  Only remove jobs tagged with this source
--path    Custom jobs storage directory (default: `.jobs`)

#### Example

```bash
aux4 jobs remove-all --source agent-a
```

```text
job 4 removed
job 5 removed
job 7 removed
```

Remove only failed jobs:

```bash
aux4 jobs remove-all --state FAILED
```
