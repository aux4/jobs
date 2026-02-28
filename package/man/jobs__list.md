#### Description

The `list` command outputs all jobs as a JSON array. Jobs are sorted by ID in ascending order. Use `--state` to filter by a specific state.

Possible job states are: `RUNNING`, `COMPLETED`, `FAILED`, `KILLED`.

#### Usage

```bash
aux4 jobs list [--state <RUNNING|COMPLETED|FAILED|KILLED>]
```

--state  Filter jobs by state (default: show all)

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
