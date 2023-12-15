package create

import (
	"github.com/invopop/validation"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

type Data struct {
	PostID  int    `json:"post_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func (d Data) Validate() []response.ErrorDetails {
	details, _ := helper.Validate(
		&d, ErrCodeInvalidBody,
		validation.Field(
			&d.PostID,
			validation.Required,
		),
		validation.Field(
			&d.Author,
			validation.Required,
			validation.Length(1, 255), //nolint:gomnd
		),
		validation.Field(
			&d.Content,
			validation.Required,
			validation.Length(1, 65535), //nolint:gomnd
		),
	)

	return details
}
