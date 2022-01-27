package repository

type DatabaseTesterRepository interface {
	Open() error
	Ping() error
	CreateDatabase(name string) error
	DropDatabase(name string) error
	SwitchDatabase(name string) error
	CreateTable(name string, fields []string) error
	TruncateTable(name string) error
	DropTable(name string) error
	Insert(tableName string, columns []string, values []map[string]interface{}) error
	Close() error
}
