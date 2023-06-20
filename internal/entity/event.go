package entity

type EventOperation string

const (
	EventTypeCreated EventOperation = "created"
	EventTypeUpdated EventOperation = "updated"
	EventTypeDeleted EventOperation = "deleted"
)

type Event struct {
	Operation EventOperation `json:"operation"`
	Company   *Company       `json:"company,omitempty"`
}
