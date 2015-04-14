package libchorizo

import (
	"code.google.com/p/gcfg"
	"fmt"
)

type ChorizoConfig interface {
	InterpolateConfig(config_file_path string) (Config)
}

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

type inConfig struct{}

func (ChorizoConfig inConfig) InterpolateConfig(config_file_path string) (Config) {
	var cfg Config
	return cfg
}

func (ChorizoConfig Config) InterpolateConfig(config_file_path string) (Config) {
	var cfg Config
	return cfg
}

func ParseConfig(exec_path string) (Config, error) {
	//exec_path, _ := os.Getwd()
	// In https://github.com/rtucker-mozilla/chorizo/issues/15
	// going to specify a specific config file path
	var CONFIGPATH = fmt.Sprintf("%s/chorizo.gcfg", exec_path)
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, CONFIGPATH)
	cfg.InterpolateConfig(CONFIGPATH)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
