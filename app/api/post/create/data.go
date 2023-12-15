package create

import (
	"github.com/invopop/validation"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

type Data struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (d Data) Validate() []response.ErrorDetails {
	details, _ := helper.Validate(
		&d, ErrCodeInvalidBody,
		validation.Field(
			&d.Title,
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
