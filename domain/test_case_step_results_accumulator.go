package domain

import (
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
)

type TestCaseStepResultsAccumulator struct {
	testCaseStep *TestCaseStep
	// TODO Refactor onto interface
	metricsMap map[MetricMeta][]float64
	errors     []string
}

func NewTestCaseStepResultsAccumulator(tcs *TestCaseStep) *TestCaseStepResultsAccumulator {
	r := new(TestCaseStepResultsAccumulator)
	r.testCaseStep = tcs
	r.metricsMap = make(map[MetricMeta][]float64)
	return r
}

func (r *TestCaseStepResultsAccumulator) AddMetric(meta *MetricMeta, value float64) {
	logrus.WithFields(logrus.Fields{"meta": *meta, "value": value}).Debug("add test case step result metric")
	if values, ok := r.metricsMap[*meta]; ok {
		r.metricsMap[*meta] = append(values, value)
	} else {
		r.metricsMap[*meta] = []float64{value}
	}
}

func (r *TestCaseStepResultsAccumulator) AddError(err string) {
	r.errors = append(r.errors, err)
}

func (r *TestCaseStepResultsAccumulator) ToTestCaseStepResults() *TestCaseStepResults {
	var metrics []Metric

	for metricMeta, values := range r.metricsMap {
		metrics = append(metrics, Metric{Meta: metricMeta, Value: stat.Mean(values, nil)})
	}

	return &TestCaseStepResults{
		TestCaseStep: *r.testCaseStep,
		Metrics:      metrics,
		Errors:       r.errors,
	}
}
