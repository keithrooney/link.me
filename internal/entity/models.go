package entity

import (
	"encoding/json"

	"github.com/keithrooney/anchorly/internal/utils"
)

type User struct {
	utils.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		utils.Model
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		user.Model,
		user.Username,
		user.Email,
	})
}
