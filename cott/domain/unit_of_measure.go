package domain

type UnitOfMeasurePrefix uint32

const (
	UnitOfMeasurePrefix_NA = iota
	UnitOfMeasurePrefix_Nano
	UnitOfMeasurePrefix_Micro
	UnitOfMeasurePrefix_Milli
	UnitOfMeasurePrefix_None
	UnitOfMeasurePrefix_Kilo
	UnitOfMeasurePrefix_Mega
	UnitOfMeasurePrefix_Tera
	UnitOfMeasurePrefix_Peta
)

type UnitOfMeasure uint32

const (
	UnitOfMeasure_NA = iota
	UnitOfMeasure_Byte
	UnitOfMeasure_Second
	UnitOfMeasure_Piece
)
