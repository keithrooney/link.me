package media

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/keithrooney/anchorly/internal/entity"
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

func TestLinkRepository(t *testing.T) {

	userRepository := entity.NewUserRepository(database)

	user, err := userRepository.Create(entity.User{
		Username: "jeandoe",
		Email:    "jeandoe@anchorly.com",
		Password: "this.is.my.password!. ",
	})
	if err != nil {
		t.Fatal(err)
	}

	linkRepository := NewLinkRepository(database)

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
