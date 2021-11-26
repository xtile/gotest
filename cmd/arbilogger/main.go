package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/xtile/gotest/internal/app/arbilogger"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/arbilogger.toml", "path to config file")
}

func main() {

	flag.Parse()

	config := arbilogger.NewConfig()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := arbilogger.New(config)

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}

}
