package internal

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Model struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt sql.NullTime `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}

type Link struct {
	Model
	Href  string `json:"href"`
	Title string `json:"title"`
	User  User   `json:"-"`
}

type User struct {
	Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Model
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		user.Model,
		user.Username,
		user.Email,
	})
}
