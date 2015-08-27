package config

import (
	"encoding/json"
	"github.com/syhlion/gopusher/module/log"
	"os"
)

type Config struct {
	AuthAccount  string `json:"auth_account"`
	AuthPassword string `json:"auth_password"`
	Environment  string `json:"environment"`
	LogDir       string `json:"logdir"`
	Listen       string `json:"listen"`
	SqlDir       string `json:"sqldir"`
}

func GetConfig(configDir string) *Config {
	file, err := os.OpenFile(configDir, os.O_RDONLY, 0655)
	defer file.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Logger.Fatal(err)
	}

	return &config
}
