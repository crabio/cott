package domain

type Report struct {
	TestCaseResults []*TestCaseResults `json:"test-case-results"`
}

func NewReport() *Report {
	r := new(Report)
	return r
}

func (r *Report) AddTestCaseResults(tcr *TestCaseResults) {
	r.TestCaseResults = append(r.TestCaseResults, tcr)
}
