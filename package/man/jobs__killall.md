#### Description

The `killall` command terminates all currently running jobs. It iterates through all jobs, finds those with state `RUNNING`, and kills each one. Each killed job is printed to stdout. If no jobs are running, it prints `no running jobs`.

#### Usage

```bash
aux4 jobs killall
```

#### Example

```bash
aux4 jobs killall
```

```text
job 1 killed
job 3 killed
```
