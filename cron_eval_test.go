package main
// To run this test
// go test cron_eval_test.go cron_eval.go

import (
	"testing"
	"fmt"
	"time"
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

func TestCalcSleepSecondsInput1(t *testing.T){
	r_sec := CalcSleepSeconds(1)
	if r_sec != 3 {
		t.Error("Sleep seconds should return 3 if input is 1.")
	}
}

func TestCalcSleepSecondsInput3(t *testing.T){
	r_sec := CalcSleepSeconds(3)
	if r_sec != 3 {
		t.Error("Sleep seconds should return 3 if input is 3.")
	}
}

func TestCalcSleepSecondsInput4(t *testing.T){
	r_sec := CalcSleepSeconds(4)
	if r_sec != 4 {
		t.Error("Sleep seconds should return 4 if input is 4.")
	}
}

func TestCalcSleepSecondsInput86400(t *testing.T){
	r_sec := CalcSleepSeconds(86400)
	if r_sec != 86400 {
		t.Error("Sleep seconds should return 86400 if input is 86400.")
	}
}
func TestGetCurrentMinute(t *testing.T){
	time_s := time.Now()
	current_minute := GetCurrentMinute(time_s)
	if current_minute.Second() != 0 {
		t.Error(fmt.Sprintf("Current Seconds should be 0 but is %s", current_minute.Second()))

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
