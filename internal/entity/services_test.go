package entity

import (
	"testing"
)

func TestUserService(t *testing.T) {

	service := repositoryUserService{
		userRepository: TestingUserRepository{
			Users: make(map[string]User),
		},
	}

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

}
