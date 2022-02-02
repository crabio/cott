package domain

type TestCaseResultsAccumulator struct {
	TestCase                        *TestCase
	testCaseStepResultsAccumulators []*TestCaseStepResultsAccumulator
}

func NewTestCaseResultsAccumulator(tc *TestCase) *TestCaseResultsAccumulator {
	r := new(TestCaseResultsAccumulator)
	r.TestCase = tc
	return r
}

func (r *TestCaseResultsAccumulator) AddTestCaseStepResultsAccumulator(tcsra *TestCaseStepResultsAccumulator) {
	r.testCaseStepResultsAccumulators = append(r.testCaseStepResultsAccumulators, tcsra)
}

func (r *TestCaseResultsAccumulator) ToTestCaseResults() *TestCaseResults {
	tcr := new(TestCaseResults)

	for _, v := range r.testCaseStepResultsAccumulators {
		tcr.StepsResults = append(tcr.StepsResults, v.ToTestCaseStepResults())
	}

	return tcr
}
