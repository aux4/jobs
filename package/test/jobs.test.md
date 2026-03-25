# jobs

```beforeAll
rm -rf .jobs
```

```afterAll
rm -rf .jobs
```

## run

### should start a background job

```execute
aux4 jobs run "echo hello world" | jq .
```

```expect:partial
{
  "id": "1",
  "command": "echo hello world",
  "pid": 0,
  "state": "RUNNING",
  "exitCode": 0,
  **
}
```

## output

### should show job stdout after completion

```execute
sleep 1 && aux4 jobs output 1
```

```expect
hello world
```

## status

### should show completed status

```execute
aux4 jobs status 1 | jq .
```

```expect:partial
{
  "id": "1",
  "command": "echo hello world",
  "pid": *?,
  "state": "COMPLETED",
  "exitCode": 0,
  "startTime": "*",
  "endTime": "*",
  "duration": "*",
  "dir": "*"
}
```

### should return error for unknown job

```execute
aux4 jobs status 999
```

```error:partial
Error: job 999 not found
```

## list

### should list all jobs

```execute
aux4 jobs list | jq .
```

```expect:partial
[
  {
    "id": "1",
    "command": "echo hello world",
    "pid": *?,
    "state": "COMPLETED",
    "exitCode": 0,
    **
  }
]
```

### should filter by state

```execute
aux4 jobs list --state COMPLETED | jq .
```

```expect:partial
[
  {
    "id": "1",
    **
    "state": "COMPLETED",
    **
  }
]
```

### should return empty array when no match

```execute
aux4 jobs list --state RUNNING | jq .
```

```expect
[]
```

## kill

### should kill a running job

```execute
aux4 jobs run "sleep 300" && sleep 1 && aux4 jobs kill 2
```

```expect:partial
job 2 killed
```

### should show killed state

```execute
aux4 jobs status 2 | jq .
```

```expect:partial
{
  "id": "2",
  "command": "sleep 300",
  "pid": *?,
  "state": "KILLED",
  "exitCode": -1,
  **
}
```

## killall

### should kill all running jobs

```execute
aux4 jobs run "sleep 301" && aux4 jobs run "sleep 302" && sleep 1 && aux4 jobs killall
```

```expect:partial
job 3 killed
job 4 killed
```

### should report no running jobs when none active

```execute
aux4 jobs killall
```

```expect
no running jobs
```

## failed job

### should capture non-zero exit code

```execute
aux4 jobs run "exit 42" && sleep 1 && aux4 jobs status 5 | jq .
```

```expect:partial
{
  "id": "5",
  "command": "exit 42",
  "pid": *?,
  "state": "FAILED",
  "exitCode": 42,
  **
}
```

## callbacks

### should run onSuccess callback on success

```execute
aux4 jobs run "echo hello" --onSuccess "echo onSuccess=called > .jobs/cb-success.txt" > /dev/null && sleep 1 && cat .jobs/cb-success.txt
```

```expect
onSuccess=called
```

### should run onFailure callback on failure

```execute
aux4 jobs run "exit 1" --onFailure "echo onFailure=called > .jobs/cb-failure.txt" > /dev/null && sleep 1 && cat .jobs/cb-failure.txt
```

```expect
onFailure=called
```

### should run onComplete callback on any state

```execute
aux4 jobs run "echo done" --onComplete "echo onComplete=called > .jobs/cb-complete.txt" > /dev/null && sleep 1 && cat .jobs/cb-complete.txt
```

```expect
onComplete=called
```

### should not run onSuccess on failure

```execute
aux4 jobs run "exit 1" --onSuccess "echo onSuccess=called > .jobs/cb-wrong.txt" > /dev/null && sleep 1 && test ! -f .jobs/cb-wrong.txt && echo "onSuccess=not called"
```

```expect
onSuccess=not called
```

### should set environment variables in callback

```execute
aux4 jobs run "echo hi" --onComplete 'echo state=$AUX4_JOB_STATE > .jobs/cb-env.txt' > /dev/null && sleep 1 && cat .jobs/cb-env.txt
```

```expect
state=COMPLETED
```
