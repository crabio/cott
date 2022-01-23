package domain

type ComponentType uint32

const (
	ComponentType_NA = iota
	ComponentType_Postgres
	ComponentType_Kafka
)

type TestCase struct {
	ComponentType ComponentType
	Host          string
	Port          uint16
	User          string
	Password      string
}
