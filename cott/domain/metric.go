package domain

type MetricMeta struct {
	Name                string              `json:"name"`
	UnitOfMeasurePrefix UnitOfMeasurePrefix `json:"uom-prefix"`
	UnitOfMeasure       UnitOfMeasure       `json:"uom"`
}

type Metric struct {
	Meta  MetricMeta `json:"meta"`
	Value float64    `json:"value"`
}
