package usecase

import (
	cl_usecase "github.com/iakrevetkho/components-tests/cott/container_launcher/usecase"
	dt_usecase "github.com/iakrevetkho/components-tests/cott/database_tester/usecase"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

type TesterUsecase interface {
	RunCases(tcs []domain.TestCase) (*domain.Report, error)
}

type testerUsecase struct {
	cluc cl_usecase.ContainerLauncherUsecase
	dtuc dt_usecase.DatabaseTesterUsecase
}

func NewTesterUsecase(cluc cl_usecase.ContainerLauncherUsecase, dtuc dt_usecase.DatabaseTesterUsecase) TesterUsecase {
	tuc := new(testerUsecase)
	tuc.cluc = cluc
	tuc.dtuc = dtuc
	return tuc
}

func (tuc *testerUsecase) RunCases(tcs []domain.TestCase) (*domain.Report, error) {
	r := domain.NewReport()

	for _, tc := range tcs {

		switch tc.ComponentType {

		case domain.ComponentType_Postgres:
			containerId, err := tuc.cluc.LaunchContainer(tc.Image, tc.EnvVars, tc.Port)
			if err != nil {
				return nil, err
			}

			tcr, err := tuc.dtuc.RunCase(&tc)
			if err != nil {
				return nil, err
			}
			r.AddTestCaseResults(tcr)
			logrus.WithField("testResults", tcr).Debug("added test results")

			if err := tuc.cluc.StopContainer(*containerId); err != nil {
				return nil, err
			}

			if err := tuc.cluc.RemoveContainer(*containerId); err != nil {
				return nil, err
			}

		default:
			return nil, domain.UNKNOWN_COMPONENT_FOR_TESTING
		}
	}

	return r, nil
}
