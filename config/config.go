package config

import (
	"path"
	"strings"

	"github.com/areYouLazy/ethanol/flags"
	"github.com/areYouLazy/ethanol/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	// check if we have a custom config file
	if flags.ConfigFile != "" {
		// check if we have a custom path to parse
		if strings.Contains(flags.ConfigFile, "/") {
			// split path (absolute or relative) and filename
			logrus.WithFields(logrus.Fields{
				"file_name": flags.ConfigFile,
			}).Debug("found path in filename provided by cli flag ", utils.Bold("-config-file"))
			fpath, fname := path.Split(flags.ConfigFile)
			viper.SetConfigFile(fname)
			viper.AddConfigPath(fpath)
		} else {
			// we have only a filename
			logrus.WithFields(logrus.Fields{
				"file_name": flags.ConfigFile,
			}).Debug("no path in filename provided by cli flag ", utils.Bold("-config-file"))
			viper.SetConfigName(flags.ConfigFile)
			viper.AddConfigPath(".")
		}
	} else {
		// fallback to default
		logrus.WithFields(logrus.Fields{
			"file_name": "config.yml",
		}).Debug("using default configuration file")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	// set configuration type
	if flags.ConfigJSON {
		logrus.Debug("parsing configuration file as json because of flag ", utils.Bold("-config-json"))
		viper.SetConfigType("json")
	} else {
		logrus.Debug("parsing configuration file as yaml")
		viper.SetConfigType("yaml")
	}

	// logrus.Debug("parsing configuration file as yaml")
	// viper.SetConfigType("yaml")

	logrus.WithFields(logrus.Fields{
		"file_name": viper.GetViper().ConfigFileUsed(),
	}).Debug("reading configuration file")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("fatal error reading configuration file")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithFields(logrus.Fields{
			"file_name": e.Name,
		}).Debug("configuration file changed")
	})
}
