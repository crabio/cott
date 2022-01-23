package usecase

import (
	database_tester_usecase "github.com/iakrevetkho/components-tests/cott/database_tester/usecase"
	"github.com/iakrevetkho/components-tests/cott/domain"
)

type TesterUsecase interface {
	RunCase(tc *domain.TestCase) (*domain.Report, error)
}

type testerUsecase struct {
	dtuc database_tester_usecase.DatabaseTesterUsecase
}

func NewTesterUsecase(dtuc database_tester_usecase.DatabaseTesterUsecase) TesterUsecase {
	tuc := new(testerUsecase)
	tuc.dtuc = dtuc
	return tuc
}

func (tuc *testerUsecase) RunCase(tc *domain.TestCase) (*domain.Report, error) {
	switch tc.ComponentType {
	case domain.ComponentType_Postgres:
		return tuc.dtuc.RunCase(tc)

	default:
		return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
	}
}
