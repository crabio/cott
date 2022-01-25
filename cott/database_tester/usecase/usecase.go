package usecase

import (
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

const (
	POSTGRES_USER_ENV_VAR     = "POSTGRES_USER"
	POSTGRES_PASSWORD_ENV_VAR = "POSTGRES_PASSWORD"
)

type DatabaseTesterUsecase interface {
	RunCase(tcra *domain.TestCaseResultsAccumulator) error
}

type databaseTesterUsecase struct {
	databaseName string
}

func NewDatabaseTesterUsecase(databaseName string) DatabaseTesterUsecase {
	dtuc := new(databaseTesterUsecase)
	dtuc.databaseName = databaseName
	return dtuc
}

func (dtuc *databaseTesterUsecase) RunCase(tcra *domain.TestCaseResultsAccumulator) error {
	r, err := dtuc.createDatabaseRepository(tcra.TestCase)
	if err != nil {
		return err
	}

	const tableName = "test_table"

	if err := dtuc.calcStepDuration(func() error { return r.Open() }, "openConnection", tcra); err != nil {
		return nil
	}

	// Await for DB ready
	if err := dtuc.calcStepDuration(func() error {
		// await 30 second
		for i := 0; i < 300; i++ {
			if err := r.Ping(); err != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				// Success
				return nil
			}
		}
		return domain.CONNECTION_WAS_NOT_ESTABLISHED
	}, "startUp", tcra); err != nil {
		logrus.WithError(err).Debug("couldn't ping database")
		time.Sleep(time.Second)
	}

	if err := r.DropDatabase(dtuc.databaseName); err != nil {
		logrus.WithError(err).Debug("couldn't drop database")
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateDatabase(dtuc.databaseName) }, "createDatabase", tcra); err != nil {
		return nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.SwitchDatabase(dtuc.databaseName) }, "switchDatabase", tcra); err != nil {
		return nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateTable(tableName) }, "createTable", tcra); err != nil {
		return nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.DropTable(tableName) }, "dropTable", tcra); err != nil {
		return nil
	}

	if err := r.SwitchDatabase(""); err != nil {
		tcra.AddError(err.Error())
		return nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.DropDatabase(dtuc.databaseName) }, "dropDatabase", tcra); err != nil {
		return nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.Close() }, "closeConnection", tcra); err != nil {
		return nil
	}

	return nil

	// Create tables speed
	// Single insert speed
	// Multiple insert speed
	// Random select
	// Single insert foreign key
	// Multiple foreign insert speed
	// Random foreign select
	// Join speed
	// Drop table speed

}

func (dtuc *databaseTesterUsecase) createDatabaseRepository(tc *domain.TestCase) (repository.DatabaseTesterRepository, error) {
	switch tc.ComponentType {

	case domain.ComponentType_Postgres:
		// Get user from env vars
		user, ok := tc.EnvVars[POSTGRES_USER_ENV_VAR]
		if !ok {
			logrus.WithField("envVarName", POSTGRES_USER_ENV_VAR).Error(domain.NO_REQUIRED_ENV_VAR_KEY)
			return nil, domain.NO_REQUIRED_ENV_VAR_KEY
		}
		// Get password from env vars
		password, ok := tc.EnvVars[POSTGRES_PASSWORD_ENV_VAR]
		if !ok {
			logrus.WithField("envVarName", POSTGRES_PASSWORD_ENV_VAR).Error(domain.NO_REQUIRED_ENV_VAR_KEY)
			return nil, domain.NO_REQUIRED_ENV_VAR_KEY
		}

		return repository.NewPostgresDatabaseTesterRepository(tc.Port, "localhost", user, password), nil

	default:
		return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
	}
}

func (dtuc *databaseTesterUsecase) calcStepDuration(f func() error, name string, tcra *domain.TestCaseResultsAccumulator) error {
	start := time.Now()
	if err := f(); err != nil {
		tcra.AddError(name + ". " + err.Error())
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "name": name}).Debug("step finished")
	tcra.AddMetric(name+"Duration", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}
