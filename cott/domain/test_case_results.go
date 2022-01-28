package domain

type TestCaseResults struct {
	TestCase TestCase `json:"test-case"`
	Score    float32  `json:"score"`
	Metrics  []Metric `json:"metrics,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}
