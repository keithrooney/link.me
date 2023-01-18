package internal

import (
	"database/sql"
	"errors"
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
		"insert into links (id, title, href, created_at, user_id) values ($1, $2, $3, $4, $5)",
		link.ID,
		link.Title,
		link.Href,
		link.CreatedAt,
		link.User.ID,
	); err != nil {
		return Link{}, errors.New("internal server error")
	}
	return link, nil
}

func (r databaseLinkRepository) GetById(id string) (Link, error) {
	link := Link{}
	row, err := r.database.QueryRow(
		`select 
			links.href,
			links.title, 
			links.id, 
			links.created_at, 
			links.updated_at, 
			links.deleted_at, 
			users.username, 
			users.email, 
			users.password,
			users.id, 
			users.created_at, 
			users.updated_at, 
			users.deleted_at 
		from 
			links 
		left join 
			users
		on 
			links.user_id = users.id and 
			links.id = $1
		`,
		id,
	)
	if err != nil {
		return Link{}, err
	}
	if err := row.Scan(
		&link.Href,
		&link.Title,
		&link.Model.ID,
		&link.Model.CreatedAt,
		&link.Model.UpdatedAt,
		&link.Model.DeletedAt,
		&link.User.Username,
		&link.User.Email,
		&link.User.Password,
		&link.User.Model.ID,
		&link.User.Model.CreatedAt,
		&link.User.Model.UpdatedAt,
		&link.User.Model.DeletedAt,
	); err != nil {
		return Link{}, errors.New("object not found")
	}
	return link, nil
}

func NewLinkRepository(database Database) LinkRepository {
	return databaseLinkRepository{
		database: database,
	}
}

type UserRepository interface {
	Create(user User) (User, error)
	GetById(id string) (User, error)
	GetByEmail(email string) (User, error)
}

type databaseUserRepository struct {
	database Database
}

func (r databaseUserRepository) Create(user User) (User, error) {
	user.Model.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	if _, err := r.database.Exec(
		"insert into users (id, username, email, password, created_at) values ($1, $2, $3, $4, $5)",
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
	); err != nil {
		return User{}, errors.New("internal server error")
	}
	return user, nil
}

func (r databaseUserRepository) GetById(id string) (User, error) {
	return r.getBy(
		"select username, email, password, id, created_at, updated_at, deleted_at from users where id = $1",
		id,
	)
}

func (r databaseUserRepository) GetByEmail(email string) (User, error) {
	return r.getBy(
		"select username, email, password, id, created_at, updated_at, deleted_at from users where email = $1",
		email,
	)
}

func (r databaseUserRepository) getBy(query string, value string) (User, error) {
	user := User{}
	row, err := r.database.QueryRow(query, value)
	if err != nil {
		return User{}, err
	}
	if err := row.Scan(
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Model.ID,
		&user.Model.CreatedAt,
		&user.Model.UpdatedAt,
		&user.Model.DeletedAt,
	); err != nil {
		return User{}, errors.New("object not found")
	}
	return user, nil
}

func NewUserRepository(database Database) UserRepository {
	return databaseUserRepository{
		database: database,
	}
}
