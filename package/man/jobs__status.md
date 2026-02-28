#### Description

The `status` command shows detailed information about a specific job as JSON. It includes the job ID, command, PID, state, exit code, start time, end time, duration, and working directory.

If the job has finished but the monitor process didn't update the state (e.g., it was killed), the status command detects this by checking whether the PID is still running and updates the state accordingly.

#### Usage

```bash
aux4 jobs status <id>
```

--id  The job ID (positional argument)

#### Example

```bash
aux4 jobs status 1
```

```json
{
  "id": "1",
  "command": "echo hello",
  "pid": 12345,
  "state": "COMPLETED",
  "exitCode": 0,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-01-01T00:00:00Z",
  "duration": "50ms",
  "dir": "/Users/me/project"
}
```
