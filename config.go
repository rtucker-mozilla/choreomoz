package main

import (
	"code.google.com/p/gcfg"
	"fmt"
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
	//exec_path, _ := os.Getwd()
	// In https://github.com/rtucker-mozilla/chorizo/issues/15
	// going to specify a specific config file path
	exec_path := "/etc/chorizo"
	var CONFIGPATH = fmt.Sprintf("%s/chorizo.gcfg", exec_path)
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, CONFIGPATH)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
