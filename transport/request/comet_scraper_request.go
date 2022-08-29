package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// CreateCometScraperReq represent create comet request body
type CreateCometScraperReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request CreateCometScraperReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
	)
}
