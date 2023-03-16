package media

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/keithrooney/anchorly/internal/utils"
)

type LinkRepository interface {
	Create(link Link) (Link, error)
	GetById(id string) (Link, error)
}

type databaseLinkRepository struct {
	database utils.Database
}

func (r databaseLinkRepository) Create(link Link) (Link, error) {
	link.Model.ID = uuid.NewString()
	link.CreatedAt = time.Now()
	if _, err := r.database.Exec(
		"insert into links (id, title, href, created_at, user_id) values ($1, $2, $3, $4, $5)",
		link.ID,
		link.Title,
		link.Href,
		link.CreatedAt,
		link.User.ID,
	); err != nil {
		return Link{}, errors.New("internal server error")
	}
	return link, nil
}

func (r databaseLinkRepository) GetById(id string) (Link, error) {
	link := Link{}
	row, err := r.database.QueryRow(
		`select 
			links.href,
			links.title, 
			links.id, 
			links.created_at, 
			links.updated_at, 
			links.deleted_at, 
			users.username, 
			users.email, 
			users.password,
			users.id, 
			users.created_at, 
			users.updated_at, 
			users.deleted_at 
		from 
			links 
		left join 
			users
		on 
			links.user_id = users.id and 
			links.id = $1
		`,
		id,
	)
	if err != nil {
		return Link{}, err
	}
	if err := row.Scan(
		&link.Href,
		&link.Title,
		&link.Model.ID,
		&link.Model.CreatedAt,
		&link.Model.UpdatedAt,
		&link.Model.DeletedAt,
		&link.User.Username,
		&link.User.Email,
		&link.User.Password,
		&link.User.Model.ID,
		&link.User.Model.CreatedAt,
		&link.User.Model.UpdatedAt,
		&link.User.Model.DeletedAt,
	); err != nil {
		return Link{}, errors.New("object not found")
	}
	return link, nil
}

func NewLinkRepository(database utils.Database) LinkRepository {
	return databaseLinkRepository{
		database: database,
	}
}
