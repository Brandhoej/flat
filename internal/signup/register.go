package signup

import (
	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/event/memory"
	"github.com/labstack/echo/v4"
)

type configuration struct {
	store             event.Store
	bus               event.Bus
	defferedBus       event.DeferredBus
	http              *echo.Echo
	emailReceptionist emailReceptionist
	emailRepository   emailRepository
}

type option func(*configuration)

func (config *configuration) configure(options ...option) {
	for _, option := range options {
		option(config)
	}

	if config.store == nil {
		config.store = memory.CreateStore()
	}

	if config.bus == nil {
		config.bus = memory.CreateBus()
	}

	if config.defferedBus == nil {
		config.defferedBus = memory.CreateDeferredBus(
			config.bus,
		)
	}

	if config.emailReceptionist == nil {
		config.emailReceptionist = createMemoryEmailReceptionist()
	}

	if config.emailRepository == nil {
		config.emailRepository = createMemoryEmailRepository()
	}
}

func WithEventStore(store event.Store) option {
	return func(configuration *configuration) {
		configuration.store = store
	}
}

func WithEventBus(bus event.Bus) option {
	return func(configuration *configuration) {
		configuration.bus = bus
	}
}

func WithDeferredEventBus(bus event.DeferredBus) option {
	return func(configuration *configuration) {
		configuration.defferedBus = bus
	}
}

func WithHttp(echo *echo.Echo) option {
	return func(configuration *configuration) {
		configuration.http = echo
	}
}

func Register(options ...option) {
	var configuration configuration
	configuration.configure(options...)

	configuration.bus.Subscribe(&event.ListenSelf{
		Store:   configuration.store,
		UseCase: signUpUseCaseID,
	})

	if configuration.http != nil {
		registerHttp(
			configuration.http,
			configuration.defferedBus,
			configuration.emailReceptionist,
			configuration.emailRepository,
		)
	}
}
