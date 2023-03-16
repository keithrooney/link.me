package entity

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/keithrooney/anchorly/internal/utils"
)

func TestUser(t *testing.T) {

	t.Run("TestMarshalJSON", func(t *testing.T) {

		timestamp := time.Now()
		user := &User{
			utils.Model{
				ID:        "asdf234234",
				CreatedAt: timestamp,
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
			"John Doe",
			"johndoe@delivr.co",
			"my.secret.password!.",
		}

		bytes, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
		}

		body := string(bytes)

		if strings.Contains(body, "password") {
			t.Error("expected field not to exist")
		}
		if strings.Contains(body, "deletedAt") || strings.Contains(body, "DeletedAt") {
			t.Error("expected field not to exist")
		}
		if strings.Contains(body, "updatedAt\":{") {
			t.Error("expected field not to exist")
		}

	})

}
