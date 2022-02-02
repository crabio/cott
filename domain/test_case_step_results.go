package domain

type TestCaseStepResults struct {
	TestCaseStep TestCaseStep `json:"step"`
	Metrics      []Metric     `json:"metrics,omitempty"`
	Errors       []string     `json:"errors,omitempty"`
}
