package main

import (
	"flag"
	"log"

	"github.com/nyan2d/rkndogs/app"
	"github.com/nyan2d/rkndogs/offline"
)

func main() {
	var address string
	var confpath string

	flag.StringVar(&address, "net", "", "network address")
	flag.StringVar(&confpath, "config", "config.yaml", "config path")
	flag.Parse()

	if address != "" {
		a := app.NewApp()
		a.Listen(address)
	} else {
		log.Println(confpath)
		offline.Do(confpath)
	}
}
