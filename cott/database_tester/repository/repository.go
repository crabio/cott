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
	SelectById(tableName string, id uint64) error
	SelectByConditions(tableName string, conditions string) error
	Close() error
}
