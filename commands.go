package main

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"encoding/json"
	log "libchorizo/log"
	"github.com/streadway/amqp"
)
type ParseCommand struct {
	Hash    string
	Args    []string
	Command string
	Channel amqp.Channel
}

type CommandResp interface {
	Response() (string, error)
}

type PingResp struct {
	Hash 						string
	Command 					string
	ReturnString 				string
}


func NewParseCommand(input_string string) (*ParseCommand, error){
	pc := &ParseCommand{}
	input_string_byte := []byte(input_string)
	err := json.Unmarshal(input_string_byte, &pc)
	return pc, err
} 

func (pc *ParseCommand) ExecuteCommand() (CommandResp, error){
	switch pc.Command {
		case "ping":
			pr := &PingResp{}
			pr.Hash = pc.Hash
			pr.Command = "pong"
			pr.ReturnString = "Here is the return string"
			return pr, nil
			break
		default:
			return nil, nil
			break
	}
	return nil, nil

}

func (pr *PingResp) Response() (string, error){
	jsonResp, err := json.Marshal(pr)
	return string(jsonResp), err
}

func (pc *ParseCommand) Response() (string, error){
	return "asdf", nil
}

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
