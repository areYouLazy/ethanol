package flags

import "flag"

var (
	LogLevelDebug bool
	LogFormatJSON bool
	LogCaller     bool

	ConfigFile string
	ConfigJSON bool
)

func Init() {
	flag.BoolVar(&LogLevelDebug, "debug", false, "print debug log messages (default false)")
	flag.BoolVar(&LogFormatJSON, "json", false, "print log messages in json format (default false)")
	flag.BoolVar(&LogCaller, "caller", false, "print log messages caller (default false)")

	flag.StringVar(&ConfigFile, "config-file", "", "use a custom configuration file")
	flag.BoolVar(&ConfigJSON, "config-json", false, "read configuration file as json")

	flag.Parse()
}
