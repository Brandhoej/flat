package event

type StoreQueryParams struct {
	Start Index
	End   Index
}

type Store interface {
	Query(aggregateID string, params StoreQueryParams) ([]Event, error)
	Store(event Event) error
}
