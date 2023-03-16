package entity

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/keithrooney/anchorly/internal/utils"
	"github.com/ory/dockertest"
	"github.com/pressly/goose"
)

const (
	ACHORLY_DATABASE_USERNAME     = "admin"
	ACHORLY_DATABASE_PASSWORD     = "password"
	ACHORLY_DATABASE_HOST         = "localhost"
	DATABASE_MIGRATIONS_DIRECTORY = "assets/db/migrations/postgresql"
)

// This is set at runtime
var (
	database utils.Database
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

	dataSource := utils.NewDataSource(
		ACHORLY_DATABASE_USERNAME,
		ACHORLY_DATABASE_PASSWORD,
		ACHORLY_DATABASE_HOST,
		resource.GetPort("5432/tcp"),
	)

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

	migrations := filepath.Join(filepath.Dir("../../."), DATABASE_MIGRATIONS_DIRECTORY)
	if err := goose.Up(db, migrations); err != nil {
		log.Fatal(err)
	}

	database = utils.NewDatabase(dataSource)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)

}

func TestUserRepository(t *testing.T) {

	repository := NewUserRepository(database)

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
