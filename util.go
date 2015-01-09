package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}

func ReadCronFile(filename string) (string, error) {
	fh, err := ioutil.ReadFile(filename)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to open scheduler file at: %s", filename))
	}
	return string(fh), err
}

func HasScriptPath(filename string) bool {
	return FileExists(filename)
}

func HasLockFile(filename string) bool {
	return FileExists(filename)
}

func HasStateFile(filename string) bool {
	return FileExists(filename)
}

func WriteLockFile(filename string) {
	ioutil.WriteFile(filename, []byte("a"), 0600)
}

func DeleteLockFile(filename string) {
	if HasLockFile(filename) == true {
		os.Remove(filename)
	}
}

func DeleteStateFile(filename string) {
	if HasStateFile(filename) == true {
		os.Remove(filename)
	}
}

// ScriptValid confirms that the script is executable and should be executed
// @TODO: Refactor this so that it is more robust
func ScriptValid(script_path string) (retval bool) {
	retval = false
	info, info_err := os.Stat(script_path)
	if info_err != nil {
		return
	}
	// Get the filemode of the script_path
	// If the mode does not include the execute bit, set retval = false
	file_mode := info.Mode()
	file_mode_str := file_mode.String()
	regex_mode_match, regex_mode_match_err := regexp.MatchString("x", file_mode_str)
	if regex_mode_match && regex_mode_match_err == nil {
		retval = true
	}
	// We have confirmed the script is executable. Now confirm that the script follows the
	// proper naming convention
	if retval == true {
		// Valid filenames are in the form 01_name\.[sh|py|etc]
		short_path_slice := strings.Split(script_path, "/")
		short_path := short_path_slice[len(short_path_slice)-1]
		regex_mode_match, regex_mode_match_err := regexp.MatchString("^\\d+.*", short_path)
		if regex_mode_match && regex_mode_match_err == nil {
			retval = true
		} else {
			retval = false
		}
	}
	return
}
