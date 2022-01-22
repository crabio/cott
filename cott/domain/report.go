package domain

type Report struct {
	TestCase *TestCase
	Score    float32
	Metrics  []*Metric
}
