package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// secure is true if ListenAndServeTLS is used
	secure = false

	// scheme is "https" if ListeAndServeTLS is used
	scheme = "http"
)

type Core struct{}

func NewCore() *Core {
	// allocate a new core instance
	c := &Core{}

	// return new core instance
	return c
}

func (c *Core) Init() {
	// tls
	err := checkEthanolTLSCertificateKeyPair()
	if err != nil {
		generateEthanolTLSCertificateKeyPair()
	}

	// get logrus writer to be used in webserver configuration
	w := logrus.New().Writer()
	defer w.Close()

	// define server instance
	srv := &http.Server{
		Addr:         viper.GetString("ethanol.server.address") + ":" + viper.GetString("ethanol.server.port"),
		WriteTimeout: time.Duration(viper.GetInt("ethanol.server.writetimeout")) * time.Second,
		ReadTimeout:  time.Duration(viper.GetInt("ethanol.server.readtimeout")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("ethanol.server.idletimeout")) * time.Second,
		Handler:      getCoreEngine(),
		TLSConfig: &tls.Config{
			Certificates: getEthanolTLSCertificateKeyPairAsSlice(),
		},
		ErrorLog: log.New(w, "", 0),
	}

	// check if we need to server as a TLS server
	secure := viper.GetBool("ethanol.server.tls.enabled")

	// use goroutines to start services
	go func() {
		if secure {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				if err != http.ErrServerClosed {
					logrus.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Error("error in webserver execution")
				} else {
					logrus.Info("webserver closed")
				}
			}
		} else {
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					logrus.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Error("error in webserver execution")
				} else {
					logrus.Info("webserver closed")
				}
			}
		}
	}()

	// create a channel of lenght 1 to intercept signals
	signalChannel := make(chan os.Signal, 1)

	// listen for SIGINT
	signal.Notify(signalChannel, os.Interrupt)

	// adjust scheme according to webserver configuration
	scheme := "http"
	if secure {
		scheme = "https"
	}

	// greetings
	logrus.WithFields(logrus.Fields{
		"server_address": fmt.Sprintf("%s://%s", scheme, srv.Addr),
	}).Info("ready to serve")

	// debug configuration
	logrus.WithFields(logrus.Fields{
		"server_write_timeout": srv.WriteTimeout,
		"server_read_timeout":  srv.ReadTimeout,
		"server_idle_timeout":  srv.IdleTimeout,
	}).Debug("web server configuration")

	// wait for events on signal channel
	<-signalChannel

	// load timeout value from config
	wait := viper.GetDuration("ethanol.server.shutdowntimeout")

	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), wait*time.Second)
	defer cancel()

	// shutdown services with context
	srv.Shutdown(ctx)

	// greetings
	logrus.Info("bye")
}
