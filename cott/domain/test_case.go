package domain

type ComponentType string

const (
	ComponentType_NA       = ""
	ComponentType_Postgres = "postgres"
	ComponentType_Kafka    = "kafka"
)

type TestCase struct {
	ComponentType ComponentType `json:"component-type"`
	Host          string        `json:"host"`
	Port          uint16        `json:"port"`
	User          string        `json:"user"`
	Password      string        `json:"password"`
}
