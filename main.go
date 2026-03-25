package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Job struct {
	ID         string `json:"id"`
	Command    string `json:"command"`
	PID        int    `json:"pid"`
	State      string `json:"state"`
	ExitCode   int    `json:"exitCode"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime,omitempty"`
	Dir        string `json:"dir"`
	OnSuccess  string `json:"onSuccess,omitempty"`
	OnFailure  string `json:"onFailure,omitempty"`
	OnComplete string `json:"onComplete,omitempty"`
}

func baseDir() string {
	return ".jobs"
}

func jobDir(id string) string {
	return filepath.Join(baseDir(), id)
}

func loadJob(id string) (*Job, error) {
	data, err := os.ReadFile(filepath.Join(jobDir(id), "job.json"))
	if err != nil {
		return nil, fmt.Errorf("job %s not found", id)
	}
	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func saveJob(job *Job) error {
	dir := jobDir(job.ID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "job.json"), data, 0644)
}

func allocateID() string {
	base := baseDir()
	os.MkdirAll(base, 0755)
	for i := 1; ; i++ {
		id := strconv.Itoa(i)
		if err := os.Mkdir(filepath.Join(base, id), 0755); err == nil {
			return id
		}
	}
}

func refreshState(job *Job) {
	if job.State == "RUNNING" && job.PID > 0 && !isProcessRunning(job.PID) {
		job.State = "FAILED"
		job.ExitCode = -1
		job.EndTime = time.Now().UTC().Format(time.RFC3339)
		saveJob(job)
	}
}

func formatDuration(job *Job) string {
	start, err := time.Parse(time.RFC3339, job.StartTime)
	if err != nil {
		return "-"
	}
	var d time.Duration
	if job.EndTime != "" {
		end, err := time.Parse(time.RFC3339, job.EndTime)
		if err != nil {
			return "-"
		}
		d = end.Sub(start)
	} else {
		d = time.Since(start)
	}
	return d.Truncate(time.Millisecond).String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: aux4-jobs <command> [args...]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "run":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: aux4-jobs run <command> [onSuccess] [onFailure] [onComplete]")
			os.Exit(1)
		}
		onSuccess := ""
		onFailure := ""
		onComplete := ""
		if len(args) >= 2 {
			onSuccess = args[1]
		}
		if len(args) >= 3 {
			onFailure = args[2]
		}
		if len(args) >= 4 {
			onComplete = args[3]
		}
		runJob(args[0], onSuccess, onFailure, onComplete)
	case "_monitor":
		if len(args) < 2 {
			os.Exit(1)
		}
		monitorJob(args[0], args[1])
	case "list":
		state := ""
		if len(args) >= 1 {
			state = args[0]
		}
		listJobs(state)
	case "status":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: aux4-jobs status <id>")
			os.Exit(1)
		}
		statusJob(args[0])
	case "output":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: aux4-jobs output <id> [stdout|stderr]")
			os.Exit(1)
		}
		stream := "stdout"
		if len(args) >= 2 {
			stream = args[1]
		}
		outputJob(args[0], stream)
	case "tail":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: aux4-jobs tail <id> [stdout|stderr]")
			os.Exit(1)
		}
		stream := "stdout"
		if len(args) >= 2 {
			stream = args[1]
		}
		tailJob(args[0], stream)
	case "kill":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: aux4-jobs kill <id>")
			os.Exit(1)
		}
		killJob(args[0])
	case "killall":
		killAllJobs()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runJob(command, onSuccess, onFailure, onComplete string) {
	id := allocateID()
	cwd, _ := os.Getwd()

	job := &Job{
		ID:         id,
		Command:    command,
		State:      "RUNNING",
		StartTime:  time.Now().UTC().Format(time.RFC3339),
		Dir:        cwd,
		OnSuccess:  onSuccess,
		OnFailure:  onFailure,
		OnComplete: onComplete,
	}
	if err := saveJob(job); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	monitor := exec.Command(self, "_monitor", id, command)
	monitor.Dir = cwd
	monitor.Stdin = nil
	monitor.Stdout = nil
	monitor.Stderr = nil
	setSysProcAttr(monitor)

	if err := monitor.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	monitor.Process.Release()

	data, _ := json.Marshal(job)
	fmt.Println(string(data))
}

func monitorJob(id, command string) {
	job, err := loadJob(id)
	if err != nil {
		os.Exit(1)
	}

	dir := jobDir(id)

	stdoutFile, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		os.Exit(1)
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(filepath.Join(dir, "stderr"))
	if err != nil {
		os.Exit(1)
	}
	defer stderrFile.Close()

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Dir = job.Dir
	setSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		job.State = "FAILED"
		job.ExitCode = 1
		now := time.Now().UTC().Format(time.RFC3339)
		job.EndTime = now
		saveJob(job)
		fmt.Fprintln(stderrFile, err.Error())
		os.Exit(1)
	}

	job.PID = cmd.Process.Pid
	saveJob(job)

	cmdErr := cmd.Wait()

	job, _ = loadJob(id)
	if job.State == "KILLED" {
		executeCallback(job, job.OnComplete)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	job.EndTime = now

	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			job.ExitCode = exitErr.ExitCode()
		} else {
			job.ExitCode = 1
		}
		job.State = "FAILED"
	} else {
		job.ExitCode = 0
		job.State = "COMPLETED"
	}
	saveJob(job)

	if job.State == "COMPLETED" {
		executeCallback(job, job.OnSuccess)
	} else if job.State == "FAILED" {
		executeCallback(job, job.OnFailure)
	}
	executeCallback(job, job.OnComplete)
}

