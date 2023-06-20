package service

import (
	"context"
	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
)

//go:generate mockgen -source=event.go -package=usecases -destination=event_mock.go

type eventRepository interface {
	Send(ctx context.Context, event *entity.Event) error
}

type EventService struct {
	eventRepository eventRepository
	logger          log.Logger
}

func NewEventService(eventRepository eventRepository, logger log.Logger) *EventService {
	return &EventService{eventRepository: eventRepository, logger: logger}
}

func (u *EventService) CompanyCreated(ctx context.Context, company *entity.Company) error {
	event := &entity.Event{
		Operation: entity.EventTypeCreated,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}

func (u *EventService) CompanyUpdated(ctx context.Context, company *entity.Company) error {
	event := &entity.Event{
		Operation: entity.EventTypeUpdated,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}

func (u *EventService) CompanyDeleted(ctx context.Context, company *entity.Company) error {
	event := &entity.Event{
		Operation: entity.EventTypeDeleted,
		Company:   company,
	}
	if err := u.eventRepository.Send(ctx, event); err != nil {
		return err
	}
	return nil
}
