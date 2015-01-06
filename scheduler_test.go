package main

import (
	"testing"
	//"io/ioutil"
	//"fmt"
	"time"
)

func TestShouldExecReturns(t *testing.T) {
	//This will only pass if it is not the 59th minute of the hour
	//current_minute := time.Now()
	// @TODO refactor to allow mocking a time.Time object to inject to EvalCronLine
	exec_ret := SystemReboot(false)
	if exec_ret != false {
		t.Error("SystemReboot not returning false")
	}
}
func TestGuidHash(t *testing.T) {
	//This will only pass if it is not the 59th minute of the hour
	//current_minute := time.Now()
	// @TODO refactor to allow mocking a time.Time object to inject to EvalCronLine
	the_hash := GUIDHash("test_hostname")
	time.Sleep(10 * time.Millisecond)
	the_hash_iter2 := GUIDHash("test_hostname")
	if the_hash == "" {
		t.Error("GUIDHash returning empty string")
	}
	if len(the_hash) != 40 {
		t.Error("GUIDHash returning incorrect hash length:", len(the_hash))
	}
	if the_hash == the_hash_iter2 {
		t.Error("GUIDHash matches even with additional sleep")
	}
}
