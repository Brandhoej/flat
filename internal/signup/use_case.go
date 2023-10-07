package signup

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
	"github.com/brandhoej/flat/pkg/pii"
	"github.com/oklog/ulid/v2"
)

const signUpUseCaseID = event.Index(1)

var errCouldNotPublishSignUpEvent = errors.New("could not publish sign in event")

type Command struct {
	ActorID  string
	PII      []byte
	Password string
	Email    string
}

type Event struct {
	ID    string
	Email string
}

type useCase struct {
	bus event.DeferredBus
}

func (useCase *useCase) guest(
	password string,
	email string,
	receptionist emailReceptionist,
	emails emailRepository,
) cqrs.Command[identity.Guest] {
	emailReservationGuard := emailReservationGuard{
		receptionist: receptionist,
	}
	uniqueEmailGuard := uniqueEmailGuard{
		emails: emails,
	}
	passwordStrengthGuard := passwordStrengthGuard{}

	return cqrs.Guard[identity.Guest](
		func(guest identity.Guest, context context.Context) error {
			return useCase.perform(context, Command{
				ActorID:  guest.ID,
				PII:      guest.Actor.PII,
				Password: password,
				Email:    email,
			})
		},
		passwordStrengthGuard.guest(password),
		uniqueEmailGuard.guest(email),
		emailReservationGuard.guest(email),
	)
}

func (useCase *useCase) perform(
	context context.Context, command Command,
) error {
	pii := pii.Create(command.PII)

	signedUp := Event{
		ID:    ulid.Make().String(),
		Email: pii.Encrypt(command.Email),
	}

	event := event.Create(
		command.ActorID,
		signUpUseCaseID,
		signedUp,
		event.WithPII(command.PII),
	)

	useCase.bus.Enqueue(event)

	if err := useCase.bus.Flush(); err != nil {
		return errors.Join(errCouldNotPublishSignUpEvent, err)
	}

	return nil
}
