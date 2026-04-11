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

### Callbacks

Run a command automatically when a job finishes using `--onSuccess`, `--onFailure`, or `--onComplete`:

```bash
aux4 jobs run "npm test" --onSuccess "echo Tests passed" --onFailure "echo Tests FAILED"
```

Use `--onComplete` to run a command on any terminal state (success, failure, or killed):

```bash
aux4 jobs run "deploy.sh" --onComplete "echo Job finished with state $AUX4_JOB_STATE"
```

Callbacks run via `sh -c` in the job's working directory. If both a specific callback (`onSuccess`/`onFailure`) and `onComplete` are set, the specific one runs first.

#### Callback Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `AUX4_JOB_ID` | Job ID | `3` |
| `AUX4_JOB_STATE` | Terminal state | `COMPLETED`, `FAILED`, `KILLED` |
| `AUX4_JOB_EXIT_CODE` | Exit code | `0`, `1`, `-1` |
| `AUX4_JOB_COMMAND` | Original command | `npm test` |
| `AUX4_JOB_DIR` | Working directory | `/home/user/project` |

### Source tags

Tag jobs with a source identifier so multiple agents sharing the same jobs directory can list only their own jobs:

```bash
aux4 jobs run "long-task.sh" --source agent-a
aux4 jobs list --source agent-a
```

### Auto-cleanup

Use `--cleanup true` to remove the job directory automatically after callbacks finish:

```bash
aux4 jobs run "send-notification.sh" --cleanup true
```

This is useful for fire-and-forget jobs whose output doesn't need to be retained.

### Custom storage path

Use `--path` on any command to use a different storage directory than the default `.jobs`:

```bash
aux4 jobs run "build.sh" --path .my-jobs
aux4 jobs list --path .my-jobs
```

This lets multiple agents fully isolate their jobs from each other.

### Removing jobs

Remove a single finished job:

```bash
aux4 jobs remove 3
```

Force-remove a running job (kills it first):

```bash
aux4 jobs remove 3 --force true
```

Remove all finished jobs created by a specific source:

```bash
aux4 jobs remove-all --source agent-a
```

Remove only failed jobs:

```bash
aux4 jobs remove-all --state FAILED
```

## Commands

| Command | Description |
|---------|-------------|
| `aux4 jobs run <command> [--onSuccess] [--onFailure] [--onComplete] [--source] [--cleanup] [--path]` | Run a command in the background with optional callbacks, source tag, auto-cleanup, and storage path |
| `aux4 jobs list [--state] [--source] [--path]` | List jobs, optionally filtered by state and/or source |
| `aux4 jobs status <id> [--path]` | Show job status with exit code and duration |
| `aux4 jobs output <id> [--stream] [--path]` | Show full job output |
| `aux4 jobs tail <id> [--stream] [--path]` | Tail job output in real time |
| `aux4 jobs kill <id> [--path]` | Kill a running job |
| `aux4 jobs killall [--path]` | Kill all running jobs |
| `aux4 jobs remove <id> [--force] [--path]` | Remove a finished job from storage |
| `aux4 jobs remove-all [--state] [--source] [--path]` | Remove all finished jobs (optionally filtered) |
| `aux4 jobs on <id> [--success] [--failure] [--complete]` | Register callbacks on a running job |
| `aux4 jobs attach <pid> <command> [--stdout] [--stderr] [--source] [--path]` | Attach an external process to the jobs system |

## How It Works

Jobs are stored in `.jobs/` in the current directory by default (override with `--path`). Each job gets its own directory containing:

- `job.json` — job metadata (command, PID, state, timestamps, source, cleanup flag)
- `stdout` — captured standard output
- `stderr` — captured standard error

A background monitor process waits for the command to finish and updates the job state with the exit code and end time. If `cleanup: true` was set, the directory is removed after callbacks run.
