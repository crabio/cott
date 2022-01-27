package usecase

import (
	"math/rand"
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

const (
	DATABASE_NAME = "cott_db"
)

type DatabaseTesterUsecase interface {
	RunCase(tcra *domain.TestCaseResultsAccumulator) error
}

type databaseTesterUsecase struct {
	databaseName string
}

func NewDatabaseTesterUsecase() DatabaseTesterUsecase {
	dtuc := new(databaseTesterUsecase)
	dtuc.databaseName = DATABASE_NAME
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
		logrus.WithError(err).WithField("name", name).Warn("error on step execution")
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
		keyValueTableName   = "test_table"
		keyValueTableFields = []string{
			"id BIGSERIAL PRIMARY KEY",
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

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncateEmptyTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1))
	}, "1xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate1XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(10))
	}, "10xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate10XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(100))
	}, "100xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate100XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000))
	}, "1000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate1000XTable", tcra)

	// Use for loop because postgress 65536 params
	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate10000XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate100000XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate1000000XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate10000000XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate100000000XTable", tcra)

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "1000000000xInsertEmptyTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.TruncateTable(keyValueTableName) }, "truncate1000000000XTable", tcra)

	// Single inserts into table with 1000000 rows
	for i := 0; i < 1000; i++ {
		if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
			tcra.AddError(err.Error())
			return
		}
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1))
	}, "1xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(10))
	}, "10xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(100))
	}, "100xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000))
	}, "1000xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsert1000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsert1000000xTable", tcra); err != nil {
		return
	}

	// Single inserts into table with 10000000 rows
	for i := 0; i < 10000; i++ {
		if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
			tcra.AddError(err.Error())
			return
		}
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1))
	}, "1xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(10))
	}, "10xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(100))
	}, "100xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000))
	}, "1000xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsert10000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000000xInsert10000000xTable", tcra); err != nil {
		return
	}

	// Single inserts into table with 100000000 rows
	for i := 0; i < 100000; i++ {
		if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
			tcra.AddError(err.Error())
			return
		}
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1))
	}, "1xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(10))
	}, "10xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(100))
	}, "100xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		return r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000))
	}, "1000xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 1000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "1000000xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 10000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "10000000xInsert100000000xTable", tcra); err != nil {
		return
	}

	if err := dtuc.calcStepDuration(func() error {
		for i := 0; i < 100000; i++ {
			if err := r.Insert(keyValueTableName, keyValueTableColumns, dtuc.generateTableData(1000)); err != nil {
				return err
			}
		}
		return nil
	}, "100000000xInsert100000000xTable", tcra); err != nil {
		return
	}

	dtuc.calcStepDuration(func() error { return r.DropTable(keyValueTableName) }, "dropTable", tcra)
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
func (dtuc *databaseTesterUsecase) generateTableData(count int) []map[string]interface{} {
	var values []map[string]interface{}

	for i := 0; i < count; i++ {
		valuesSet := make(map[string]interface{})

		// "f1 BIGINT",
		valuesSet["f1"] = rand.Intn(255)
		// "f2 BIGSERIAL",
		valuesSet["f2"] = rand.Intn(255)
		// "f3 BOOLEAN",
		valuesSet["f3"] = rand.Intn(255) > 128
		// "f4 DATE",
		valuesSet["f4"] = time.Now()
		// "f5 FLOAT",
		valuesSet["f5"] = rand.Float32()
		// "f6 REAL",
		valuesSet["f6"] = rand.Float64()
		// "f7 INTEGER",
		valuesSet["f7"] = rand.Intn(255)
		// "f8 NUMERIC",
		valuesSet["f8"] = rand.Intn(255)
		// "f9 SMALLINT",
		valuesSet["f9"] = rand.Intn(255)
		// "f10 SMALLSERIAL",
		valuesSet["f10"] = rand.Intn(255)
		// "f11 SERIAL",
		valuesSet["f11"] = rand.Intn(255)

		values = append(values, valuesSet)
	}

	return values
}
