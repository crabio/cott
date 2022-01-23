package usecase

import (
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

type DatabaseTesterUsecase interface {
	RunCase(tc *domain.TestCase) (*domain.Report, error)
}

type databaseTesterUsecase struct {
	databaseName string
}

func NewDatabaseTesterUsecase(databaseName string) DatabaseTesterUsecase {
	dtuc := new(databaseTesterUsecase)
	dtuc.databaseName = databaseName
	return dtuc
}

func (dtuc *databaseTesterUsecase) RunCase(tc *domain.TestCase) (*domain.Report, error) {
	r, err := dtuc.createDatabaseRepository(tc)
	if err != nil {
		return nil, err
	}

	report := domain.NewReport(tc)

	if err := dtuc.calcStepDuration(func() error { return r.Open() }, "openConnection", report); err != nil {
		return nil, err
	}

	if err := dtuc.calcStepDuration(func() error { return r.CreateDatabase(dtuc.databaseName) }, "createDatabase", report); err != nil {
		return nil, err
	}

	if err := dtuc.calcStepDuration(func() error { return r.DropDatabase(dtuc.databaseName) }, "dropDatabase", report); err != nil {
		return nil, err
	}

	if err := dtuc.calcStepDuration(func() error { return r.Close() }, "closeConnection", report); err != nil {
		return nil, err
	}

	return report, nil

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
		return repository.NewPostgresDatabaseTesterRepository(tc.Port, tc.Host, tc.User, tc.Password), nil

	default:
		return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
	}
}

func (dtuc *databaseTesterUsecase) calcStepDuration(f func() error, name string, report *domain.Report) error {
	start := time.Now()
	if err := f(); err != nil {
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "name": name}).Debug("step finished")
	report.AddMetric(name+"Duration", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}
