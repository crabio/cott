package domain

type TestCaseStep struct {
	Name             string
	ContainerID      string
	StepFunc         func() error
	ContainerMetrics []*MetricMeta
}
