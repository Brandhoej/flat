package message

import (
	"context"
	"errors"
	"fmt"

	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
	"github.com/brandhoej/flat/pkg/pii"
	"github.com/oklog/ulid/v2"
)

var ErrMissingPIIKey = errors.New("the event did not have the pii key in metadata")

type SignUpListener struct {
	message useCase
}

type signUpEvent struct {
	ID    string
	Email string
}

func (handler *SignUpListener) Handle(signedUp *event.Event) error {
	data, err := event.EventData[signUpEvent](*signedUp)
	if err != nil {
		return err
	}

	piiKey, ok := event.Metadata[[]byte]("pii", *signedUp)
	if !ok {
		return ErrMissingPIIKey
	}

	actor := identity.Actor{
		ID:  signedUp.Actor,
		PII: piiKey,
	}

	// FIXME: Remove this as it is just a test of the pii encrypt/decrypt.
	fmt.Println(data.ID, pii.Create(actor.PII).Decrypt(data.Email))

	context := context.Background()
	message := handler.message.actor(
		ulid.Make().String(),
		"Please press YES to confirm your account signup with ID"+data.ID,
	)
	return message(actor, context)
}
