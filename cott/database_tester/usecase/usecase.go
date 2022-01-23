package usecase

import (
	"time"

	"github.com/iakrevetkho/components-tests/cott/database_tester/repository"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

const DATABASE_NAME = "cott_db"

type DatabaseTesterUsecase interface {
	RunCase(tc *domain.TestCase) (*domain.Report, error)
}

type databaseTesterUsecase struct {
}

func NewDatabaseTesterUsecase() DatabaseTesterUsecase {
	return &databaseTesterUsecase{}
}

func (tuc *databaseTesterUsecase) RunCase(tc *domain.TestCase) (*domain.Report, error) {
	switch tc.ComponentType {
	case domain.ComponentType_Postgres:
		r := repository.NewPostgresDatabaseTesterRepository(tc.Port, tc.Host, tc.User, tc.Password)

		report := domain.NewReport(tc)

		start := time.Now()
		if err := r.Open(); err != nil {
			return nil, err
		}
		duration := time.Since(start)
		logrus.WithField("duration", duration).Debug("open connection")
		report.AddMetric("open connection", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))

		start = time.Now()
		if err := r.CreateDatabase(DATABASE_NAME); err != nil {
			return nil, err
		}
		duration = time.Since(start)
		logrus.WithField("duration", duration).Debug("create database")
		report.AddMetric("create database", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))

		start = time.Now()
		if err := r.DropDatabase(DATABASE_NAME); err != nil {
			return nil, err
		}
		duration = time.Since(start)
		logrus.WithField("duration", duration).Debug("drop database")
		report.AddMetric("drop database", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))

		start = time.Now()
		if err := r.Close(); err != nil {
			return nil, err
		}
		duration = time.Since(start)
		logrus.WithField("duration", duration).Debug("close connection")
		report.AddMetric("close connection", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))

		return report, nil

	// Open connection speed

	// Create database speed
	// Create tables speed
	// Single insert speed
	// Multiple insert speed
	// Random select
	// Single insert foreign key
	// Multiple foreign insert speed
	// Random foreign select
	// Join speed
	// Drop table speed
	// Drop database speed

	default:
		return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
	}
}
