package config

import (
	"os"
	"gopkg.in/yaml.v2"
	"log"
)

type config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Email struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}

func InitConfig() {
	f, err := os.Open("../src/config.yml")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal(err)
	}
}

var Config config


