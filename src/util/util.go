package util
import (
	"fmt"
	"io/ioutil"
	"os"
	"errors"
)

func FileExists(filename string) bool{
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

func WriteLockFile(filename string){
	ioutil.WriteFile(filename, []byte("a"), 0600)
}

func DeleteLockFile(filename string){
	if HasLockFile(filename) == true {
		os.Remove(filename)
	}
}

func DeleteStateFile(filename string){
	if HasStateFile(filename) == true {
		os.Remove(filename)
	}
}
