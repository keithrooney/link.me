package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/keithrooney/anchorly/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

var secret []byte

func init() {
	if value, found := os.LookupEnv("ANCHORLY_TOKEN_KEY"); !found {
		panic("Expected environment variable to be configured.")
	} else {
		secret = []byte(value)
	}
}

type LoginService interface {
	Login(user entity.User) (Token, error)
}

type loginService struct {
	userService entity.UserService
}

func (s loginService) Login(other entity.User) (Token, error) {
	user, err := s.userService.GetByEmail(other.Email)
	if err != nil {
		return Token{}, errors.New("bad request")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(other.Password)); err != nil {
		return Token{}, errors.New("permission denied")
	}
	claims := jwt.MapClaims{
		"iss": "anchorly.com",
		"sub": user.Model.ID,
		"aud": user.Model.ID,
		"exp": time.Now().Add(time.Hour * 3).UnixMilli(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	value, err := token.SignedString(secret)
	if err != nil {
		return Token{}, errors.New("internal server error")
	}
	return Token{
		Value: value,
	}, nil
}

func NewLoginService(service entity.UserService) LoginService {
	return loginService{
		userService: service,
	}
}
