package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
	memSol     string
)

type Config struct {
	BindAddr           string `toml:"bind_addr"`
	LogLevel           string `toml:"log_level"`
	Token              string `toml:"token"`
	RequestFreq        int    `toml:"requestfreq"`
	SendMessageTimeout int    `toml:"sendmessagetimeout"`
	Storage            struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Database string `toml:"database"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Attempts int    `toml:"attempts2con"`
	} `toml:"storage"`
}

func parseFlags() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatal("Wrong binary parameters, try -help")
	}
}

func NewConfig() (*Config, error) {
	parseFlags()
	config := &Config{}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
		return config, err
	}
	return config, nil
}
