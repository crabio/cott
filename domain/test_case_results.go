package domain

type TestCaseResults struct {
	TestCase     TestCase               `json:"test-case"`
	Score        float32                `json:"score"`
	StepsResults []*TestCaseStepResults `json:"steps-results,omitempty"`
}
