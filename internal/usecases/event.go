package usecases

import (
	"context"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
)

type EventUseCase struct {
	eventRepository repositories.EventRepository
	logger          log.Logger
}

func NewEventUseCase(eventRepository repositories.EventRepository, logger log.Logger) usecases.EventUseCase {
	return &EventUseCase{eventRepository: eventRepository, logger: logger}
}

func (u *EventUseCase) CompanyCreated(ctx context.Context, company *models.Company) error {
	event := &models.Event{
		Operation: models.EventTypeCreated,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}

func (u *EventUseCase) CompanyUpdated(ctx context.Context, company *models.Company) error {
	event := &models.Event{
		Operation: models.EventTypeUpdated,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}

func (u *EventUseCase) CompanyDeleted(ctx context.Context, company *models.Company) error {
	event := &models.Event{
		Operation: models.EventTypeDeleted,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}
