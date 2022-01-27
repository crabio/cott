package usecase

import (
	"math/rand"
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

	dtuc.testTable(tcra, r)

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

func (dtuc *databaseTesterUsecase) testTable(tcra *domain.TestCaseResultsAccumulator, r repository.DatabaseTesterRepository) {
	var (
		keyValueTableName   = "key_value"
		keyValueTableFields = []string{
			"f1 BIGINT",
			"f2 BIGSERIAL",
			"f3 BOOLEAN",
			"f4 DATE",
			"f5 FLOAT",
			"f6 REAL",
			"f7 INTEGER",
			"f8 NUMERIC",
			"f9 SMALLINT",
			"f10 SMALLSERIAL",
			"f11 SERIAL",
		}
		keyValueTableColumns = []string{
			"f1",
			"f2",
			"f3",
			"f4",
			"f5",
			"f6",
			"f7",
			"f8",
			"f9",
			"f10",
			"f11",
		}
	)

	if err := dtuc.calcStepDuration(func() error { return r.CreateTable(keyValueTableName, keyValueTableFields) }, "createTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.SingleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1))
	}, "singleInsertTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1)); err != nil {
				return err
			}
		}
		return nil
	}, "10xsingleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1)); err != nil {
				return err
			}
		}
		return nil
	}, "100xsingleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.SingleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1)); err != nil {
				return err
			}
		}
		return nil
	}, "1000xsingleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		if err := r.MultipleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(10)); err != nil {
			return err
		}
		return nil
	}, "10xmultipleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		if err := r.MultipleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(100)); err != nil {
			return err
		}
		return nil
	}, "100xmultipleInsertKeyValueTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		if err := r.MultipleInsert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
			return err
		}
		return nil
	}, "1000xmultipleInsertKeyValueTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.DropTable(keyValueTableName) }, "dropKeyValueTable", tcra)
}

// Method geerates data set for:
/*
keyValueTableFields = []string{
	"f1 BIGINT",
	"f2 BIGSERIAL",
	"f3 BOOLEAN",
	"f4 DATE",
	"f5 FLOAT",
	"f6 REAL",
	"f7 INTEGER",
	"f8 NUMERIC",
	"f9 SMALLINT",
	"f10 SMALLSERIAL",
	"f11 SERIAL",
}
*/
func (dtuc *databaseTesterUsecase) generateTableData(count int) []interface{} {
	var buf []interface{}

	for i := 0; i < count; i++ {
		// "f1 BIGINT",
		buf = append(buf, rand.Int())
		// "f2 BIGSERIAL",
		buf = append(buf, rand.Uint64())
		// "f3 BOOLEAN",
		buf = append(buf, true)
		// "f4 DATE",
		buf = append(buf, time.Now())
		// "f5 FLOAT",
		buf = append(buf, rand.Float32())
		// "f6 REAL",
		buf = append(buf, rand.Float64())
		// "f7 INTEGER",
		buf = append(buf, rand.Int())
		// "f8 NUMERIC",
		buf = append(buf, rand.Int())
		// "f9 SMALLINT",
		buf = append(buf, rand.Intn(255))
		// "f10 SMALLSERIAL",
		buf = append(buf, rand.Intn(255))
		// "f11 SERIAL",
		buf = append(buf, rand.Uint32())
	}

	return buf
}
