#### Description

The `run` command starts a command in the background and returns immediately with a JSON object containing the job details. The command runs in a separate process with its stdout and stderr captured to files. A monitor process tracks the command and updates the job state when it completes.

The full command string is passed to `sh -c`, so shell features like pipes, redirects, and `&&` chains work as expected.

Optional callback commands can be specified to execute automatically when the job reaches a terminal state. Callbacks run via `sh -c` in the job's working directory and receive job metadata through environment variables.

#### Usage

```bash
aux4 jobs run <command> [--onSuccess <command>] [--onFailure <command>] [--onComplete <command>]
```

--command     The full command to run in the background (positional argument)
--onSuccess   Command to run when the job succeeds (exit code 0). Default: empty
--onFailure   Command to run when the job fails (non-zero exit code). Default: empty
--onComplete  Command to run when the job finishes, regardless of outcome. Default: empty

#### Callback Execution Order

1. If the job succeeds, `onSuccess` runs first, then `onComplete`.
2. If the job fails, `onFailure` runs first, then `onComplete`.
3. If the job is killed, only `onComplete` runs.

Callback errors are logged to the job's stderr file but do not change the job state.

#### Environment Variables

Callback commands receive the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `AUX4_JOB_ID` | Job ID | `3` |
| `AUX4_JOB_STATE` | Terminal state | `COMPLETED`, `FAILED`, `KILLED` |
| `AUX4_JOB_EXIT_CODE` | Exit code | `0`, `1`, `-1` |
| `AUX4_JOB_COMMAND` | Original command | `npm test` |
| `AUX4_JOB_DIR` | Working directory | `/home/user/project` |

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

Run with a success callback:

```bash
aux4 jobs run "npm test" --onSuccess "echo Tests passed"
```

Run with different actions for success and failure:

```bash
aux4 jobs run "deploy.sh" \
  --onSuccess "curl -X POST https://slack.com/webhook -d 'Deploy succeeded'" \
  --onFailure "curl -X POST https://slack.com/webhook -d 'Deploy FAILED'"
```

Run with an onComplete callback that uses environment variables:

```bash
aux4 jobs run "npm test" --onComplete "echo Job $AUX4_JOB_ID finished with state $AUX4_JOB_STATE"
```
