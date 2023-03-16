package utils

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type DataSource interface {
	Connect() (*sql.DB, error)
}

type postgresDataSource struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (s postgresDataSource) Connect() (*sql.DB, error) {
	source := fmt.Sprintf("postgresql://%s:%s@%s:%s/anchorly?sslmode=disable", s.Username, s.Password, s.Host, s.Port)
	return sql.Open("postgres", source)
}

func NewDataSource(username string, password string, host string, port string) DataSource {
	return postgresDataSource{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) (*sql.Row, error)
}

type dataSourceDatabase struct {
	dataSource DataSource
}

func (d dataSourceDatabase) Exec(query string, args ...any) (sql.Result, error) {
	db, err := d.dataSource.Connect()
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return db.Exec(query, args...)
}

func (d dataSourceDatabase) QueryRow(query string, args ...any) (*sql.Row, error) {
	db, err := d.dataSource.Connect()
	if err != nil {
		return nil, err
	}
	return db.QueryRow(query, args...), nil
}

func NewDatabase(dataSource DataSource) Database {
	return dataSourceDatabase{
		dataSource: dataSource,
	}
}
