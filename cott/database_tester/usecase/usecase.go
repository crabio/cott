package usecase

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
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

	dtuc.testKeyValueTable(tcra, r)

	dtuc.testMeasurmentsTable(tcra, r)

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
		const (
			POSTGRES_USER_ENV_VAR     = "POSTGRES_USER"
			POSTGRES_PASSWORD_ENV_VAR = "POSTGRES_PASSWORD"
		)

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

func (dtuc *databaseTesterUsecase) testKeyValueTable(tcra *domain.TestCaseResultsAccumulator, r repository.DatabaseTesterRepository) {
	var (
		keyValueTableName    = "key_value"
		keyValueTableFields  = []string{"key VARCHAR(255)", "value VARCHAR(255)"}
		keyValueTableColumns = []string{"key", "value"}
	)

	if err := dtuc.calcStepDuration(func() error { return r.CreateTable(keyValueTableName, keyValueTableFields) }, "createKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"})
	}, "singleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"}); err != nil {
				return err
			}
		}
		return nil
	}, "100xInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"}); err != nil {
				return err
			}
		}
		return nil
	}, "1000xInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10000; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"}); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100000; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"}); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000000; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, []interface{}{"key", "value"}); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsertKeyValueTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.DropTable(keyValueTableName) }, "dropKeyValueTable", tcra)
}

func (dtuc *databaseTesterUsecase) testMeasurmentsTable(tcra *domain.TestCaseResultsAccumulator, r repository.DatabaseTesterRepository) {
	var (
		measurmentsTableName        = "measurement"
		measurmentsTableFieldsCount = 100
		measurmentsTableFields      = []string{"id SERIAL PRIMARY KEY"}
		measurmentsTableColumns     = []string{}
		measurmentsTableData        = []interface{}{}
	)
	for i := 0; i < measurmentsTableFieldsCount; i++ {
		measurmentsTableFields = append(measurmentsTableFields, "s"+strconv.FormatInt(int64(i), 10)+" REAL")
		measurmentsTableColumns = append(measurmentsTableColumns, "s"+strconv.FormatInt(int64(i), 10))
		measurmentsTableData = append(measurmentsTableData, rand.Float64())
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateTable(measurmentsTableName, measurmentsTableFields) }, "createMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData)
	}, "singleInsertMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData); err != nil {
				return err
			}
		}
		return nil
	}, "100xInsertMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData); err != nil {
				return err
			}
		}
		return nil
	}, "1000xInsertMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10000; i++ {
			if err := r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsertMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100000; i++ {
			if err := r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsertMeasurementsTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000000; i++ {
			if err := r.SingleInsert(measurmentsTableName, measurmentsTableColumns, measurmentsTableData); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsertMeasurementsTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.DropTable(measurmentsTableName) }, "dropMeasurementsTable", tcra)
}
