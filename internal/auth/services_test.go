package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/keithrooney/anchorly/internal/entity"
)

type testUserService struct {
	users map[string]entity.User
}

func (s testUserService) Create(user entity.User) (entity.User, error) {
	user.Model.ID = uuid.NewString()
	user.Model.CreatedAt = time.Now()
	s.users[user.ID] = user
	s.users[user.Email] = user
	return user, nil
}

func (s testUserService) GetById(id string) (entity.User, error) {
	if user, ok := s.users[id]; ok {
		return user, nil
	} else {
		return entity.User{}, errors.New("record not found")
	}
}

func (s testUserService) GetByEmail(email string) (entity.User, error) {
	if user, ok := s.users[email]; ok {
		return user, nil
	} else {
		return entity.User{}, errors.New("record not found")
	}
}

func (s testUserService) Exists(id string) bool {
	_, present := s.users[id]
	return present
}

func TestLoginService(t *testing.T) {

	userRepository := entity.TestingUserRepository{
		Users: make(map[string]entity.User),
	}

	userService := entity.NewUserService(userRepository)

	loginService := NewLoginService(userService)

	t.Run("TestLogin", func(t *testing.T) {

		user := entity.User{
			Username: "Dr. Strange",
			Email:    "drstrange@marvel.com",
			Password: "82392342342",
		}

		if _, err := userService.Create(user); err != nil {
			t.Fatal(err)
		}

		if _, err := loginService.Login(user); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("TestLoginReturnsBadRequest", func(t *testing.T) {

		user := entity.User{
			Username: "Dr. Strange",
			Email:    "mickocollins@marvel.com",
			Password: "82392342342",
		}

		_, err := loginService.Login(user)
		if err == nil {
			t.Fatal(err)
		}
		if err.Error() != "bad request" {
			t.Fatal("Expected a different error to be created.")
		}

	})

	t.Run("TestLoginReturnsPermissionDenied", func(t *testing.T) {

		user := entity.User{
			Username: "Dr. Strange",
			Email:    "drstrange@marvel.com",
			Password: "82392342342",
		}

		if _, err := userService.Create(user); err != nil {
			t.Fatal(err)
		}

		another := entity.User{
			Username: user.Username,
			Email:    user.Email,
			Password: "abcdefghijklmno",
		}

		_, err := loginService.Login(another)
		if err == nil {
			t.Fatal(err)
		}
		if err.Error() != "permission denied" {
			t.Fatal("Expected a different error to be created.")
		}

	})

}
