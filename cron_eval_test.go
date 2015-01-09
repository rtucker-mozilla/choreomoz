package main

import (
	"testing"
	//"io/ioutil"
	"fmt"
	//"time"
)

func TestCronRunEveryMinuteEveryDay(t *testing.T) {
	fmt.Println("Starting tests of cron_eval")
	//current_minute := time.Now()
	run_now, _, sleep_seconds := EvalCronLine("* * * * *")
	if run_now != true {
		t.Error("Incorrect compare time for * * * * *")
	}
	if sleep_seconds != 0 {
		t.Error("Incorrect sleep seconds for * * * * *")
	}
}
func TestCronRun59thMinuteOfTheHour(t *testing.T) {
	//This will only pass if it is not the 59th minute of the hour
	//current_minute := time.Now()
	// @TODO refactor to allow mocking a time.Time object to inject to EvalCronLine
	run_now, _, sleep_seconds := EvalCronLine("59 * * * *")
	if run_now != false {
		t.Error("Incorrect compare time for 59 * * * *")
	}
	if sleep_seconds == 0 {
		fmt.Println(fmt.Sprintf("Return value is %d", sleep_seconds))
		t.Error("Incorrect sleep seconds for 00 * * * *")
	}
}
