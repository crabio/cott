package usecase

import (
	"math/rand"
	"strconv"
	"time"

	container_launcher "github.com/iakrevetkho/components-tests/cott/container_launcher/usecase"
	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	metrics_collector "github.com/iakrevetkho/components-tests/cott/metrics_collector/usecase"
	"github.com/sirupsen/logrus"
)

const (
	DATABASE_NAME = "cott_db"
)

type DatabaseTesterUsecase interface {
	RunCase(tcra *domain.TestCaseResultsAccumulator, containerId string) error
}

type databaseTesterUsecase struct {
	databaseName string
	cluc         container_launcher.ContainerLauncherUsecase
}

func NewDatabaseTesterUsecase(cluc container_launcher.ContainerLauncherUsecase) DatabaseTesterUsecase {
	dtuc := new(databaseTesterUsecase)
	dtuc.databaseName = DATABASE_NAME
	dtuc.cluc = cluc
	return dtuc
}

func (dtuc *databaseTesterUsecase) RunCase(tcra *domain.TestCaseResultsAccumulator, containerId string) error {
	r, err := dtuc.createDatabaseRepository(tcra.TestCase)
	if err != nil {
		return err
	}

	mcuc := metrics_collector.NewMetricsCollectorUsecase(tcra, dtuc.cluc)

	step := &domain.TestCaseStep{Name: "openConnection", StepFunc: func() error { return r.Open() }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return nil
	}

	// Await for DB ready
	step = &domain.TestCaseStep{Name: "startUp", StepFunc: func() error {
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
	}}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		logrus.WithError(err).Debug("couldn't ping database")
		time.Sleep(time.Second)
	}

	if err := r.DropDatabase(dtuc.databaseName); err != nil {
		logrus.WithError(err).Debug("couldn't drop database")
	}

	step = &domain.TestCaseStep{Name: "createDatabase", StepFunc: func() error { return r.CreateDatabase(dtuc.databaseName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return nil
	}

	step = &domain.TestCaseStep{Name: "switchDatabase", StepFunc: func() error { return r.SwitchDatabase(dtuc.databaseName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return nil
	}

	dtuc.testTable(mcuc, r, containerId)

	if err := r.SwitchDatabase(""); err != nil {
		return err
	}

	step = &domain.TestCaseStep{Name: "dropDatabase", StepFunc: func() error { return r.DropDatabase(dtuc.databaseName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return nil
	}

	step = &domain.TestCaseStep{Name: "closeConnection", StepFunc: func() error { return r.Close() }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
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

func (dtuc *databaseTesterUsecase) testTable(mcuc metrics_collector.MetricsCollectorUsecase, r repository.DatabaseTesterRepository, containerId string) {
	var (
		tableName           = "test_table"
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
		tableColumns     = []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11"}
		selectConditions = "f1>1 AND f2>1 AND f3 AND F5>0.5 AND f6>0.5 AND f7>1 AND f8>1 AND f9>1 AND f10>1 AND f11>1"
	)

	step := &domain.TestCaseStep{Name: "createTable", StepFunc: func() error { return r.CreateTable(tableName, keyValueTableFields) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return
	}

	step = &domain.TestCaseStep{Name: "truncateEmptyTable", StepFunc: func() error { return r.TruncateTable(tableName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return
	}

	for i := 1; i <= 10000000; i *= 10 {
		if err := dtuc.testTableInsertSelect(mcuc, r, tableName, tableColumns, selectConditions, i); err != nil {
			return
		}
	}

	step = &domain.TestCaseStep{Name: "dropTable", StepFunc: func() error { return r.DropTable(tableName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return
	}

	step = &domain.TestCaseStep{Name: "dropTable", StepFunc: func() error { return r.DropTable(tableName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return
	}
}

func (dtuc *databaseTesterUsecase) testTableInsertSelect(mcuc metrics_collector.MetricsCollectorUsecase, r repository.DatabaseTesterRepository, tableName string, tableColumns []string, selectConditions string, dataCount int) error {
	testPrefix := strconv.FormatInt(int64(dataCount), 10) + "x"

	step := &domain.TestCaseStep{Name: testPrefix + "InsertEmptyTable", StepFunc: func() error {
		if dataCount > 1000 {
			// Postgres bulk insert support max 65536 params
			// Split insert by 1000 rows
			for i := dataCount / 1000; i > 0; i-- {
				if err := r.Insert(tableName, tableColumns, dtuc.generateTableData(1000)); err != nil {
					return err
				}
			}
		} else {
			return r.Insert(tableName, tableColumns, dtuc.generateTableData(dataCount))
		}

		return nil
	}}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return err
	}

	step = &domain.TestCaseStep{Name: "selectById" + testPrefix + "Table", StepFunc: func() error { return r.SelectById(tableName, dataCount/2) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return err
	}

	step = &domain.TestCaseStep{Name: "selectByConditions" + testPrefix + "Table", StepFunc: func() error { return r.SelectByConditions(tableName, selectConditions) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return err
	}

	// Inserts into full table
	if dataCount >= 1000 {
		for i := 1000; i >= 1; i /= 10 {
			insertTestPrefix := strconv.FormatInt(int64(i), 10) + "x"

			step = &domain.TestCaseStep{Name: insertTestPrefix + "Insert" + testPrefix + "Table", StepFunc: func() error { return r.Insert(tableName, tableColumns, dtuc.generateTableData(i)) }}
			if err := mcuc.CollectStepMetrics(step); err != nil {
				return err
			}
		}
	}

	step = &domain.TestCaseStep{Name: "truncate" + testPrefix + "Table", StepFunc: func() error { return r.TruncateTable(tableName) }}
	if err := mcuc.CollectStepMetrics(step); err != nil {
		return err
	}

	return nil
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
