package domain

import (
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
