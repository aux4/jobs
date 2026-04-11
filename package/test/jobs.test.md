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

## source

### should tag job with source

```execute
aux4 jobs run "echo tagged" --source agent-x | jq -r .source
```

```expect
agent-x
```

### should filter list by source

```execute
aux4 jobs run "echo y1" --source agent-y > /dev/null && sleep 1 && aux4 jobs list --source agent-y | jq '. | length'
```

```expect:regex
^[1-9][0-9]*$
```

### should return empty when no source matches

```execute
aux4 jobs list --source nonexistent
```

```expect
[]
```

## remove

### should remove a finished job

```execute
aux4 jobs run "echo to-remove" > /tmp/jobs-remove-id.json && sleep 1 && ID=$(jq -r .id /tmp/jobs-remove-id.json) && aux4 jobs remove $ID && rm -f /tmp/jobs-remove-id.json
```

```expect:partial
job *? removed
```

### should fail to remove a running job without force

```execute
aux4 jobs run "sleep 5" > /tmp/jobs-running-id.json && ID=$(jq -r .id /tmp/jobs-running-id.json) && aux4 jobs remove $ID 2>&1 ; aux4 jobs kill $ID > /dev/null ; rm -f /tmp/jobs-running-id.json
```

```expect:partial
Error: job *? is still running
```

### should remove a running job with force

```execute
aux4 jobs run "sleep 10" > /tmp/jobs-force-id.json && ID=$(jq -r .id /tmp/jobs-force-id.json) && aux4 jobs remove $ID --force true && rm -f /tmp/jobs-force-id.json
```

```expect:partial
job *? removed
```

## cleanup

### should auto-remove job after completion when cleanup is true

```execute
aux4 jobs run "echo will-vanish" --source vanish-test --cleanup true > /dev/null && sleep 1 && aux4 jobs list --source vanish-test
```

```expect
[]
```

## custom path

### should use custom storage path

```execute
rm -rf /tmp/custom-jobs-test && aux4 jobs run "echo isolated" --path /tmp/custom-jobs-test > /dev/null && sleep 1 && aux4 jobs list --path /tmp/custom-jobs-test | jq '. | length'
```

```expect
1
```

### should not pollute default jobs dir when using custom path

```execute
ls /tmp/custom-jobs-test/1/ | sort && rm -rf /tmp/custom-jobs-test
```

```expect
job.json
stderr
stdout
```

## remove-all

### should remove all completed jobs by source

```execute
aux4 jobs run "echo r1" --source remove-batch > /dev/null && aux4 jobs run "echo r2" --source remove-batch > /dev/null && sleep 1 && aux4 jobs remove-all --source remove-batch && aux4 jobs list --source remove-batch
```

```expect:partial
removed
**[]
```
