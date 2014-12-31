package auto_updater

import (
	"code.google.com/p/gcfg"
	"fmt"
	"os"
)

type Config struct {
	Main struct {
		Dbfile     string
		Debug      bool
		Cronfile   string
		Statefile  string
		Lockfile   string
		Scriptpath string
		APIUrl     string
		Loglevel   string
	}
}

func ParseConfig() (Config, error) {
	exec_path, _ := os.Getwd()
	var CONFIGPATH = fmt.Sprintf("%s/updater.gcfg", exec_path)
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, CONFIGPATH)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
