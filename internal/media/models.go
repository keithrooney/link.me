package media

import (
	"github.com/keithrooney/anchorly/internal/entity"
	"github.com/keithrooney/anchorly/internal/utils"
)

type Link struct {
	utils.Model
	Href  string      `json:"href"`
	Title string      `json:"title"`
	User  entity.User `json:"-"`
}
