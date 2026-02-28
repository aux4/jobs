#### Description

The `kill` command terminates a running job. On Unix systems, it sends SIGTERM to the entire process group, ensuring that the command and all its child processes are stopped. On Windows, it uses `taskkill /F /T` to forcefully terminate the process tree.

After killing, the job state is set to `KILLED` with exit code `-1`. If the job is not running, the command returns an error.

#### Usage

```bash
aux4 jobs kill <id>
```

--id  The job ID (positional argument)

#### Example

```bash
aux4 jobs kill 2
```

```text
job 2 killed
```

Verify the state:

```bash
aux4 jobs status 2
```

```json
{
  "id": "2",
  "state": "KILLED",
  "exitCode": -1
}
```
