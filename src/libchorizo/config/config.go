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
		ExecPath   string
		ConfigPath string
		RabbitmqHost string
		RabbitmqPort string
		RabbitmqUser string
		RabbitmqPass string
		UseTls bool
	}
}

func (cfg Config) NewConfig(exec_path string) (Config){
	cfg.Main.ExecPath = exec_path
	cfg.Main.ConfigPath = fmt.Sprintf("%s/chorizo.gcfg", cfg.Main.ExecPath)
	err := gcfg.ReadFileInto(&cfg, cfg.Main.ConfigPath)
	cfg.Main.Dbfile = fmt.Sprintf("%s/%s", exec_path, cfg.Main.Dbfile)
	cfg.Main.Cronfile = fmt.Sprintf("%s/%s", exec_path, cfg.Main.Cronfile)
	cfg.Main.Scriptpath = fmt.Sprintf("%s/%s", exec_path, cfg.Main.Scriptpath)
	cfg.Main.Statefile = fmt.Sprintf("%s/%s", exec_path, cfg.Main.Statefile)
	if err != nil {
		panic(err)
	}
	return cfg

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
