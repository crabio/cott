package repository

import (
	"bytes"
	"database/sql"
	"strconv"

	"github.com/iakrevetkho/components-tests/cott/domain"

	_ "github.com/lib/pq"
)

type postgresDatabaseTesterRepository struct {
	db       *sql.DB
	port     uint16
	host     string
	user     string
	password string
}

func NewPostgresDatabaseTesterRepository(port uint16, host, user, password string) DatabaseTesterRepository {
	r := new(postgresDatabaseTesterRepository)
	r.port = port
	r.host = host
	r.user = user
	r.password = password
	return r
}

func (r *postgresDatabaseTesterRepository) Open() error {
	var err error
	r.db, err = sql.Open("postgres", r.createConnString(r.port, r.host, r.user, r.password))
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresDatabaseTesterRepository) CreateDatabase(name string) error {
	if r.db == nil {
		return domain.CONNECTION_WAS_NOT_ESTABLISHED
	}

	var buf bytes.Buffer
	buf.WriteString("CREATE DATABASE ")
	buf.WriteString(name)

	_, err := r.db.Exec(buf.String())
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresDatabaseTesterRepository) DropDatabase(name string) error {
	if r.db == nil {
		return domain.CONNECTION_WAS_NOT_ESTABLISHED
	}

	var buf bytes.Buffer
	buf.WriteString("DROP DATABASE ")
	buf.WriteString(name)

	_, err := r.db.Exec(buf.String())
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresDatabaseTesterRepository) Close() error {
	if r.db == nil {
		return domain.CONNECTION_WAS_NOT_ESTABLISHED
	}

	if err := r.db.Close(); err != nil {
		return err
	}

	return nil
}

func (r *postgresDatabaseTesterRepository) createConnString(port uint16, host, user, password string) string {
	var buf bytes.Buffer

	buf.WriteString("host=")
	buf.WriteString(host)
	buf.WriteString(" port=")
	buf.WriteString(strconv.FormatUint(uint64(port), 10))
	buf.WriteString(" user=")
	buf.WriteString(user)
	buf.WriteString(" password=")
	buf.WriteString(password)
	buf.WriteString(" sslmode=disable")

	return buf.String()
}
