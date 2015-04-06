package main

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	log "libchorizo/log"
)

// ExecCommand executes the provided command string.
// It returns the exit_code, stdout, stderr.
func ExecCommand(cmd_string string) (int, string, string) {
	log := log.GetLogger()
	log.Error("cmd_string:", cmd_string)
	var exit_code = 0
	cmd := exec.Command(cmd_string)
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	cmd.Run()
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	exit_code = waitStatus.ExitStatus()
	stdOut_ret := stdOut.String()
	stdOut_ret = strings.Trim(stdOut_ret, " \n")
	stdErr_ret := stdErr.String()
	stdErr_ret = strings.Trim(stdErr_ret, " \n")
	return exit_code, stdOut_ret, stdErr_ret
}
