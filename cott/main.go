package main

import (
	"encoding/json"

	"github.com/iakrevetkho/components-tests/cott/config"
	"github.com/iakrevetkho/components-tests/cott/internal/helpers"

	tester_usecase "github.com/iakrevetkho/components-tests/cott/tester/usecase"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

var shutdownCh chan bool
var cfg config.Config

func init() {
	shutdownCh = helpers.AwaitProcSignals()

	if err := configor.Load(&cfg, "config.yaml"); err != nil {
		logrus.WithError(err).Fatal("Can't parse conf")
	}

	if err := helpers.InitLogger(&cfg); err != nil {
		logrus.WithError(err).Fatal("Couldn't init logger")
	}

	if cfgJson, err := json.Marshal(cfg); err != nil {
		logrus.WithError(err).Fatal("Couldn't serialize config to JSON")
	} else {
		// Use Infof to prevent \" symbols if using WithField
		logrus.Infof("Loaded config: %s", cfgJson)
	}
}

func main() {
	tester_usecase.NewTesterUsecase()

	logrus.Info("Awaiting signal.")
	<-shutdownCh
	logrus.Info("Exit.")
}