func executeCallback(job *Job, callback string) {
	if callback == "" {
		return
	}
	cmd := exec.Command("sh", "-c", callback)
	cmd.Dir = job.Dir
	cmd.Env = append(os.Environ(),
		"AUX4_JOB_ID="+job.ID,
		"AUX4_JOB_STATE="+job.State,
		"AUX4_JOB_EXIT_CODE="+strconv.Itoa(job.ExitCode),
		"AUX4_JOB_COMMAND="+job.Command,
		"AUX4_JOB_DIR="+job.Dir,
	)
	stderrPath := filepath.Join(jobDir(job.ID), "stderr")
	if f, err := os.OpenFile(stderrPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		cmd.Stderr = f
		defer f.Close()
	}
	cmd.Run()
}

func loadAllJobs() []*Job {
	base := baseDir()
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil
	}

	type entry struct {
		id  int
		job *Job
	}
	var jobs []entry

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		job, err := loadJob(e.Name())
		if err != nil {
			continue
		}
		refreshState(job)
		id, _ := strconv.Atoi(job.ID)
		jobs = append(jobs, entry{id, job})
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].id < jobs[j].id
	})

	var result []*Job
	for _, e := range jobs {
		result = append(result, e.job)
	}
	return result
}

func listJobs(state string) {
	allJobs := loadAllJobs()

	var result []Job
	for _, job := range allJobs {
		if state == "" || job.State == state {
			result = append(result, *job)
		}
	}
	if result == nil {
		result = []Job{}
	}

	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}

func statusJob(id string) {
	job, err := loadJob(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	refreshState(job)

	type Status struct {
		ID        string `json:"id"`
		Command   string `json:"command"`
		PID       int    `json:"pid"`
		State     string `json:"state"`
		ExitCode  int    `json:"exitCode"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime,omitempty"`
		Duration  string `json:"duration"`
		Dir       string `json:"dir"`
	}

	result := Status{
		ID:        job.ID,
		Command:   job.Command,
		PID:       job.PID,
		State:     job.State,
		ExitCode:  job.ExitCode,
		StartTime: job.StartTime,
		EndTime:   job.EndTime,
		Duration:  formatDuration(job),
		Dir:       job.Dir,
	}

	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}

func outputJob(id, stream string) {
	if stream != "stdout" && stream != "stderr" {
		fmt.Fprintln(os.Stderr, "Error: stream must be stdout or stderr")
		os.Exit(1)
	}

	path := filepath.Join(jobDir(id), stream)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(data)
}

func tailJob(id, stream string) {
	if stream != "stdout" && stream != "stderr" {
		fmt.Fprintln(os.Stderr, "Error: stream must be stdout or stderr")
		os.Exit(1)
	}

	path := filepath.Join(jobDir(id), stream)
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
		}
		if readErr == io.EOF {
			job, jobErr := loadJob(id)
			if jobErr != nil || job.State != "RUNNING" {
				io.Copy(os.Stdout, f)
				break
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if readErr != nil {
			break
		}
	}
}

func killJob(id string) {
	job, err := loadJob(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if job.State != "RUNNING" {
		fmt.Fprintf(os.Stderr, "Error: job %s is not running (state: %s)\n", id, job.State)
		os.Exit(1)
	}

	if job.PID > 0 {
		if err := killProcess(job.PID); err != nil {
			p, findErr := os.FindProcess(job.PID)
			if findErr == nil {
				p.Kill()
			}
		}
	}

	job.State = "KILLED"
	job.ExitCode = -1
	job.EndTime = time.Now().UTC().Format(time.RFC3339)
	saveJob(job)

	fmt.Printf("job %s killed\n", id)
}

func killAllJobs() {
	allJobs := loadAllJobs()

	killed := 0
	for _, job := range allJobs {
		if job.State != "RUNNING" {
			continue
		}
		if job.PID > 0 {
			if err := killProcess(job.PID); err != nil {
				p, findErr := os.FindProcess(job.PID)
				if findErr == nil {
					p.Kill()
				}
			}
		}
		job.State = "KILLED"
		job.ExitCode = -1
		job.EndTime = time.Now().UTC().Format(time.RFC3339)
		saveJob(job)
		fmt.Printf("job %s killed\n", job.ID)
		killed++
	}

	if killed == 0 {
		fmt.Println("no running jobs")
	}
}
