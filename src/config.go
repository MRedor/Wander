package main

import (
	"os"
	"log"
	"src/gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Table    string `yaml:"table"`
	} `yaml:"database"`
}

func readConfig() {
	f, err := os.Open("config.yml")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}
