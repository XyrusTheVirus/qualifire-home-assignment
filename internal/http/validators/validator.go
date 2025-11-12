package validators

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/models"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(c *gin.Context) models.Model
}

func ConvertDto[T any](c *gin.Context, rule T) *T {

	if reflect.TypeOf(rule).Kind() != reflect.Struct {
		panic(errors.Validation{}.GetError("rule is not a struct", http.StatusInternalServerError))
	}

	v := validator.New()
	err := v.Struct(rule)

	if err != nil {
		panic(errors.Validation{}.GetError(err.Error(), http.StatusBadRequest))
	}

	return &rule

}
