package signin

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
	"github.com/oklog/ulid/v2"
)

const signInUseCaseID = event.Index(0)

var errCouldNotPublishSignInEvent = errors.New("could not publish sign in event")

type Command struct {
	SessionID string
	UserID    string
}

type Event struct {
	UserID    string
	SessionID string
}

type useCase struct {
	bus event.DeferredBus
}

func (useCase *useCase) guest() cqrs.Command[identity.Guest] {
	return func(guest identity.Guest, context context.Context) error {
		return useCase.perform(context, Command{
			UserID:    guest.ID,
			SessionID: ulid.Make().String(),
		})
	}
}

func (useCase *useCase) perform(
	context context.Context, command Command,
) error {
	signedIn := Event{
		UserID:    command.UserID,
		SessionID: command.SessionID,
	}

	event := event.Create(
		command.UserID,
		signInUseCaseID,
		signedIn,
	)

	useCase.bus.Enqueue(event)

	if err := useCase.bus.Flush(); err != nil {
		return errors.Join(errCouldNotPublishSignInEvent, err)
	}

	return nil
}
