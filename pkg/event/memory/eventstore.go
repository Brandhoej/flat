package memory

import (
	"errors"

	"github.com/brandhoej/flat/pkg/event"
)

var (
	ErrEventsStartsAtZero        = errors.New("aggregate events must start with index 0")
	ErrEventsIncrementsByOne     = errors.New("events must be exactly one increment above the current of the aggregate")
	ErrEventTimesMustBeAscending = errors.New("the creation time of events must be in ascending order")
	ErrAggreagteHasNoEvents      = errors.New("the aggregate has no events")
)

type inMemoryEventStore struct {
	events map[string][]event.Event
}

func CreateStore() event.Store {
	return &inMemoryEventStore{
		make(map[string][]event.Event),
	}
}

func (store *inMemoryEventStore) Query(
	aggregateID string,
	params event.StoreQueryParams,
) ([]event.Event, error) {
	events, exists := store.events[aggregateID]
	if !exists {
		return nil, ErrAggreagteHasNoEvents
	}

	beginning, end := params.Start, params.End

	return events[beginning:end], nil
}

func (store *inMemoryEventStore) Store(event event.Event) error {
	store.events[event.Actor] = append(store.events[event.Actor], event)

	return nil
}
