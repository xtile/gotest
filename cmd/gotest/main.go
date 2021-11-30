package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/xtile/gotest/internal/app/gotest"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/arbilogger.toml", "path to config file")
}

func main() {

	flag.Parse()

	config := gotest.NewConfig()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := gotest.New(config)

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}

	log.Fatal("Finishing app...")
}
