package memory

import "github.com/brandhoej/flat/pkg/event"

type DeferredBus struct {
	bus   event.Bus
	queue fifoQueue[event.Event]
}

func (deferred *DeferredBus) Enqueue(event event.Event) {
	deferred.queue.enqueue(event)
}

func (deferred *DeferredBus) Flush() error {
	for {
		event, found := deferred.queue.dequeue()

		if !found {
			break
		}

		if err := deferred.bus.Publish(event); err != nil {
			return err
		}
	}

	return nil
}

func CreateDeferredBus(bus event.Bus) event.DeferredBus {
	return &DeferredBus{
		bus,
		createFifoQueue[event.Event](),
	}
}
