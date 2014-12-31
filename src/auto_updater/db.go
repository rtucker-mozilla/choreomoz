package auto_updater

import (
	"database/sql"
	"fmt"
	"time"
	"util"
)

type LogObject struct {
	id          int
	synced      int
	stdout      string
	stderr      string
	log_time    time.Time
	return_code int
	update_guid string
	script      string
}

type UpdateIdToUpdateGuid struct {
	id          int
	update_id   int
	update_guid string
}

// CreateDbIfNotExists creates a sqlite database at db_file path if it doesn't exist.
// returns a bool of the success of creating the database file.
func CreateDbIfNotExists(db_file string) bool {
	if !util.FileExists(db_file) {
		db, err := sql.Open("sqlite3", db_file)
		defer db.Close()
		if err != nil {
			panic(err)
		}
		sqlStmt := `create table system_logs (
			id integer not null primary key,
			stdout text,
			stderr text,
			script text,
			log_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			synced int DEFAULT 0,
			update_guid text,
			return_code int);
			create table version (version int);
			insert into version (version) values (1);
			create table state (
				id integer not null primary key,
				update_guid text,
				last_script_completed text,
				finished int default 0
			);
			create table update_id_to_update_guid (
				id integer not null primary key,
				update_id int,
				update_guid text
			);
			`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			panic(err)
			return false
		}
		db.Close()
	}
	return true
}

// GetUpdateId returns the mapping of the local update to the centalized API mapping
// If a mapping does not exist, return 0
func GetUpdateId(db *sql.DB, update_guid string) int {
	var return_id = 0
	rows, err := db.Query("select id, update_id, update_guid from update_id_to_update_guid where update_guid = ?", update_guid)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var update_mapping UpdateIdToUpdateGuid
		rows.Scan(&update_mapping.id, &update_mapping.update_id, &update_mapping.update_guid)
		return_id = update_mapping.update_id
	}
	rows.Close()
	return return_id

}

// CreateIdToGUIDMapping inserts into the local database a mapping of the remote id to the local update_guid
func CreateIdToGUIDMapping(db *sql.DB, update_id int, update_guid string) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare("insert into update_id_to_update_guid(update_id, update_guid) values(?, ?)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec(update_id, update_guid)
	stmt.Close()
	tx.Commit()
}

// SetToSynced updates the entry in system_logs to 1 meaning that it has been synced
func SetToSynced(db *sql.DB, system_log_id_slice []int) bool {
	for i := range system_log_id_slice {
		tx, err := db.Begin()
		if err != nil {
			return false
		}
		stmt, err := tx.Prepare("update system_logs set synced = 1 where id = ?")
		system_log_id := system_log_id_slice[i]
		if err != nil {
			panic(err)
		}
		ret, err := stmt.Exec(system_log_id)
		fmt.Println("ret:", ret)
		if err != nil {
			panic(err)
		}
		tx.Commit()
		stmt.Close()
	}
	return true
}

// DBPoll is a separate go routine that queries the database for logs to push to the centralized API
func DBPoll(db *sql.DB, HOSTNAME string, API_URL string, system_id int) {
	log := GetLogger()
	for {
		should_create := false
		current_update_guid := ""
		if system_id == 0 {
			rest_system_id, rest_err := GetSystemId(API_URL, HOSTNAME)
			if rest_err != nil {
				fmt.Sprintf("couldn't contact api to get id")
			}
			system_id = rest_system_id
		}

		// If system_id is still 0, there's no reason to iterate over the recorded events because the API is down or unavailable
		if system_id != 0 {
			system_log_id_slice := []int{}
			active_rows := []*LogObject{}
			rows, err := db.Query("select id, synced, log_time, return_code, stdout, stderr, update_guid, script from system_logs where synced = 0")
			if err != nil {
				panic(err)
			}
			for rows.Next() {
				var lo LogObject
				rows.Scan(&lo.id, &lo.synced, &lo.log_time, &lo.return_code, &lo.stdout, &lo.stderr, &lo.update_guid, &lo.script)
				active_rows = append(active_rows, &lo)
			}
			rows.Close()
			for i := range active_rows {
				lo := active_rows[i]
				var system_update_id = GetUpdateId(db, lo.update_guid)
				state, _ := GetStateByGUID(db, lo.update_guid)
				log.Error("state:", state)
				log.Error("system_update_id:", system_update_id)
				log.Error("current_update_guid:", current_update_guid)
				if lo.update_guid != current_update_guid {
					should_create = true
				} else {
					should_create = false
				}
				if system_update_id == 0 || should_create == true {
					log.Error("Starting System Update")
					system_update_id_tmp, _ := CreateSytemUpdate(API_URL, system_id)
					system_update_id = system_update_id_tmp
					CreateIdToGUIDMapping(db, system_update_id, lo.update_guid)
					should_create = false
					current_update_guid = lo.update_guid
				}
				log.Error("Logging System Script Output")
				is_logged := LogCapture(API_URL, system_id, system_update_id, lo)
				if is_logged == true {
					system_log_id_slice = append(system_log_id_slice, lo.id)
					current_update_guid = lo.update_guid
				}
				SetToSynced(db, system_log_id_slice)
				if state.Finished == 1 && state.Last_script_completed == lo.script {
					log.Error("Setting to finished")
					FinishSystemUpdate(API_URL, system_id)
				}
				system_log_id_slice = []int{}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
