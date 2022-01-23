package repository

type DatabaseTesterRepository interface {
	Open() error
	CreateDatabase(name string) error
	DropDatabase(name string) error
	Close() error
}
