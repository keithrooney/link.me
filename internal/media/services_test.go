package media

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/keithrooney/anchorly/internal/entity"
	"github.com/keithrooney/anchorly/internal/utils"
)

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

	linkService := repositoryLinkService{
		linkRepository: testLinkRepository{
			links: make(map[string]Link),
		},
	}

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
			User: entity.User{
				Model: utils.Model{
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
