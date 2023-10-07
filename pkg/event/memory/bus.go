package memory

import "github.com/brandhoej/flat/pkg/event"

type inMemoryBus struct {
	handlers []event.Handler
}

func (bus *inMemoryBus) Publish(event event.Event) error {
	for _, handler := range bus.handlers {
		handler.Handle(&event)
	}
	return nil
}

func (bus *inMemoryBus) Subscribe(handler event.Handler) {
	bus.handlers = append(bus.handlers, handler)
}

func CreateBus() event.Bus {
	return &inMemoryBus{
		handlers: make([]event.Handler, 0),
	}
}
