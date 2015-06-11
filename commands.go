package main

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"fmt"
	"os"
	"encoding/json"
	log "libchorizo/log"
	"github.com/streadway/amqp"
)
type ParseCommand struct {
	Hash    					string
	Hostname    				string
	Args    					[]string
	Command 					string
	GroupId 					int
	Channel 					amqp.Channel
}

type CommandResp interface {
	Response() (string, error)
}

type PingResp struct {
	Hash 						string
	Hostname    				string
	Command 					string
	ReturnString 				string
	GroupId 					int
}

type StartRebootResp struct {
	Hash 						string
	Hostname    				string
	Command 					string
	ReturnString 				string
	GroupId 					int
}

type StartUpdateResp struct {
	Hash 						string
	Hostname    				string
	Command 					string
	GroupId 					int
	ReturnString 				string
	Args						[]string
}

type ExecResp struct {
	Hash 						string
	Hostname    				string
	Command 					string
	ReturnString 				string
	Args						[]string
	ExitCode 					int
	StdOut   					string
	StdErr   					string
	GroupId 					int
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
			pr.Hostname, _ = os.Hostname()
			pr.Command = "pong"
			pr.ReturnString = "Here is the return string"
			return pr, nil
			break
		case "start_update":
			sur := &StartUpdateResp{}
			sur.Hash = pc.Hash
			sur.Hostname, _ = os.Hostname()
			sur.Command = "start_update_resp"
			sur.ReturnString = "Starting Update"
			sur.Args = pc.Args
			fmt.Println("start_update: pc.GroupId: ", pc.GroupId)
			sur.GroupId = pc.GroupId
			return sur, nil
			break
		case "exec":
			execr := &ExecResp{}
			execr.Hash = pc.Hash
			execr.Hostname, _ = os.Hostname()
			execr.Command = "exec_response"
			execr.Args = pc.Args
			execr.GroupId = pc.GroupId
			return execr, nil
			break
		case "start_reboot":
			start_reboot_response := &StartRebootResp{}
			start_reboot_response.Hash = pc.Hash
			start_reboot_response.Hostname, _ = os.Hostname()
			start_reboot_response.Command = "start_reboot_resp"
			start_reboot_response.GroupId = pc.GroupId
			return start_reboot_response, nil
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

func (pr *StartRebootResp) Response() (string, error){
	jsonResp, err := json.Marshal(pr)
	return string(jsonResp), err
}

func (execr *ExecResp) Response() (string, error){
	action_string := strings.Join(execr.Args," ")
	exit_code, outStr, errStr := ExecCommand(action_string)
	execr.ExitCode = exit_code
	execr.StdOut = outStr
	execr.StdErr = errStr
	jsonResp, err := json.Marshal(execr)

	return string(jsonResp), err
}

func (sur *StartUpdateResp) Response() (string, error){
	jsonResp, err := json.Marshal(sur)
	return string(jsonResp), err
}

func (pc *ParseCommand) Response() (string, error){
	return "", nil
}

// ExecCommand executes the provided command string.
// It returns the exit_code, stdout, stderr.
func ExecCommand(cmd_string string) (int, string, string) {
	log := log.GetLogger()
	log.Error("cmd_string:", cmd_string)
	exec_string := fmt.Sprintf("/etc/chorizo/scripts/%s", cmd_string)
	var exit_code = 0
	cmd := exec.Command("sh", "-c", exec_string)
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
