package domain

type UnitOfMeasurePrefix string

const (
	UnitOfMeasurePrefix_Nano  = "nano"
	UnitOfMeasurePrefix_Micro = "micro"
	UnitOfMeasurePrefix_Milli = "milli"
	UnitOfMeasurePrefix_None  = ""
	UnitOfMeasurePrefix_Kilo  = "kilo"
	UnitOfMeasurePrefix_Mega  = "mega"
	UnitOfMeasurePrefix_Tera  = "tera"
	UnitOfMeasurePrefix_Peta  = "peta"
)

type UnitOfMeasure string

const (
	UnitOfMeasure_Byte   = "byte"
	UnitOfMeasure_Second = "second"
	UnitOfMeasure_Piece  = "piece"
)
