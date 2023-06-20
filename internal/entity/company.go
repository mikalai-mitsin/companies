package entity

import (
	"time"

	"github.com/018bf/companies/internal/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Company's permissions.
const (
	PermissionIDCompanyList   PermissionID = "company_list"
	PermissionIDCompanyDetail PermissionID = "company_detail"
	PermissionIDCompanyCreate PermissionID = "company_create"
	PermissionIDCompanyUpdate PermissionID = "company_update"
	PermissionIDCompanyDelete PermissionID = "company_delete"
)

const (
	CompanyTypeCorporations CompanyType = iota + 1
	CompanyTypeNonProfit
	CompanyTypeCooperative
	CompanyTypeSoleProprietorship
)

type CompanyType uint8

func (c CompanyType) Validate() error {
	err := validation.Validate(uint8(c), validation.In(
		uint8(CompanyTypeCorporations),
		uint8(CompanyTypeNonProfit),
		uint8(CompanyTypeCooperative),
		uint8(CompanyTypeSoleProprietorship),
	))
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type Company struct {
	ID                UUID        `json:"id"`
	UpdatedAt         time.Time   `json:"updated_at"`
	CreatedAt         time.Time   `json:"created_at"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	AmountOfEmployees int         `json:"amount_of_employees"`
	Registered        bool        `json:"registered"`
	Type              CompanyType `json:"type"`
}

func (m *Company) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.Name, validation.Required, validation.RuneLength(1, 15)),
		validation.Field(&m.Description, validation.RuneLength(0, 3000)),
		validation.Field(&m.AmountOfEmployees, validation.Required),
		validation.Field(&m.Registered),
		validation.Field(&m.Type, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CompanyCreate struct {
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	AmountOfEmployees int         `json:"amount_of_employees"`
	Registered        bool        `json:"registered"`
	Type              CompanyType `json:"type"`
}

func (m *CompanyCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required, validation.RuneLength(1, 15)),
		validation.Field(&m.Description, validation.RuneLength(0, 3000)),
		validation.Field(&m.AmountOfEmployees, validation.Required),
		validation.Field(&m.Registered),
		validation.Field(&m.Type, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CompanyUpdate struct {
	ID                UUID         `json:"id"`
	Name              *string      `json:"name"`
	Description       *string      `json:"description"`
	AmountOfEmployees *int         `json:"amount_of_employees"`
	Registered        *bool        `json:"registered"`
	Type              *CompanyType `json:"type"`
}

func (m *CompanyUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Name, validation.RuneLength(1, 15)),
		validation.Field(&m.Description, validation.RuneLength(0, 3000)),
		validation.Field(&m.AmountOfEmployees),
		validation.Field(&m.Registered),
		validation.Field(&m.Type),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CompanyFilter struct {
	IDs        []UUID        `json:"ids" form:"ids"`
	PageSize   *uint64       `json:"page_size" form:"page_size"`
	PageNumber *uint64       `json:"page_number" form:"page_number"`
	OrderBy    []string      `json:"order_by" form:"order_by"`
	Search     *string       `json:"search" form:"search"`
	Types      []CompanyType `json:"types" form:"types"`
	Registered *bool         `json:"registered" form:"registered"`
}

func (m *CompanyFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.IDs),
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.OrderBy, validation.Each(validation.In(
			"id ASC", "id DESC",
			"updated_at ASC", "updated_at DESC",
			"created_at ASC", "created_at DESC",
			"name ASC", "name DESC",
			"description ASC", "description DESC",
			"amount_of_employees ASC", "amount_of_employees DESC",
			"registered ASC", "registered DESC",
			"type ASC", "type DESC",
		))),
		validation.Field(&m.Search),
		validation.Field(&m.Types),
		validation.Field(&m.Registered),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
