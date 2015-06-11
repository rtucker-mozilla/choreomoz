package libchorizo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

func ScriptValid(script_path string) bool {
	if !FileExists(script_path) {
		return false
	}
	sv := NewScriptValidator(script_path)
	return sv.IsExecutable && sv.ValidFileName
}

func HostnameToQueueName(input_string string) string {
	queueName := strings.Replace(input_string, ".", "-", -1)
	return queueName
}
func HostnameToBindingKey(input_string string) string {
	keyName := strings.Replace(input_string, ".", "-", -1)
	keyName = fmt.Sprintf("%s.host", keyName)
	return keyName
}

