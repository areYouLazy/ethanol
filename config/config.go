package config

import (
	"github.com/areYouLazy/ethanol/flags"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	// set filename
	viper.SetConfigFile(flags.Config)

	// log
	logrus.WithFields(logrus.Fields{
		"file_name": viper.GetViper().ConfigFileUsed(),
	}).Info("reading configuration file")

	// read configuration file
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("fatal error reading configuration file")
	}

	// enable watch config, for now just for monitoring
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithFields(logrus.Fields{
			"file_name":     e.Name,
			"event_message": e.String(),
		}).Debug("configuration file changed")
	})
}
