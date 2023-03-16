package utils

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
