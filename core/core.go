package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/areYouLazy/ethanol/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Core struct{}

func NewCore() *Core {
	// allocate a new core instance
	c := &Core{}

	// return new core instance
	return c
}

func (c *Core) Init() {
	var wait time.Duration

	// proxy
	proxy.Setup()

	// define server instance
	srv := &http.Server{
		Addr:         viper.GetString("Ethanol.Server.Address") + ":" + viper.GetString("Ethanol.Server.Port"),
		WriteTimeout: time.Duration(viper.GetInt("Ethanol.Server.WriteTimeout")) * time.Second,
		ReadTimeout:  time.Duration(viper.GetInt("Ethanol.Server.ReadTimeout")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("Ethanol.Server.IdleTimeout")) * time.Second,
		Handler:      getCoreEngine(),
	}

	// use goroutines to start services
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err.Error())
		}
	}()

	// create a channel of lenght 1 to intercept signals
	signalChannel := make(chan os.Signal, 1)

	// listen for SIGINT
	signal.Notify(signalChannel, os.Interrupt)

	// greetings
	logrus.WithFields(logrus.Fields{
		"server_address": srv.Addr,
	}).Info("ready to serve")

	// debug configuration
	logrus.WithFields(logrus.Fields{
		"server_write_timeout": srv.WriteTimeout,
		"server_read_timeout":  srv.ReadTimeout,
		"server_idle_timeout":  srv.IdleTimeout,
	}).Debug("web server configuration")

	// wait for events on signal channel
	<-signalChannel

	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// shutdown services with context
	srv.Shutdown(ctx)

	// greetings
	logrus.Info("bye")
}
