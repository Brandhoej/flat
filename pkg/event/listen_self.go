package event

type ListenSelf struct {
	Store   Store
	UseCase Index
}

func (handler *ListenSelf) Handle(event *Event) error {
	if event.UseCase != handler.UseCase {
		return nil
	}

	return handler.Store.Store(*event)
}
