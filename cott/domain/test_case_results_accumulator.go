package domain

import (
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
)

type TestCaseResultsAccumulator struct {
	TestCase   *TestCase
	MetricsMap map[MetricMeta][]float64
	Errors     []string
}

func NewTestCaseResultsAccumulator(tc *TestCase) *TestCaseResultsAccumulator {
	r := new(TestCaseResultsAccumulator)
	r.TestCase = tc
	r.MetricsMap = make(map[MetricMeta][]float64)
	return r
}

func (r *TestCaseResultsAccumulator) AddMetric(name string, uofp UnitOfMeasurePrefix, uom UnitOfMeasure, value float64) {
	metricMeta := MetricMeta{
		Name:                name,
		UnitOfMeasurePrefix: uofp,
		UnitOfMeasure:       uom,
	}

	if values, ok := r.MetricsMap[metricMeta]; ok {
		r.MetricsMap[metricMeta] = append(values, value)
	} else {
		r.MetricsMap[metricMeta] = []float64{value}
	}
}

func (r *TestCaseResultsAccumulator) CalcStepDuration(stepFunc func() error, stepName string) error {
	start := time.Now()
	if err := stepFunc(); err != nil {
		logrus.WithError(err).WithField("stepName", stepName).Warn("error on step execution")
		r.AddError(stepName + ". " + err.Error())
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "stepName": stepName}).Debug("step finished")
	r.AddMetric(stepName+"Duration", UnitOfMeasurePrefix_Micro, UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}

func (r *TestCaseResultsAccumulator) AddError(err string) {
	r.Errors = append(r.Errors, err)
}

func (r *TestCaseResultsAccumulator) ToTestCaseResults() *TestCaseResults {
	var metrics []Metric

	for metricMeta, values := range r.MetricsMap {
		metrics = append(metrics, Metric{Meta: metricMeta, Value: stat.Mean(values, nil)})
	}

	return &TestCaseResults{
		TestCase: *r.TestCase,
		Metrics:  metrics,
		Errors:   r.Errors,
	}
}
