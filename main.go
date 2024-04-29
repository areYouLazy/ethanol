package main

import (
	"github.com/sirupsen/logrus"

	"github.com/areYouLazy/ethanol/config"
	"github.com/areYouLazy/ethanol/core"
	"github.com/areYouLazy/ethanol/flags"
	"github.com/areYouLazy/ethanol/log"
)

func init() {
	// parse flags
	flags.Init()

	// init log subsystem
	log.Init()

	// load configuration
	config.Init()
}

func main() {
	// greetings
	logrus.Info("welcome to ethanol")

	// generate and init a new core instance
	c := core.NewCore()
	c.Init()
}
