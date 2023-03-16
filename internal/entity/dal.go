package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/keithrooney/anchorly/internal/utils"
)

type UserRepository interface {
	Create(user User) (User, error)
	GetById(id string) (User, error)
	GetByEmail(email string) (User, error)
}

type databaseUserRepository struct {
	database utils.Database
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

func NewUserRepository(database utils.Database) UserRepository {
	return databaseUserRepository{
		database: database,
	}
}
