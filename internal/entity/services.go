package entity

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user User) (User, error)
	GetById(id string) (User, error)
	GetByEmail(email string) (User, error)
	Exists(id string) bool
}

type repositoryUserService struct {
	userRepository UserRepository
}

func (s repositoryUserService) Create(user User) (User, error) {
	if err := validation.Validate(
		user.Username,
		validation.Required,
		validation.Length(4, 250),
	); err != nil {
		return User{}, errors.New("username is invalid")
	}
	if err := validation.Validate(
		user.Email,
		validation.Required,
		is.Email,
	); err != nil {
		return User{}, errors.New("email is invalid")
	}
	password := user.Password
	if err := validation.Validate(
		password,
		validation.Required,
		validation.Length(8, 500),
	); err != nil {
		return User{}, errors.New("password is invalid")
	}
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("internal server error")
	}
	clone, err := s.userRepository.Create(User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hp),
	})
	if err != nil {
		return User{}, errors.New("internal server error")
	}
	return clone, nil
}

func (s repositoryUserService) GetById(id string) (User, error) {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		return User{}, errors.New("object not found")
	}
	return user, nil
}

func (s repositoryUserService) GetByEmail(email string) (User, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return User{}, errors.New("object not found")
	}
	return user, nil
}

func (s repositoryUserService) Exists(id string) bool {
	_, err := s.GetById(id)
	return err == nil
}

func NewUserService(repository UserRepository) UserService {
	return repositoryUserService{
		userRepository: repository,
	}
}
