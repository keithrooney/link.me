package media

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

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

func NewLinkService(repository LinkRepository) LinkService {
	return repositoryLinkService{
		linkRepository: repository,
	}
}
