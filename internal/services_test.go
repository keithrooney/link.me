package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type testUserRepository struct {
	users map[string]User
}

func (s testUserRepository) Create(user User) (User, error) {
	user.Model.ID = uuid.NewString()
	user.Model.CreatedAt = time.Now()
	s.users[user.ID] = user
	s.users[user.Email] = user
	return user, nil
}

func (s testUserRepository) GetById(id string) (User, error) {
	if user, ok := s.users[id]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}

func (s testUserRepository) GetByEmail(email string) (User, error) {
	if user, ok := s.users[email]; ok {
		return user, nil
	} else {
		return User{}, errors.New("record not found")
	}
}

func TestUserService(t *testing.T) {

	service := NewUserService(
		testUserRepository{
			users: make(map[string]User),
		},
	)

	t.Run("TestCreateWithInvalidFields", func(t *testing.T) {

		users := []User{
			{},
			{
				Username: "",
			},
			{
				Username: "Kei",
			},
			{
				Username: "Batman",
			},
			{
				Username: "Batman",
				Email:    "foobar",
			},
			{
				Username: "clarky",
				Email:    "superman@dc.com",
			},
			{
				Username: "clarky",
				Email:    "superman@dc.com",
			},
			{
				Username: "Clarky",
				Email:    "superman@dc.com",
				Password: "1234567",
			},
		}

		for _, user := range users {
			_, err := service.Create(user)
			if err == nil {
				t.Fatal("object should be invalid")
			}
		}

	})

	t.Run("TestGetById", func(t *testing.T) {

		user := User{
			Username: "starky",
			Email:    "ironman@marvel.com",
			Password: "1234567890",
		}

		other, err := service.Create(user)
		if err != nil {
			t.Fatal(err)
		}

		if !service.Exists(other.ID) {
			t.Fatal("object not found")
		}

		another, _ := service.GetById(other.ID)
		if other != another {
			t.Fatal("expected != actual")
		}

	})

	t.Run("TestGetByEmail", func(t *testing.T) {

		user := User{
			Username: "Wolfy",
			Email:    "wolverine@marvel.com",
			Password: "82392342342",
		}

		other, err := service.Create(user)
		if err != nil {
			t.Fatal(err)
		}

		if !service.Exists(other.ID) {
			t.Fatal("object not found")
		}

		another, _ := service.GetByEmail(other.Email)
		if other != another {
			t.Fatal("expected != actual")
		}

	})

	t.Run("TestLogin", func(t *testing.T) {

		user := User{
			Username: "Dr. Strange",
			Email:    "drstrange@marvel.com",
			Password: "82392342342",
		}

		if _, err := service.Create(user); err != nil {
			t.Fatal(err)
		}

		if _, err := service.(LoginService).Login(user); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("TestLoginReturnsBadRequest", func(t *testing.T) {

		user := User{
			Username: "Dr. Strange",
			Email:    "mickocollins@marvel.com",
			Password: "82392342342",
		}

		_, err := service.(LoginService).Login(user)
		if err == nil {
			t.Fatal(err)
		}
		if err.Error() != "bad request" {
			t.Fatal("Expected a different error to be created.")
		}

	})

	t.Run("TestLoginReturnsPermissionDenied", func(t *testing.T) {

		user := User{
			Username: "Dr. Strange",
			Email:    "drstrange@marvel.com",
			Password: "82392342342",
		}

		if _, err := service.Create(user); err != nil {
			t.Fatal(err)
		}

		another := User{
			Username: user.Username,
			Email:    user.Email,
			Password: "abcdefghijklmno",
		}

		_, err := service.(LoginService).Login(another)
		if err == nil {
			t.Fatal(err)
		}
		if err.Error() != "permission denied" {
			t.Fatal("Expected a different error to be created.")
		}

	})

}

type testLinkRepository struct {
	links map[string]Link
}

func (r testLinkRepository) Create(link Link) (Link, error) {
	link.Model.ID = uuid.NewString()
	link.CreatedAt = time.Now()
	r.links[link.ID] = link
	return link, nil
}

func (r testLinkRepository) GetById(id string) (Link, error) {
	if user, ok := r.links[id]; ok {
		return user, nil
	} else {
		return Link{}, errors.New("object not found")
	}
}

func TestLinkService(t *testing.T) {

	linkService := NewLinkService(testLinkRepository{
		links: make(map[string]Link),
	})

	t.Run("TestCreateWithInvalidFields", func(t *testing.T) {

		links := []Link{
			{},
			{
				Title: "123",
			},
			{
				Title: "This is our site!",
			},
			{
				Title: "This is our site!",
				Href:  "site",
			},
			{
				Title: "This is our site!",
				Href:  "www.site.com",
			},
		}

		for _, link := range links {
			if _, err := linkService.Create(link); err == nil {
				t.Fatal("this test should not fail")
			}
		}

	})

	t.Run("TestCreate", func(t *testing.T) {

		link, err := linkService.Create(Link{
			Title: "This is our site!",
			Href:  "www.site.com",
			User: User{
				Model: Model{
					ID: uuid.NewString(),
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		other, err := linkService.GetById(link.Model.ID)
		if err != nil {
			t.Fatal(err)
		}

		other.Model.CreatedAt = link.Model.CreatedAt
		other.Model.UpdatedAt = link.Model.UpdatedAt
		other.Model.DeletedAt = link.Model.DeletedAt

		other.User.Model.CreatedAt = link.User.Model.CreatedAt
		other.User.Model.UpdatedAt = link.User.Model.UpdatedAt
		other.User.Model.DeletedAt = link.User.Model.DeletedAt

		if link != other {
			t.Fatal("expected != actual")
		}

	})

	t.Run("TestGetByIdReturnsNotFound", func(t *testing.T) {
		if _, err := linkService.GetById(uuid.NewString()); err == nil {
			t.Fatal("object should not exist")
		}
	})

}
