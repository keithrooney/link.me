package internal

import (
	"database/sql"
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
}
