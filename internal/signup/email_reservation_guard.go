package signup

import (
	"context"
	"errors"
	"time"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/identity"
)

var errEmailCouldNotBeReserved = errors.New("could not reserve email")

type emailReservationGuard struct {
	receptionist emailReceptionist
}

func (emailGuard emailReservationGuard) guest(email string) cqrs.Command[identity.Guest] {
	return func(actor identity.Guest, context context.Context) error {
		expiration := time.Now().Add(time.Hour)
		if err := emailGuard.receptionist.reserve(email, expiration); err != nil {
			return errEmailCouldNotBeReserved
		}

		return nil
	}
}
