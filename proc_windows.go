//go:build windows

package main

import (
	"fmt"
	"os/exec"
)

func setSysProcAttr(cmd *exec.Cmd) {}

func isProcessRunning(pid int) bool {
	return true
}

func killProcess(pid int) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid)).Run()
}
