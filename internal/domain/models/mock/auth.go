package mock_models // nolint:stylecheck

import (
	"testing"

	"github.com/018bf/companies/internal/domain/models"

	"github.com/jaswdr/faker"
)

func NewToken(t *testing.T) models.Token {
	t.Helper()
	return models.Token(faker.New().Internet().Password())
}
