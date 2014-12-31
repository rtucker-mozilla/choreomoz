package cron_eval
import (
	"fmt"
	"time"
	"math"
	"github.com/gorhill/cronexpr"
)

func GetCurrentMinute(now time.Time) time.Time {
	fixed_now := time.Date(
		now.Year(),
		now.Month(),
		int(now.Day()),
		now.Hour(),
		now.Minute(),
		0,
		0,
		now.Location())
	return fixed_now

}

func GetNextCronRun(cron_line string, current_minute time.Time) time.Time {
	nextTime := cronexpr.MustParse(cron_line).Next(current_minute)
	return nextTime
}

func EvalCronLine(cron_line string) (bool, bool, int){
	now := time.Now()
	var DEBUG = true
	current_minute := GetCurrentMinute(now)
	next_cron_minute := current_minute.Add(-1 * time.Minute)
	next_cron_run := GetNextCronRun(cron_line, next_cron_minute)
	var run_now = false
	var run_after = false
    var sleep_seconds = 0
	if current_minute == next_cron_run {
		run_now = true
		sleep_seconds = 0
	} else {
		time_diff := time.Since(next_cron_run)
		sleep_seconds = int(math.Abs(time_diff.Seconds()))
		if DEBUG {
			fmt.Println("Sleep Seconds: ", sleep_seconds)
			fmt.Println("current_minute: ", current_minute)
			fmt.Println("next_cron_minute: ", next_cron_minute)
		}
		if sleep_seconds <= 86400 {
			run_after = true
		}
		if sleep_seconds <= 3 {
			sleep_seconds = 3
		}
	}
	return run_now, run_after, sleep_seconds
}