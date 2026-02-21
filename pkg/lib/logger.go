package lib

import (
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
)

func SetupLogger(debug bool) log.Logger {
	var logger log.Logger

	{
		logger = log.NewLogfmtLogger(os.Stderr)

		if debug {
			logger = level.NewFilter(logger, level.AllowDebug())
		} else {
			logger = level.NewFilter(logger, level.AllowInfo())
		}

		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	return logger
}
