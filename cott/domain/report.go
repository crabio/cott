package domain

type Report struct {
	TestCase *TestCase `json:"test-case"`
	Score    float32   `json:"score"`
	Metrics  []*Metric `json:"metrics"`
}

func NewReport(tc *TestCase) *Report {
	r := new(Report)
	r.TestCase = tc
	return r
}

func (r *Report) AddMetric(name string, uofp UnitOfMeasurePrefix, uom UnitOfMeasure, value float64) {
	r.Metrics = append(r.Metrics, &Metric{
		Name:                name,
		UnitOfMeasurePrefix: uofp,
		UnitOfMeasure:       uom,
		Value:               value,
	})
}
