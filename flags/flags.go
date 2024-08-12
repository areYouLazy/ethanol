package flags

import "flag"

var (
	LogLevelDebug bool
	LogFormatJSON bool
	LogCaller     bool

	Config string
)

func Init() {
	flag.BoolVar(&LogLevelDebug, "debug", false, "print debug log messages (default false)")
	flag.BoolVar(&LogFormatJSON, "json", false, "print log messages in json format (default false)")
	flag.BoolVar(&LogCaller, "caller", false, "print log messages caller (default false)")

	flag.StringVar(&Config, "config", "config.yml", "use a custom configuration file")

	flag.Parse()
}
