package usecase

import (
	database_tester_usecase "github.com/iakrevetkho/components-tests/cott/database_tester/usecase"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

type TesterUsecase interface {
	RunCases(tcs []domain.TestCase) (*domain.Report, error)
}

type testerUsecase struct {
	dtuc database_tester_usecase.DatabaseTesterUsecase
}

func NewTesterUsecase(dtuc database_tester_usecase.DatabaseTesterUsecase) TesterUsecase {
	tuc := new(testerUsecase)
	tuc.dtuc = dtuc
	return tuc
}

func (tuc *testerUsecase) RunCases(tcs []domain.TestCase) (*domain.Report, error) {
	r := domain.NewReport()

	for _, tc := range tcs {

		switch tc.ComponentType {

		case domain.ComponentType_Postgres:
			tcr, err := tuc.dtuc.RunCase(&tc)
			if err != nil {
				return nil, err
			}
			r.AddTestCaseResults(tcr)
			logrus.WithField("testResults", tcr).Debug("added test results")

		default:
			return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
		}
	}

	return r, nil
}
