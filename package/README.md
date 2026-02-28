# @aux4/jobs

Background job manager for aux4. Run commands in the background, monitor their status, and retrieve their output.

## Installation

```bash
aux4 aux4 pkger install aux4/jobs
```

## Usage

### Run a background job

```bash
aux4 jobs run "sleep 10 && echo done"
```

```json
{
  "id": "1",
  "command": "sleep 10 && echo done",
  "pid": 0,
  "state": "RUNNING",
  "exitCode": 0,
  "startTime": "2025-01-01T00:00:00Z",
  "dir": "/Users/me/project"
}
```

### List all jobs

```bash
aux4 jobs list
```

```json
[
  {
    "id": "1",
    "command": "sleep 10 && echo done",
    "pid": 12345,
    "state": "RUNNING",
    "exitCode": 0,
    "startTime": "2025-01-01T00:00:00Z",
    "dir": "/Users/me/project"
  }
]
```

Filter by state:

```bash
aux4 jobs list --state RUNNING
```

### Check job status

```bash
aux4 jobs status 1
```

```json
{
  "id": "1",
  "command": "sleep 10 && echo done",
  "pid": 12345,
  "state": "COMPLETED",
  "exitCode": 0,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-01-01T00:00:10Z",
  "duration": "10s",
  "dir": "/Users/me/project"
}
```

### Get job output

```bash
aux4 jobs output 1
```

```text
done
```

Show stderr instead:

```bash
aux4 jobs output 1 --stream stderr
```

### Tail job output

Stream output in real time while the job is running:

```bash
aux4 jobs tail 1
```

Tail stderr:

```bash
aux4 jobs tail 1 --stream stderr
```

### Kill a running job

```bash
aux4 jobs kill 1
```

```text
job 1 killed
```

### Kill all running jobs

```bash
aux4 jobs killall
```

```text
job 1 killed
job 3 killed
```

## Commands

| Command | Description |
|---------|-------------|
| `aux4 jobs run <command>` | Run a command in the background |
| `aux4 jobs list [--state <RUNNING\|COMPLETED\|FAILED\|KILLED>]` | List jobs, optionally filtered by state |
| `aux4 jobs status <id>` | Show job status with exit code and duration |
| `aux4 jobs output <id> [--stream stdout\|stderr]` | Show full job output |
| `aux4 jobs tail <id> [--stream stdout\|stderr]` | Tail job output in real time |
| `aux4 jobs kill <id>` | Kill a running job |
| `aux4 jobs killall` | Kill all running jobs |

## How It Works

Jobs are stored in `.jobs/` in the current directory. Each job gets its own directory containing:

- `job.json` — job metadata (command, PID, state, timestamps)
- `stdout` — captured standard output
- `stderr` — captured standard error

A background monitor process waits for the command to finish and updates the job state with the exit code and end time.
