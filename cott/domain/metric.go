package domain

type Metric struct {
	Name                string              `json:"name"`
	UnitOfMeasurePrefix UnitOfMeasurePrefix `json:"uom-prefix"`
	UnitOfMeasure       UnitOfMeasure       `json:"uom"`
	Value               float64             `json:"value"`
}
