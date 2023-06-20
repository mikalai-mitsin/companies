package mock_models // nolint:stylecheck

import (
	"github.com/jaswdr/faker"
	"testing"

	"github.com/018bf/companies/internal/entity"
)

func NewEvent(t *testing.T) *entity.Event {
	t.Helper()
	return &entity.Event{
		Operation: entity.EventOperation(
			faker.New().RandomStringElement([]string{"created", "updated", "deleted"}),
		),
		Company: NewCompany(t),
	}
}
