package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TestingUserRepository struct {
	Users  map[string]User
	Errors map[string]error
}

func (r TestingUserRepository) Create(user User) (User, error) {
	if err, present := r.Errors["create"]; present {
		return User{}, err
	} else {
		user.Model.ID = uuid.NewString()
		user.Model.CreatedAt = time.Now()
		r.Users[user.ID] = user
		r.Users[user.Email] = user
		return user, nil
	}
}

func (r TestingUserRepository) GetById(id string) (User, error) {
	if err, present := r.Errors["getById"]; present {
		return User{}, err
	}
	if user, ok := r.Users[id]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}

func (r TestingUserRepository) GetByEmail(email string) (User, error) {
	if err, present := r.Errors["getByEmail"]; present {
		return User{}, err
	}
	if user, ok := r.Users[email]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}

func (r TestingUserRepository) Exists(id string) bool {
	_, err := r.GetById(id)
	return err == nil
}
