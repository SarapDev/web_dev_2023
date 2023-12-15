package create

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_Validate(t *testing.T) {
	t.Run("валидные данные", func(t *testing.T) {
		data := Data{
			Title:   "title",
			Content: "content",
		}
		assert.Empty(t, data.Validate())
	})
	t.Run("невалидные данные", func(t *testing.T) {
		data := Data{
			Title:   strings.Repeat("too long", 100),
			Content: "",
		}
		errs := data.Validate()
		assert.Len(t, errs, 2)
		assert.Equal(t, "title: the length must be between 1 and 255", errs[0].Detail)
		assert.Equal(t, "content: cannot be blank", errs[1].Detail)
	})
}
