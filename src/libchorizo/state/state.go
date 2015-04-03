package libchorizo

import (
	"database/sql"
	log "libchorizo/log"
)

type State struct {
	Id                    int
	Update_guid           string
	Last_script_completed string
	Finished              int
}

// Create inserts into the database a new state entry
// takes the guid as an argument and sets completed = 0 (false)
func (state *State) Create(db *sql.DB) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("insert into state(update_guid) values(?)")
	defer stmt.Close()
	_, err = stmt.Exec(state.Update_guid)
	stmt.Close()
	tx.Commit()
	return true, err
}

// SetLastScriptCompleted updates the state entry to the last successfully completed script that ran
// this is used to track state after a reboot so the update process can continue
func (state *State) SetLastScriptCompleted(db *sql.DB, last_script string) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("update state set last_script_completed = ? where update_guid = ?")
	defer stmt.Close()
	_, err = stmt.Exec(last_script, state.Update_guid)
	stmt.Close()
	tx.Commit()
	return true, err
}

// Finish sets the finished column to 1 marking the update process as completed
func (state *State) Finish(db *sql.DB) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("update state set finished = 1 where update_guid = ?")
	defer stmt.Close()
	_, err = stmt.Exec(state.Update_guid)
	stmt.Close()
	tx.Commit()
	return true, err
}

// GetStateByGUID returns a State struct as found by it's guid
func GetStateByGUID(db *sql.DB, guid string) (State, error) {
	rows, err := db.Query("select id, update_guid, last_script_completed, finished from state where update_guid = ?", guid)
	var state State
	for rows.Next() {
		rows.Scan(&state.Id, &state.Update_guid, &state.Last_script_completed, &state.Finished)
	}
	rows.Close()
	return state, err

}

// GetStateByGUID returns a State struct as found by it's guid
func GetMostRecentState(db *sql.DB) (State, error) {
	s_log := log.GetLogger()
	defer func() {
		if e := recover(); e != nil {
			s_log.Error(e)
		}
	}()
	var state State
	err := db.QueryRow("select id, update_guid, last_script_completed, finished from state order by id desc limit 1").Scan(&state.Id, &state.Update_guid, &state.Last_script_completed, &state.Finished)
	return state, err
}
