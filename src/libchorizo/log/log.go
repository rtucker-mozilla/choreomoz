package libchorizo

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"parse_update_script"
	"strings"
	"time"
	"libchorizo/config"
)

func GetLogger() *logrus.Logger {
	var CONFIGPATH = "/etc/chorizo/chorizo.gcfg"
	config, _ := libchorizo.ParseConfig(CONFIGPATH)
	var level_string = config.Main.Loglevel
	log := logrus.New()
	switch level_string {
	case "DEBUG":
		log.Level = logrus.DebugLevel
	case "INFO":
		log.Level = logrus.InfoLevel
	case "WARNING":
		log.Level = logrus.WarnLevel
	case "ERROR":
		log.Level = logrus.ErrorLevel
	}
	return log
}

func FixLogLine(line string) string {
	line = strings.Replace(line, "\t", "    ", -1)
	line = strings.Trim(line, " ")
	return line
}

func WriteLogLine(logline string, logline_type string, ret_code int, filepath string) {
	now := time.Now()
	time_string := now.Format("20060102150405")
	if logline != "" {
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}
		logline_split := strings.Split(logline, "\n")
		for i := 0; i < len(logline_split); i++ {
			logline := logline_split[i]
			stdout := ""
			stderr := ""
			if logline_type == "stdout" {
				stdout = logline
			} else {
				stderr = logline
			}
			write_string := fmt.Sprintf("%s\t%d\t%s\t%s\n", time_string, ret_code, stdout, stderr)
			f.WriteString(write_string)
		}
		f.Close()
	}

}
func ProcessLog(update_guid string, ret_code int, stdout string, stderr string, us *parse_update_script.UpdateScript, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("insert into system_logs(update_guid, stdout, stderr, return_code, script) values(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	stmt.Exec(update_guid, stdout, stderr, ret_code, us.FilePath)
	stmt.Close()
	tx.Commit()
}
