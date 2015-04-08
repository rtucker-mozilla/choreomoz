package libchorizo

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

func FileStat(script_path string) (os.FileInfo, error) {
	info, info_err := os.Stat(script_path)
	return info, info_err
}

type ScriptValidator struct {
	Filepath string
	Filemode string
	Fileinfo os.FileInfo
	IsExecutable bool
	ValidFileName bool
}

func NewScriptValidator(filepath string) *ScriptValidator {
	sv := new(ScriptValidator)
	sv.Filepath = filepath
	sv.ValidFileName = false
	sv.getfileinfo()
	sv.getfilemode()
	sv.getisexecutable()
	sv.getvalidfilename()
	return sv
}

func (sv *ScriptValidator) getfileinfo() (error) {
	info, info_err := os.Stat(sv.Filepath)
	sv.Fileinfo = info
	return info_err
}

func (sv *ScriptValidator) getfilemode() {
	file_mode := sv.Fileinfo.Mode()
	file_mode_str := file_mode.String()
	sv.Filemode = file_mode_str
}

func (sv *ScriptValidator) getisexecutable() {
	sv.IsExecutable = false
	regex_mode_match, regex_mode_match_err := regexp.MatchString("x", sv.Filemode)
	if regex_mode_match && regex_mode_match_err == nil {
		sv.IsExecutable = true
	}
}

func (sv *ScriptValidator) getvalidfilename() {
	short_path_slice := strings.Split(sv.Filepath, "/")
	short_path := short_path_slice[len(short_path_slice)-1]
	regex_mode_match, regex_mode_match_err := regexp.MatchString("^\\d+.*", short_path)
	if regex_mode_match && regex_mode_match_err == nil {
		sv.ValidFileName = true
	}
}

func ScriptValid(script_path string) (bool) {
	if !FileExists(script_path){
		return false
	}
	sv := NewScriptValidator(script_path)
	return sv.IsExecutable && sv.ValidFileName
}
