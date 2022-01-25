package domain

type TestCaseResults struct {
	TestCase TestCase `json:"test-case"`
	Score    float32  `json:"score"`
	Metrics  []Metric `json:"metrics,omitempty"`
	Error    *string  `json:"error,omitempty"`
}
