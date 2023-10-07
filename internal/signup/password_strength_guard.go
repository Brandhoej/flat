package signup

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/identity"
)

var errPasswordNotStrongEnough = errors.New("password is not strong enought to e used")

type passwordStrengthGuard struct{}

func (passwordGuard passwordStrengthGuard) guest(password string) cqrs.Command[identity.Guest] {
	return func(actor identity.Guest, context context.Context) error {
		if len(password) < 4 {
			return errPasswordNotStrongEnough
		}

		return nil
	}
}
