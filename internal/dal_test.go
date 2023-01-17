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

	datasource := NewDataSource(ACHORLY_DATABASE_USERNAME, ACHORLY_DATABASE_PASSWORD, ACHORLY_DATABASE_HOST, resource.GetPort("5432/tcp"))
	if err := pool.Retry(func() error {
		var err error
		db, err := datasource.Connect()
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatal(err)
	}

	db, err := datasource.Connect()
	if err != nil {
		log.Fatal(err)
	}

	migrations := filepath.Join(filepath.Dir("../."), DATABASE_MIGRATIONS_DIRECTORY)
	if err := goose.Up(db, migrations); err != nil {
		log.Fatal(err)
	}

	DATABASE = NewDatabase(datasource)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)

}

func TestDataSource(t *testing.T) {

}

func TestLinkRepository(t *testing.T) {

	repository := NewLinkRepository(DATABASE)

	link := Link{
		Title: "This my great website!",
		Href:  "https://www.mygreatwebsite.com",
	}

	link, err := repository.Create(link)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestGetById", func(t *testing.T) {

		other, err := repository.GetById(link.ID)
		if err != nil {
			t.Fatal(err)
		}

		other.Model.CreatedAt = link.Model.CreatedAt
		other.Model.UpdatedAt = link.Model.UpdatedAt
		other.Model.DeletedAt = link.Model.DeletedAt

		if link != other {
			t.Fatal("expected != actual")
		}

	})

	t.Run("TestGetByIdReturnsNil", func(t *testing.T) {

		_, err := repository.GetById(uuid.NewString())
		if err == nil {
			t.Fatal(err)
		}

	})

}
