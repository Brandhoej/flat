package message

import (
	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/event/memory"
	"github.com/labstack/echo/v4"
)

const SignUpUseCaseID = event.Index(1)

type configuration struct {
	store          event.Store
	bus            event.Bus
	defferedBus    event.DeferredBus
	http           *echo.Echo
	listenToSignUp bool
	test           bool
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

func ListenToSignUp() option {
	return func(configuration *configuration) {
		configuration.listenToSignUp = true
	}
}

func Register(options ...option) {
	var configuration configuration
	configuration.configure(options...)

	configuration.bus.Subscribe(&event.ListenSelf{
		Store:   configuration.store,
		UseCase: messageSignUpConfirmationUseCaseID,
	})

	if configuration.listenToSignUp {
		var messenger messenger
		if configuration.test {
			messenger = &mockMessenger{}
		} else {
			messenger = &emailMessenger{}
		}

		SignUpListener := event.SubscribeToUseCase(
			SignUpUseCaseID,
			&SignUpListener{
				message: useCase{
					bus:       configuration.defferedBus,
					messenger: messenger,
				},
			},
		)
		configuration.bus.Subscribe(SignUpListener)
	}
}
