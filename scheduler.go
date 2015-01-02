package main

import (
	"auto_updater"
	"cron_eval"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/exec"
	"parse_update_script"
	"path/filepath"
	"strconv"
	"time"
	"util"
)

// SystemReboot executes a shell command to reboot the host
func SystemReboot() {
	cmd := exec.Command("/sbin/shutdown", "-r", "now")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// GUIDHash returns a SHA1 hash of the hostname and timestamp to tag local logging groups
func GUIDHash(hostname string) string {
	current_time := time.Now()
	hasher := sha1.New()
	write_hash_string := []byte(fmt.Sprintf("%s%d", hostname, current_time.UnixNano()))
	hasher.Write(write_hash_string)
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha
}

// Main entry point
func main() {
	log := auto_updater.GetLogger()

	exec_path, _ := os.Getwd()
	HOSTNAME, _ := os.Hostname()
	config, config_err := auto_updater.ParseConfig()
	if config_err != nil {
		log.Error("Unable to open config file")
	}
	log.Debug("Config Loaded")
	log.Debug(HOSTNAME)
	var DB_FILE = fmt.Sprintf("%s/%s", exec_path, config.Main.Dbfile)
	log.Debug("DB_FILE: ", DB_FILE)
	var CRONFILE = fmt.Sprintf("%s/%s", exec_path, config.Main.Cronfile)
	var STATEFILE = fmt.Sprintf("%s/%s", exec_path, config.Main.Statefile)
	var LOCKFILE = fmt.Sprintf("%s/%s", exec_path, config.Main.Lockfile)
	var SCRIPTPATH = fmt.Sprintf("%s/%s", exec_path, config.Main.Scriptpath)
	var APIURL = config.Main.APIUrl
	log.Debug("GUIDHash: ", GUIDHash(HOSTNAME))
	db_created := auto_updater.CreateDbIfNotExists(DB_FILE)
	if db_created {
		log.Info("DB Created at path: ", DB_FILE)
	}
	//	db, db_open_err := sql.Open("sqlite3", DB_FILE)
	db1, db_open_err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", DB_FILE))
	db2, db_open_err2 := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", DB_FILE))
	if db_open_err != nil {
		panic("Unable to open existing database")
	}
	if db_open_err2 != nil {
		panic("Unable to open existing database")
	}
	state, _ := auto_updater.GetMostRecentState(db1)
	log.Error(state)
	var currently_locked = util.HasLockFile(LOCKFILE)
	if currently_locked == true {
		if util.HasStateFile(STATEFILE) {
			log.Debug("There exists a statefile")
		}
	}
	if !util.HasScriptPath(SCRIPTPATH) {
		log.Error(fmt.Sprintf("Script Path %s does not exist.", SCRIPTPATH))
		os.Exit(2)
	}
	go auto_updater.DBPoll(db1, HOSTNAME, APIURL, 0)
	for {
		util.DeleteLockFile(LOCKFILE)
		cron_line, cron_err := util.ReadCronFile(CRONFILE)
		if cron_err != nil {
			log.Debug(fmt.Sprintf("%s", cron_err))
			os.Exit(2)
		}

		run_now, run_after, sleep_seconds := cron_eval.EvalCronLine(cron_line)
		if run_now == false && run_after == false {
			time.Sleep(time.Duration(sleep_seconds) * time.Second)
			continue
		}

		if run_now == false && run_after == true {
			time.Sleep(time.Duration(sleep_seconds) * time.Second)
		}
		update_guid := GUIDHash(HOSTNAME)
		state := auto_updater.State{}
		state.Update_guid = update_guid
		state.Create(db1)

		scripts, _ := filepath.Glob(fmt.Sprintf("%s/*", SCRIPTPATH))
		UpdateScripts := []parse_update_script.UpdateScript{}
		for i := 0; i < len(scripts); i++ {
			script_path := scripts[i]
			var uf parse_update_script.UpdateScriptFile
			var us parse_update_script.UpdateScript
			uf.FilePath = script_path
			parse_update_script.ReadFile(&uf)
			us.ParseScript(&uf)
			UpdateScripts = append(UpdateScripts, us)
		}
		for i := 0; i < len(UpdateScripts); i++ {
			exec_script := UpdateScripts[i].FilePath
			ret_code, stdout, stderr := auto_updater.ExecCommand(exec_script)
			log.Debug("\n===== START LOG CAPTURE OF ExecCommand =====")
			log.Debug("ret_code: ", ret_code)
			log.Debug("stdout: ", stdout)
			log.Debug("stderr: ", stderr)
			log.Debug("===== END LOG CAPTURE OF ExecCommand =====\n")
			auto_updater.ProcessLog(update_guid, ret_code, stdout, stderr, &UpdateScripts[i], db2)
			exit_code_to_i, exit_code_to_i_err := strconv.Atoi(UpdateScripts[i].ScriptExitCodeReboot)
			state.SetLastScriptCompleted(db1, exec_script)

			if exit_code_to_i_err != nil {
				panic(exit_code_to_i_err)
			}
			if exit_code_to_i == ret_code {
				SystemReboot()
			}

		}
		state.Finish(db1)
		time.Sleep(5 * time.Second)
	}
}
