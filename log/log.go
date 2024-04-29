package log

import (
	"github.com/areYouLazy/ethanol/flags"
	"github.com/sirupsen/logrus"
)

func Init() {
	if flags.LogLevelDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if flags.LogCaller {
		logrus.SetReportCaller(true)
	}

	if flags.LogFormatJSON {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
