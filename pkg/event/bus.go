package event

type Handler interface {
	Handle(event *Event) error
}

type Bus interface {
	Publish(event Event) error
	Subscribe(handler Handler)
}

type UseCaseSubscription struct {
	useCase Index
	handler Handler
}

func (useCaseSubscription *UseCaseSubscription) Handle(event *Event) error {
	if event.UseCase == useCaseSubscription.useCase {
		return useCaseSubscription.handler.Handle(event)
	}

	return nil
}

func SubscribeToUseCase(useCase Index, handler Handler) Handler {
	return &UseCaseSubscription{
		useCase: useCase,
		handler: handler,
	}
}
