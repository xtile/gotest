package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/xtile/gotest/internal/app/arbi"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/arbilogger.toml", "path to config file")
}

func main() {

	flag.Parse()

	config := arbi.NewConfig()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := arbi.New(config)

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}

}
