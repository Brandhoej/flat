package event

type DeferredBus interface {
	Enqueue(event Event)
	Flush() error
}
