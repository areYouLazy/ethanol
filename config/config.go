package config

import (
	"path"

	"github.com/areYouLazy/ethanol/flags"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	// check if there's a relative or absolute path in filename
	fpath, fname := path.Split(flags.Config)

	// path.Split() returns empty fpath if there are no slashes in the provided file path
	if fpath == "" {
		// add current path if omitted
		fpath = "."
	}

	// add configuration path
	viper.AddConfigPath(fpath)

	// set filename
	viper.SetConfigFile(fname)

	// set configuration type, default to yaml
	viper.SetConfigType("yaml")

	// log
	logrus.WithFields(logrus.Fields{
		"file_name": viper.GetViper().ConfigFileUsed(),
		"file_path": fpath,
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
