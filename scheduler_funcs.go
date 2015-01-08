package main

import (
	"auto_updater"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/jmcvetta/napping"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// GUIDHash returns a SHA1 hash of the hostname and timestamp to tag local logging groups
func GUIDHash(hostname string) string {
	current_time := time.Now()
	hasher := sha1.New()
	write_hash_string := []byte(fmt.Sprintf("%s%d", hostname, current_time.UnixNano()))
	hasher.Write(write_hash_string)
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha
}

// SystemReboot executes a shell command to reboot the host

func GetSystemId(API_URL string, hostname string) (int, error) {
	log := auto_updater.GetLogger()
	// Add function here to pull from local db
	var from_cache = false
	log.Debug("from_cache: ", from_cache)
	rest_system_id, rest_err := auto_updater.GetSystemId(API_URL, hostname)
	// Set value to local db
	return rest_system_id, rest_err
}

func DBStartUpdate(db *sql.DB, update_id int) bool {
	tx, err := db.Begin()
	if err != nil {
		return false
	}
	// insert into state (update_id, finished) values (?, 0)
	stmt, err := tx.Prepare("insert into state (update_id, finished) values (?, 0)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec(update_id)
	stmt.Close()
	tx.Commit()
	return true
}

func DBGetCurrentUpdateId(db *sql.DB) int {
	var update_id = 0
	err := db.QueryRow("select update_id from state order by id desc limit 1").Scan(&update_id)
	if err != nil {
		panic(err)
		update_id = 0
	}
	return update_id
}
func DBEndUpdate(db *sql.DB, update_id int) bool {
	tx, err := db.Begin()
	if err != nil {
		return false
	}
	// insert into state (update_id, finished) values (?, 0)
	stmt, err := tx.Prepare("delete from state")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	return true
}

func DBSetLastCompletedScript(db *sql.DB, update_id int, script string) bool {
	tx, err := db.Begin()
	if err != nil {
		return false
	}
	// insert into state (update_id, finished) values (?, 0)
	stmt, err := tx.Prepare("update state set last_script_completed = ?")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec(script)
	stmt.Close()
	tx.Commit()
	return true
}

func ProcessEntry(us chan UpdateScriptResponse, db *sql.DB) {
	log := auto_updater.GetLogger()
	s := <-us
	log.Debug("\n===== START LOG CAPTURE OF ExecCommand =====")
	log.Debug("ret_code: ", s.ret_code)
	log.Debug("stdout: ", s.stdout)
	log.Debug("stderr: ", s.stderr)
	log.Debug("===== END LOG CAPTURE OF ExecCommand =====\n")
	var update_id = 0
	if s.is_start {
		update_id, _ := auto_updater.CreateSytemUpdate(s.api_url, s.system_id)
		DBStartUpdate(db, update_id)
		log.Error("CreateSystemUpdate")
	}
	if update_id == 0 {
		DBGetCurrentUpdateId(db)
	}
	if s.stderr != "" || s.stdout != "" {
		LogCapture(s.api_url, s.system_id, &s)
		log.Error("Called LogCapture")
	}
	DBSetLastCompletedScript(db, update_id, s.update_script.FilePath)
	if s.is_end {
		auto_updater.FinishSystemUpdate(s.api_url, s.system_id)
		// Here we either delete from the database or mark as completed
		// Delete is probably easier
		DBEndUpdate(db, update_id)
		log.Error("FinishSystemUpdate")
	}
	log.Debug("Leaving ProcessEntry goroutine")
}

func LogCapture(url string, system_id int, log_object *UpdateScriptResponse) bool {
	payload := struct {
		Return_code int    `json:"return_code"`
		Stdout      string `json:"stdout"`
		Stderr      string `json:"stderr"`
		System_id   int    `json:"system_id"`
	}{}
	return_value := false
	var final_url = fmt.Sprintf("%s/logcapture/", url)
	payload.System_id = system_id
	payload.Stdout = log_object.stdout
	payload.Stderr = log_object.stderr
	payload.Return_code = log_object.ret_code
	resp, err := napping.Post(final_url, &payload, nil, nil)
	if resp.Status() == 200 {
		return_value = true
	} else if err != nil {
		return_value = false
	}
	return return_value

}
