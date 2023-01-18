package internal

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

func NewUserService(userRepository UserRepository) UserService {
	return repositoryUserService{
		userRepository: userRepository,
	}
}

type LinkService interface {
	Create(link Link) (Link, error)
	GetById(id string) (Link, error)
}

type repositoryLinkService struct {
	linkRepository LinkRepository
}

func (s repositoryLinkService) Create(link Link) (Link, error) {
	if err := validation.Validate(
		link.Title,
		validation.Required,
		validation.Length(4, 250),
	); err != nil {
		return Link{}, errors.New("title is invalid")
	}
	if err := validation.Validate(
		link.Href,
		validation.Required,
		is.URL,
	); err != nil {
		return Link{}, errors.New("href is invalid")
	}
	if err := validation.Validate(
		link.User.ID,
		validation.Required,
		is.UUID,
	); err != nil {
		return Link{}, errors.New("user is required")
	}
	clone, err := s.linkRepository.Create(link)
	if err != nil {
		return Link{}, errors.New("internal server error")
	}
	return clone, nil
}

func (s repositoryLinkService) GetById(id string) (Link, error) {
	link, err := s.linkRepository.GetById(id)
	if err != nil {
		return Link{}, errors.New("object not found")
	}
	return link, nil
}

func NewLinkService(linkRepository LinkRepository) LinkService {
	return repositoryLinkService{
		linkRepository: linkRepository,
	}
}
