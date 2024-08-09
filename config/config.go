package config

import (
	"path"

	"github.com/areYouLazy/ethanol/flags"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	// check if we have a custom config file
	if flags.ConfigFile != "" {
		// check if there's a relative or absolute path in filename
		fpath, fname := path.Split(flags.ConfigFile)

		// path.Split() returns empty fpath if there are no slashes in the provided file path
		if fpath != "" {
			// we found a path
			viper.AddConfigPath(fpath)
		}

		// add default location and custom file name
		viper.AddConfigPath(".")
		viper.SetConfigFile(fname)

		// log
		logrus.WithFields(logrus.Fields{
			"filename": flags.ConfigFile,
		}).Info("custom configuration file provided")
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	// set configuration type, default to yaml
	viper.SetConfigType("yaml")

	// set configuration type as json if requested by flag
	if flags.ConfigJSON {
		viper.SetConfigType("json")
		logrus.Debug("reading configuration as a json file")
	}

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
