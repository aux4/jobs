#### Description

The `jobs` command group manages background jobs. It allows you to run commands in the background, list all jobs, check their status, retrieve output, tail output in real time, and kill running jobs.

Jobs are stored in `.jobs/` in the current directory. Each job gets a unique numeric ID and stores its stdout and stderr in separate files.

#### Usage

```bash
aux4 jobs <subcommand>
```

Available subcommands: `run`, `list`, `status`, `output`, `tail`, `kill`.

#### Example

```bash
aux4 jobs run "make build"
aux4 jobs list
aux4 jobs status 1
aux4 jobs output 1
aux4 jobs kill 1
```
