#### Description

The `list` command outputs all jobs as a JSON array. Jobs are sorted by ID in ascending order. Filter the output by state, source tag, or both.

Possible job states are: `RUNNING`, `COMPLETED`, `FAILED`, `KILLED`.

The `--source` filter is useful when multiple agents share the same jobs directory: each agent can tag its jobs with a unique source on creation (`jobs run --source agent-a`) and later list only its own jobs (`jobs list --source agent-a`).

#### Usage

```bash
aux4 jobs list [--state <RUNNING|COMPLETED|FAILED|KILLED>] [--source <tag>] [--path <dir>]
```

--state   Filter jobs by state (default: show all)
--source  Filter jobs by source tag (default: show all)
--path    Custom jobs storage directory (default: `.jobs`)

#### Example

```bash
aux4 jobs list
```

```json
[
  {
    "id": "1",
    "command": "echo hello",
    "pid": 12345,
    "state": "COMPLETED",
    "exitCode": 0,
    "startTime": "2025-01-01T00:00:00Z",
    "endTime": "2025-01-01T00:00:00Z",
    "dir": "/home/user"
  }
]
```

```bash
aux4 jobs list --state RUNNING
```

```json
[
  {
    "id": "2",
    "command": "sleep 300",
    "pid": 12346,
    "state": "RUNNING",
    "exitCode": 0,
    "startTime": "2025-01-01T00:00:00Z",
    "dir": "/home/user"
  }
]
```

Filter by source tag:

```bash
aux4 jobs list --source agent-a
```

Combine filters:

```bash
aux4 jobs list --state COMPLETED --source agent-a
```
