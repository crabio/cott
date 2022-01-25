package repository

type DatabaseTesterRepository interface {
	Open() error
	Ping() error
	CreateDatabase(name string) error
	DropDatabase(name string) error
	SwitchDatabase(name string) error
	CreateTable(name string, fields []string) error
	DropTable(name string) error
	Close() error
}
