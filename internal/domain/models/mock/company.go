package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewCompany(t *testing.T) *models.Company {
	t.Helper()
	return &models.Company{
		ID:                models.UUID(uuid.NewString()),
		UpdatedAt:         faker.New().Time().Time(time.Now()),
		CreatedAt:         faker.New().Time().Time(time.Now()),
		Name:              faker.New().Lorem().Text(15),
		Description:       faker.New().Lorem().Sentence(15),
		AmountOfEmployees: faker.New().Int(),
		Registered:        faker.New().Bool(),
		Type:              models.CompanyType(faker.New().Int8Between(1, 4)),
	}
}
func NewCompanyCreate(t *testing.T) *models.CompanyCreate {
	t.Helper()
	return &models.CompanyCreate{
		Name:              faker.New().Lorem().Text(15),
		Description:       faker.New().Lorem().Sentence(15),
		AmountOfEmployees: faker.New().Int(),
		Registered:        faker.New().Bool(),
		Type:              models.CompanyType(faker.New().Int8Between(1, 4)),
	}
}
func NewCompanyUpdate(t *testing.T) *models.CompanyUpdate {
	t.Helper()
	return &models.CompanyUpdate{
		ID:                models.UUID(uuid.NewString()),
		Name:              utils.Pointer(faker.New().Lorem().Text(15)),
		Description:       utils.Pointer(faker.New().Lorem().Sentence(15)),
		AmountOfEmployees: utils.Pointer(faker.New().Int()),
		Registered:        utils.Pointer(faker.New().Bool()),
		Type:              utils.Pointer(models.CompanyType(faker.New().Int8Between(1, 4))),
	}
}
func NewCompanyFilter(t *testing.T) *models.CompanyFilter {
	t.Helper()
	return &models.CompanyFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
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
		Types: []models.CompanyType{
			models.CompanyType(faker.New().UInt8Between(1, 3)),
			models.CompanyType(faker.New().UInt8Between(1, 3)),
		},
		Registered: utils.Pointer(faker.New().Bool()),
	}
}
