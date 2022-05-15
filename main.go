package main

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Entry struct {
	Address string `yaml:"addr"`
	Mask    string `yaml:"mask"`
}

func main() {
}

func ReadEntries(path string) []Entry {
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var obj []Entry
	err = yaml.Unmarshal(buf, &obj)
	if err != nil {
		log.Fatal(err)
	}

	return obj
}

func IsIP(address string) bool {
	rg := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`)
	return rg.MatchString(address)
}
