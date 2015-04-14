package libchorizo

// To run this test
// go test util_test.go util.go

import (
	"testing"
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
)

func RmFile(file_path string) {
	cmd := exec.Command("rm", "-f", file_path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func CreateTestConfig(filepath string){
	config_text := `[main]
dbfile = LOG.db
debug = true
cronfile = CRONFILE
scriptpath = scripts
apiurl = "https://127.0.0.1/api"
statefile = STATEFILE
loglevel = DEBUG`
	d1 := []byte(config_text)
    err := ioutil.WriteFile(filepath, d1, 0644)
	if err != nil {
		panic(err)
	}
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}

func TestNewConfigExecPath(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.ExecPath != exec_path {
		t.Error("NewConfig not setting ExecPath")
	}
}
func TestNewConfigConfigPath(t *testing.T) {
	exec_path := "/tmp"
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.ConfigPath != fmt.Sprintf("%s/chorizo.gcfg", exec_path) {
		t.Error("NewConfig not setting proper ConfigPath")
	}
}
func TestNewConfigDbfile(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.Dbfile != fmt.Sprintf("%s/LOG.db",exec_path) {
		t.Error("NewConfig not setting proper Dbfile")
	}
}
func TestNewConfigCronfile(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.Cronfile != fmt.Sprintf("%s/CRONFILE",exec_path) {
		t.Error("NewConfig not setting proper Cronfile")
	}
}
func TestNewConfigStatefile(t *testing.T) {
	// This has been deprecated, testing only for coverage
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.Statefile != fmt.Sprintf("%s/STATEFILE",exec_path) {
		t.Error("NewConfig not setting proper Statefile")
	}
}
func TestNewConfigScriptPath(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.Scriptpath != fmt.Sprintf("%s/scripts",exec_path) {
		t.Error("NewConfig not setting proper ScriptPath")
	}
}
func TestNewConfigLogLevel(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.Loglevel != "DEBUG" {
		t.Error("NewConfig not setting proper Loglevel")
	}
}
func TestNewConfigAPIUrl(t *testing.T) {
	exec_path := "/tmp"
	full_path := "/tmp/chorizo.gcfg"
	if FileExists(full_path) {
		RmFile(full_path)
	}
	CreateTestConfig(full_path)
	cfg := Config{}
	c := cfg.NewConfig(exec_path)
	if c.Main.APIUrl != "https://127.0.0.1/api" {
		t.Error("NewConfig not setting proper APIUrl")
	}
}
