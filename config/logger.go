package config

import (
	"github.com/apsdehal/go-logger"
	"os"
)

var Log *logger.Logger

func init() {
	var err error
	Log, err = logger.New("codenight", 3, os.Stdout)
	Log.Debug("Logger configured")
	if err != nil {
		panic(err)
	}
}
