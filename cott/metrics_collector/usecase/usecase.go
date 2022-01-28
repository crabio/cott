package usecase

import (
	"time"

	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

type MetricsCollectorUsecase interface {
	CalcStepDuration(stepFunc func() error, stepName string) error
}

type metricsCollectorUsecase struct {
	tcra *domain.TestCaseResultsAccumulator
}

func NewMetricsCollectorUsecase(tcra *domain.TestCaseResultsAccumulator) MetricsCollectorUsecase {
	mcuc := new(metricsCollectorUsecase)
	mcuc.tcra = tcra
	return mcuc
}

func (mcuc *metricsCollectorUsecase) CalcStepDuration(stepFunc func() error, stepName string) error {
	start := time.Now()
	if err := stepFunc(); err != nil {
		logrus.WithError(err).WithField("stepName", stepName).Warn("error on step execution")
		mcuc.tcra.AddError(stepName + ". " + err.Error())
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "stepName": stepName}).Debug("step finished")
	mcuc.tcra.AddMetric(stepName+"Duration", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}
