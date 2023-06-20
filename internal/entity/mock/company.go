package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewCompany(t *testing.T) *entity.Company {
	t.Helper()
	return &entity.Company{
		ID:                entity.UUID(uuid.NewString()),
		UpdatedAt:         faker.New().Time().Time(time.Now()),
		CreatedAt:         faker.New().Time().Time(time.Now()),
		Name:              faker.New().Address().PostCode(),
		Description:       faker.New().Lorem().Sentence(15),
		AmountOfEmployees: faker.New().Int(),
		Registered:        faker.New().Bool(),
		Type:              entity.CompanyType(faker.New().Int8Between(1, 4)),
	}
}
func NewCompanyCreate(t *testing.T) *entity.CompanyCreate {
	t.Helper()
	return &entity.CompanyCreate{
		Name:              faker.New().Lorem().Word(),
		Description:       faker.New().Lorem().Sentence(15),
		AmountOfEmployees: faker.New().Int(),
		Registered:        faker.New().Bool(),
		Type:              entity.CompanyType(faker.New().Int8Between(1, 4)),
	}
}
func NewCompanyUpdate(t *testing.T) *entity.CompanyUpdate {
	t.Helper()
	return &entity.CompanyUpdate{
		ID:                entity.UUID(uuid.NewString()),
		Name:              utils.Pointer(faker.New().Lorem().Text(15)),
		Description:       utils.Pointer(faker.New().Lorem().Sentence(15)),
		AmountOfEmployees: utils.Pointer(faker.New().Int()),
		Registered:        utils.Pointer(faker.New().Bool()),
		Type:              utils.Pointer(entity.CompanyType(faker.New().Int8Between(1, 4))),
	}
}
func NewCompanyFilter(t *testing.T) *entity.CompanyFilter {
	t.Helper()
	return &entity.CompanyFilter{
		IDs:        []entity.UUID{entity.UUID(uuid.NewString()), entity.UUID(uuid.NewString())},
		PageSize:   utils.Pointer(faker.New().UInt64()),
		PageNumber: utils.Pointer(faker.New().UInt64()),
		OrderBy: []string{faker.New().RandomStringElement([]string{
			"id ASC",
			"id DESC",
			"updated_at ASC",
			"updated_at DESC",
			"created_at ASC",
			"created_at DESC",
			"name ASC",
			"name DESC",
			"description ASC",
			"description DESC",
			"amount_of_employees ASC",
			"amount_of_employees DESC",
			"registered ASC",
			"registered DESC",
			"type ASC",
			"type DESC",
		})},
		Search: utils.Pointer(faker.New().Lorem().Sentence(15)),
		Types: []entity.CompanyType{
			entity.CompanyType(faker.New().UInt8Between(1, 3)),
			entity.CompanyType(faker.New().UInt8Between(1, 3)),
		},
		Registered: utils.Pointer(faker.New().Bool()),
	}
}
