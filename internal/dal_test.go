package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/pressly/goose"
)

type Thing struct {
	ID string
}

const (
	ACHORLY_DATABASE_USERNAME     = "admin"
	ACHORLY_DATABASE_PASSWORD     = "password"
	ACHORLY_DATABASE_HOST         = "localhost"
	DATABASE_MIGRATIONS_DIRECTORY = "assets/db/migrations/postgresql"
)

// This is set at runtime
var (
	DATABASE Database
)

func TestMain(m *testing.M) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}

	resource, err := pool.Run(
		"postgres",
		"latest",
		[]string{
			fmt.Sprintf("POSTGRES_USER=%s", ACHORLY_DATABASE_USERNAME),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", ACHORLY_DATABASE_PASSWORD),
			"POSTGRES_DB=anchorly",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	dataSource = postgresDataSource{
		Username: ACHORLY_DATABASE_USERNAME,
		Password: ACHORLY_DATABASE_PASSWORD,
		Host:     ACHORLY_DATABASE_HOST,
		Port:     resource.GetPort("5432/tcp"),
	}

	if err := pool.Retry(func() error {
		var err error
		db, err := dataSource.Connect()
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatal(err)
	}

	db, err := dataSource.Connect()
	if err != nil {
		log.Fatal(err)
	}

	migrations := filepath.Join(filepath.Dir("../."), DATABASE_MIGRATIONS_DIRECTORY)
	if err := goose.Up(db, migrations); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)

}

func TestLinkRepository(t *testing.T) {

	userRepository := NewUserRepository()

	user, err := userRepository.Create(User{
		Username: "jeandoe",
		Email:    "jeandoe@anchorly.com",
		Password: "this.is.my.password!. ",
	})
	if err != nil {
		t.Fatal(err)
	}

	linkRepository := NewLinkRepository()

	link, err := linkRepository.Create(Link{
		Title: "This my great website!",
		Href:  "https://www.mygreatwebsite.com",
		User:  user,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestGetById", func(t *testing.T) {

		other, err := linkRepository.GetById(link.ID)
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

	t.Run("TestGetByIdReturnsNil", func(t *testing.T) {

		_, err := linkRepository.GetById(uuid.NewString())
		if err == nil {
			t.Fatal(err)
		}

	})

}

func TestUserRepository(t *testing.T) {

	repository := NewUserRepository()

	user := User{
		Username: "johndoe",
		Email:    "johndoe@foobar.com",
		Password: "this.is.my.password!. ",
	}

	user, err := repository.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestGetById", func(t *testing.T) {

		other, err := repository.GetById(user.ID)
		if err != nil {
			t.Fatal(err)
		}

		other.Model.CreatedAt = user.Model.CreatedAt
		other.Model.UpdatedAt = user.Model.UpdatedAt
		other.Model.DeletedAt = user.Model.DeletedAt

		if user != other {
			t.Fatal("expected != actual")
		}

	})

	t.Run("TestGetByEmail", func(t *testing.T) {

		other, err := repository.GetByEmail(user.Email)
		if err != nil {
			t.Fatal(err)
		}

		other.Model.CreatedAt = user.Model.CreatedAt
		other.Model.UpdatedAt = user.Model.UpdatedAt
		other.Model.DeletedAt = user.Model.DeletedAt

		if user != other {
			t.Fatal("expected != actual")
		}

	})

}
