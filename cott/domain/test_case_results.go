package domain

type TestCaseResults struct {
	TestCase TestCase `json:"test-case"`
	Score    float32  `json:"score"`
	Metrics  []Metric `json:"metrics,omitempty"`
	Error    *string  `json:"error,omitempty"`
}

func NewTestCaseResults(tc *TestCase) *TestCaseResults {
	r := new(TestCaseResults)
	r.TestCase = *tc
	return r
}

func (r *TestCaseResults) AddMetric(name string, uofp UnitOfMeasurePrefix, uom UnitOfMeasure, value float64) {
	r.Metrics = append(r.Metrics, Metric{
		Name:                name,
		UnitOfMeasurePrefix: uofp,
		UnitOfMeasure:       uom,
		Value:               value,
	})
}

func (r *TestCaseResults) AddError(err string) {
	r.Error = &err
}
