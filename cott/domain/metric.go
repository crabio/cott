package domain

type Metric struct {
	Name                string
	UnitOfMeasurePrefix UnitOfMeasurePrefix
	UnitOfMeasure       UnitOfMeasure
	Value               float64
}
