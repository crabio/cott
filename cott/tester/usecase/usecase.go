package usecase

import (
	"github.com/iakrevetkho/components-tests/cott/domain"
)

type TesterUsecase interface {
	RunCase(tc *domain.TestCase) (*domain.Report, error)
}

type testerUsecase struct {
}

func NewTesterUsecase() TesterUsecase {
	return &testerUsecase{}
}

func (tuc *testerUsecase) RunCase(tc *domain.TestCase) (*domain.Report, error) {
	return nil, nil
}
