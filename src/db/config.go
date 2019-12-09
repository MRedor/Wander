package db

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Table    string `yaml:"table"`
	} `yaml:"database"`
}

func readConfig() Config {
	f, err := os.Open("./src/config.yml")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
