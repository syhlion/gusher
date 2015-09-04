package config

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	AuthAccount  string `json:"auth_account"`
	AuthPassword string `json:"auth_password"`
	Environment  string `json:"environment"`
	LogFile      string `json:"logfile"`
	Listen       string `json:"listen"`
	SqlFile      string `json:"sqlfile"`
	MaxWaitHook  int    `json:"max_wait_hook"`
}

func GetConfig(configfile string) *Config {
	file, err := os.OpenFile(configfile, os.O_RDONLY, 0655)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func Write(config Config) (err error) {
	os.Remove("./config.json")
	file, err := os.OpenFile("./config.json", os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		panic(err)
		return
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	return

}
