package libchorizo

import (
	"time"
)

type LogObject struct {
	Id          int
	Synced      int
	Stdout      string
	Stderr      string
	Log_time    time.Time
	Return_code int
	Update_guid string
	Script      string
}