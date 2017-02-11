package config

import (
	"github.com/apsdehal/go-logger"
	"os"
)

var Log *logger.Logger

func SetupLogger() {
	var err error
	Log, err = logger.New("codenight", 3, os.Stdout)
	if err != nil {
		panic(err)
	}
}
