package internal

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DataSource interface {
	Connect() (*sql.DB, error)
}

type postgresDataSource struct {
	username string
	password string
	host     string
	port     string
}

func (s postgresDataSource) Connect() (*sql.DB, error) {
	source := fmt.Sprintf("postgresql://%s:%s@%s:%s/anchorly?sslmode=disable", s.username, s.password, s.host, s.port)
	return sql.Open("postgres", source)
}

var (
	DEFAULT_DATABASE_USERNAME = "admin"
	DEFAULT_DATABASE_HOST     = "localhost"
	DEFAULT_DATABASE_PORT     = "5432"
)

func NewDataSource(username string, password string, host string, port string) DataSource {
	if username == "" {
		username = DEFAULT_DATABASE_USERNAME
	}
	if password == "" {
		password = os.Getenv("ANCHORLY_DATABASE_PASSWORD")
	}
	if host == "" {
		host = DEFAULT_DATABASE_HOST
	}
	if port == "" {
		port = DEFAULT_DATABASE_PORT
	}
	return postgresDataSource{
		username: username,
		password: password,
		host:     host,
		port:     port,
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
		return nil, err
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

type LinkRepository interface {
	Create(link Link) (Link, error)
	GetById(id string) (Link, error)
}

type databaseLinkRepository struct {
	database Database
}

func (r databaseLinkRepository) Create(link Link) (Link, error) {
	link.Model.ID = uuid.NewString()
	link.CreatedAt = time.Now()
	if _, err := r.database.Exec(
		"insert into links (id, title, href, created_at) values ($1, $2, $3, $4)",
		link.ID,
		link.Title,
		link.Href,
		link.CreatedAt,
	); err != nil {
		return Link{}, err
	}
	return link, nil
}

func (r databaseLinkRepository) GetById(id string) (Link, error) {
	link := Link{}
	row, err := r.database.QueryRow("select title, href, id, created_at, updated_at, deleted_at from links where id = $1", id)
	if err != nil {
		return Link{}, err
	}
	if err := row.Scan(
		&link.Title,
		&link.Href,
		&link.Model.ID,
		&link.Model.CreatedAt,
		&link.Model.UpdatedAt,
		&link.Model.DeletedAt,
	); err != nil {
		return Link{}, err
	}
	return link, nil
}

func NewLinkRepository(database Database) LinkRepository {
	return databaseLinkRepository{
		database: database,
	}
}
