package domain

type MetricType string

const (
	MetricType_Duration            = "duration"
	MetricType_CpuUsage            = "cpuUsage"
	MetricType_MemoryUsage         = "memoryUsage"
	MetricType_MemoryUsageDiff     = "memoryUsageDiff"
	MetricType_StorageReadUsage    = "storageReadUsage"
	MetricType_StorageWriteUsage   = "storageWriteUsage"
	MetricType_NetworkReceiveUsage = "networkReceiveUsage"
	MetricType_NetworkSendUsage    = "networkSendUsage"
)

type MetricMeta struct {
	Name                string              `json:"name"`
	UnitOfMeasurePrefix UnitOfMeasurePrefix `json:"uom-prefix"`
	UnitOfMeasure       UnitOfMeasure       `json:"uom"`
}

var (
	MetricMeta_Duration            = MetricMeta{Name: "duration", UnitOfMeasurePrefix: UnitOfMeasurePrefix_Micro, UnitOfMeasure: UnitOfMeasure_Second}
	MetricMeta_CpuUsage            = MetricMeta{Name: "cpuUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_Nano, UnitOfMeasure: UnitOfMeasure_Second}
	MetricMeta_MemoryUsage         = MetricMeta{Name: "memoryUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
	MetricMeta_MemoryUsageDiff     = MetricMeta{Name: "memoryUsageDiff", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
	MetricMeta_StorageReadUsage    = MetricMeta{Name: "storageReadUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
	MetricMeta_StorageWriteUsage   = MetricMeta{Name: "storageWriteUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
	MetricMeta_NetworkReceiveUsage = MetricMeta{Name: "networkReceiveUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
	MetricMeta_NetworkSendUsage    = MetricMeta{Name: "networkSendUsage", UnitOfMeasurePrefix: UnitOfMeasurePrefix_None, UnitOfMeasure: UnitOfMeasure_Byte}
)

type Metric struct {
	Meta  MetricMeta `json:"meta"`
	Value float64    `json:"value"`
}
