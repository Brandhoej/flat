package message

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
)

const messageSignUpConfirmationUseCaseID = event.Index(2)

var (
	errCouldNotPublishMessagesEvent = errors.New("could not publish message event")
	errCouldNotMessage              = errors.New("messenger failed messaging")
)

type Command struct {
	MessageID string
	UserID    string
	Message   string
}

type Event struct {
	MessageID string
	UserID    string
	Message   string
}

type useCase struct {
	bus       event.DeferredBus
	messenger messenger
}

func (useCase *useCase) actor(messageID string, message string) cqrs.Command[identity.Actor] {
	return func(actor identity.Actor, context context.Context) error {
		return useCase.perform(context, Command{
			UserID:    actor.ID,
			MessageID: messageID,
			Message:   message,
		})
	}
}

func (useCase *useCase) perform(
	context context.Context, command Command,
) error {
	if err := useCase.messenger.message(command, context); err != nil {
		return errors.Join(errCouldNotMessage, err)
	}

	messaged := Event{
		UserID:    command.UserID,
		MessageID: command.MessageID,
		Message:   command.Message,
	}

	event := event.Create(
		command.UserID,
		messageSignUpConfirmationUseCaseID,
		messaged,
	)

	useCase.bus.Enqueue(event)

	if err := useCase.bus.Flush(); err != nil {
		return errors.Join(errCouldNotPublishMessagesEvent, err)
	}

	return nil
}
