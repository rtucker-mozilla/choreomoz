package main

import (
	"github.com/gorhill/cronexpr"
	"math"
	"fmt"
	"time"
	log "libchorizo/log"
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

func CalcSleepSeconds(in_seconds int) int {
	ret_seconds := in_seconds
	if in_seconds <= 3 {
		ret_seconds = 3
	}
	return ret_seconds
}

func EvalCronLine(cron_line string) (bool, bool, int) {
	log_s := log.GetLogger()
	now := time.Now()
	current_minute := GetCurrentMinute(now)
	next_cron_minute := current_minute.Add(-1 * time.Minute)
	next_cron_run := GetNextCronRun(cron_line, next_cron_minute)
	var run_now, run_after = false, false
	var sleep_seconds = 0
	if current_minute == next_cron_run {
		run_now = true
		sleep_seconds = 0
	} else {
		time_diff := time.Since(next_cron_run)
		sleep_seconds = int(math.Abs(time_diff.Seconds()))

		log_s.Debug(fmt.Sprintf("Sleep Seconds: %s", sleep_seconds))
		log_s.Debug(fmt.Sprintf("current_minute: %s", current_minute))
		log_s.Debug(fmt.Sprintf("next_cron_minute: %s", next_cron_minute))
		if sleep_seconds <= 86400 {
			run_after = true
		}
		sleep_seconds = CalcSleepSeconds(sleep_seconds)
	}
	return run_now, run_after, sleep_seconds
}
