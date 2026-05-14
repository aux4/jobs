#### Description

The `clean` command removes all completed, failed, or killed jobs whose end time is older than the specified duration. Running jobs are always skipped.

This is useful for periodic housekeeping — clearing out stale jobs without affecting active work. The duration uses Go duration format (e.g. `24h`, `12h`, `30m`).

#### Usage

```bash
aux4 jobs clean [duration] [--path <dir>]
```

duration  How old a finished job must be to be removed (default: `24h`)
--path    Custom jobs storage directory (default: `.jobs`)

#### Example

```bash
aux4 jobs clean
```

Removes all finished jobs older than 24 hours.

```bash
aux4 jobs clean 12h
```

```text
job 2 removed
job 5 removed
```

Remove finished jobs older than 30 minutes:

```bash
aux4 jobs clean 30m
```
