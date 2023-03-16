package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TestingUserRepository struct {
	Users map[string]User
}

func (r TestingUserRepository) Create(user User) (User, error) {
	user.Model.ID = uuid.NewString()
	user.Model.CreatedAt = time.Now()
	r.Users[user.ID] = user
	r.Users[user.Email] = user
	return user, nil
}

func (r TestingUserRepository) GetById(id string) (User, error) {
	if user, ok := r.Users[id]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}

func (r TestingUserRepository) GetByEmail(email string) (User, error) {
	if user, ok := r.Users[email]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}
