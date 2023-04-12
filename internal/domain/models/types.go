package models

import (
	"github.com/018bf/companies/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UUID string

func (u UUID) Validate() error {
	err := validation.Validate(string(u), is.UUID, validation.Required)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
