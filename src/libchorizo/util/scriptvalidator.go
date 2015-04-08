package libchorizo

import (
	"os"
	"regexp"
	"strings"
)

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
