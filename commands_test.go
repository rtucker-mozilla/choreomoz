package main

// To run this test
// go test util_test.go util.go

import (
	"testing"
)
var ping_example = "{\"hash\": \"dde93f95d664df0c518e10bff196d9111e30e7ad\", \"args\": [\"arg1\", \"arg2\"], \"command\": \"ping\"}"

func TestParseCommandCommand(t *testing.T) {
	pc, _ := NewParseCommand(ping_example)
	if pc.Command != "ping" {
		t.Error("NewParseCommand not setting Command to ping correctly")
	}
}

func TestParseCommandArgs(t *testing.T) {
	pc, _ := NewParseCommand(ping_example)
	if len(pc.Args) != 2 {
		t.Error("NewParseCommand doesn't have proper length")
	}
	if pc.Args[0] != "arg1" {
		t.Error("NewParseCommand doesn't have proper first arg")
	}
	if pc.Args[1] != "arg2" {
		t.Error("NewParseCommand doesn't have proper second arg")
	}
}

func TestParseCommandHash(t *testing.T) {
	var proper_hash = "dde93f95d664df0c518e10bff196d9111e30e7ad"
	pc, _ := NewParseCommand(ping_example)
	if pc.Hash != proper_hash {
		t.Error("NewParseCommand doesn't set proper hash. Expected", proper_hash, "but got: ", pc.Hash )
	}
}

func TestParseCommandPing(t *testing.T) {
	pc, _ := NewParseCommand(ping_example)
	ret, err := pc.Ping()
	if err != nil {
		t.Error("err is not nil")
	}
	if ret != "asdf" {
		t.Error("NewParseCommand doesn't set proper hash")
	}
}
func TestParseCommandExecuteCommand(t *testing.T) {
	pc, _ := NewParseCommand(ping_example)
	ret, err := pc.ExecuteCommand()
	if err != nil {
		t.Error("err is not nil")
	}
	if ret != "ping" {
		t.Error("ExecuteCommand doesn't return proper function")
	}
}
