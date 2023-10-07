package signup

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/identity"
)

var (
	errEmailAlreadyInUse             = errors.New("email is already in use by another")
	errFailedCheckingEmailUniqueness = errors.New("email uniqueness check failed")
)

type uniqueEmailGuard struct {
	emails emailRepository
}

func (emailGuard uniqueEmailGuard) guest(email string) cqrs.Command[identity.Guest] {
	return func(actor identity.Guest, context context.Context) error {
		contains, err := emailGuard.emails.contains(context, email)
		if err != nil {
			return errors.Join(errFailedCheckingEmailUniqueness, err)
		}

		if contains {
			return errEmailAlreadyInUse
		}

		return nil
	}
}
