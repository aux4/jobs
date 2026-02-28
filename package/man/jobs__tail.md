#### Description

The `tail` command streams job output in real time. If the job is still running, it continuously polls for new output and prints it as it appears. When the job finishes, it prints any remaining output and exits.

If the job has already completed, `tail` prints the full output and exits immediately, behaving the same as `aux4 jobs output`.

By default it tails stdout. Use `--stream stderr` to tail standard error instead.

#### Usage

```bash
aux4 jobs tail <id> [--stream <stdout|stderr>]
```

--id      The job ID (positional argument)
--stream  Output stream to tail (default: `stdout`, options: `stdout`, `stderr`)

#### Example

```bash
aux4 jobs tail 1
```

```text
Building...
Compiling main.go
Compiling utils.go
Build complete.
```
