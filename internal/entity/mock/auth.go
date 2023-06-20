package mock_models // nolint:stylecheck

import (
	"testing"

	"github.com/018bf/companies/internal/entity"

	"github.com/jaswdr/faker"
)

func NewToken(t *testing.T) entity.Token {
	t.Helper()
	return entity.Token(faker.New().Internet().Password())
}
