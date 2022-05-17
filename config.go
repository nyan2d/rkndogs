package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Address string `yaml:"addr"`
	Mask    int    `json:"mask"`
}

func ReadConfig(path string) ([]Config, error) {
	var obj []Config

	buf, err := os.ReadFile(path)
	if err != nil {
		return obj, err
	}

	err = yaml.Unmarshal(buf, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}
