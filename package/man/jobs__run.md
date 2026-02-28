#### Description

The `run` command starts a command in the background and returns immediately with a JSON object containing the job details. The command runs in a separate process with its stdout and stderr captured to files. A monitor process tracks the command and updates the job state when it completes.

The full command string is passed to `sh -c`, so shell features like pipes, redirects, and `&&` chains work as expected.

#### Usage

```bash
aux4 jobs run <command>
```

--command  The full command to run in the background (positional argument)

#### Example

```bash
aux4 jobs run "sleep 5 && echo done"
```

```json
{
  "id": "1",
  "command": "sleep 5 && echo done",
  "pid": 0,
  "state": "RUNNING",
  "exitCode": 0,
  "startTime": "2025-01-01T00:00:00Z",
  "dir": "/Users/me/project"
}
```
