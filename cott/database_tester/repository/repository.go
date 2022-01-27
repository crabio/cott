package repository

type DatabaseTesterRepository interface {
	Open() error
	Ping() error
	CreateDatabase(name string) error
	DropDatabase(name string) error
	SwitchDatabase(name string) error
	CreateTable(name string, fields []string) error
	DropTable(name string) error
	SingleInsert(tableName string, columns []string, values []interface{}) error
	MultipleInsert(tableName string, columns []string, values []interface{}) error
	Close() error
}
