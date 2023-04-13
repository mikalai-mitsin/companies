package mock_models // nolint:stylecheck

import (
	"github.com/jaswdr/faker"
	"testing"

	"github.com/018bf/companies/internal/domain/models"
)

func NewEvent(t *testing.T) *models.Event {
	t.Helper()
	return &models.Event{
		Operation: models.EventOperation(
			faker.New().RandomStringElement([]string{"created", "updated", "deleted"}),
		),
		Company: NewCompany(t),
	}
}
