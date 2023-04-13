package repositories

import (
	"context"
	"github.com/018bf/companies/internal/domain/models"
)

// EventRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/event.go . EventRepository
type EventRepository interface {
	Send(ctx context.Context, event *models.Event) error
}
