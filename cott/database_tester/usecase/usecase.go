package usecase

import (
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

type DatabaseTesterUsecase interface {
	RunCase(tc *domain.TestCase) (*domain.TestCaseResults, error)
}

type databaseTesterUsecase struct {
	databaseName string
}

func NewDatabaseTesterUsecase(databaseName string) DatabaseTesterUsecase {
	dtuc := new(databaseTesterUsecase)
	dtuc.databaseName = databaseName
	return dtuc
}

func (dtuc *databaseTesterUsecase) RunCase(tc *domain.TestCase) (*domain.TestCaseResults, error) {
	r, err := dtuc.createDatabaseRepository(tc)
	if err != nil {
		return nil, err
	}

	const tableName = "test_table"

	tcr := domain.NewTestCaseResults(tc)

	if err := dtuc.calcStepDuration(func() error { return r.Open() }, "openConnection", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := r.DropDatabase(dtuc.databaseName); err != nil {
		logrus.WithError(err).Debug("couldn't drop database")
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateDatabase(dtuc.databaseName) }, "createDatabase", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.SwitchDatabase(dtuc.databaseName) }, "switchDatabase", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateTable(tableName) }, "createTable", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.DropTable(tableName) }, "dropTable", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := r.SwitchDatabase(""); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.DropDatabase(dtuc.databaseName) }, "dropDatabase", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	if err := dtuc.calcStepDuration(func() error { return r.Close() }, "closeConnection", tcr); err != nil {
		tcr.AddError(err)
		return tcr, nil
	}

	return tcr, nil

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
		return repository.NewPostgresDatabaseTesterRepository(tc.Port, tc.Host, tc.User, tc.Password, ""), nil

	default:
		return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
	}
}

func (dtuc *databaseTesterUsecase) calcStepDuration(f func() error, name string, tcr *domain.TestCaseResults) error {
	start := time.Now()
	if err := f(); err != nil {
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "name": name}).Debug("step finished")
	tcr.AddMetric(name+"Duration", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}
