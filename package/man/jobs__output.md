#### Description

The `output` command prints the full captured output of a job. By default it shows stdout. Use `--stream stderr` to show the standard error output instead.

This command prints the complete output at once. For real-time streaming while a job is running, use `aux4 jobs tail` instead.

#### Usage

```bash
aux4 jobs output <id> [--stream <stdout|stderr>]
```

--id      The job ID (positional argument)
--stream  Output stream to display (default: `stdout`, options: `stdout`, `stderr`)

#### Example

```bash
aux4 jobs output 1
```

```text
hello world
```

```bash
aux4 jobs output 1 --stream stderr
```

```text
warning: deprecated function used
```
