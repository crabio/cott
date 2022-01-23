package main

import (
	"encoding/json"

	"github.com/iakrevetkho/components-tests/cott/config"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/iakrevetkho/components-tests/cott/internal/helpers"

	database_tester_usecase "github.com/iakrevetkho/components-tests/cott/database_tester/usecase"
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
	dtuc := database_tester_usecase.NewDatabaseTesterUsecase()
	tuc := tester_usecase.NewTesterUsecase(dtuc)

	// Test use case
	tc := &domain.TestCase{
		ComponentType: domain.ComponentType_Postgres,
		Host:          "localhost",
		Port:          5432,
		User:          "user",
		Password:      "password",
	}

	report, err := tuc.RunCase(tc)
	if err != nil {
		logrus.WithError(err).Error("test case error")
	}
	logrus.WithField("report", report).Info("test case done")

	logrus.Info("Awaiting signal.")
	<-shutdownCh
	logrus.Info("Exit.")
}
