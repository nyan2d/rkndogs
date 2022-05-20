package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Device DeviceConfig  `yaml:"device"`
	Gate   string        `yaml:"gate"`
	Routes []RouteConfig `yaml:"routes"`
}

type DeviceConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RouteConfig struct {
	Address string `yaml:"addr"`
	Mask    int    `yaml:"mask"`
}

func ReadConfig(path string) (Config, error) {
	var obj Config

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
