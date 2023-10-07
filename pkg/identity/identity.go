package identity

import (
	"context"
	"errors"

	"github.com/brandhoej/flat/pkg/cqrs"
	"github.com/brandhoej/flat/pkg/jwt"
	"github.com/brandhoej/flat/pkg/pii"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

var (
	ErrNoAccessCookie      = errors.New("an identity could not be established as there are no access cookie")
	ErrParsingAccessCookie = errors.New("an identity could not be established as there was an error parsing the access cookie")
)

type Actor struct {
	ID  string
	PII []byte
}

type Guest struct {
	Actor
}

type Identity struct {
	Subject string
	Role    string
	Pii     []byte
}

func NewGuest() Identity {
	return Identity{
		Subject: ulid.Make().String(),
		Role:    "guest",
		Pii:     pii.NewKey(),
	}
}

func (guest Guest) Identity() Identity {
	return Identity{
		Subject: guest.ID,
		Role:    "guest",
		Pii:     guest.PII,
	}
}

func FromAccessCookie(context echo.Context) (Identity, error) {
	var identity Identity

	cookie, err := context.Cookie("access")
	if err != nil {
		return identity, errors.Join(ErrNoAccessCookie, err)
	}
	if cookie == nil {
		return identity, ErrNoAccessCookie
	}

	accessToken, err := jwt.Parse(cookie.Value)
	if err != nil {
		return identity, ErrParsingAccessCookie
	}

	identity.Role = accessToken.Role
	identity.Subject = accessToken.Subject
	identity.Pii = accessToken.Pii

	return identity, nil
}

func (id Identity) Guest() Guest {
	return Guest{
		Actor: id.Actor(),
	}
}

func (id Identity) Actor() Actor {
	return Actor{
		ID:  id.Subject,
		PII: id.Pii,
	}
}

func (id Identity) IsGuest() bool {
	return id.Role == "guest"
}

func (id Identity) CanBeGuest() bool {
	return true
}

func (id Identity) AsGuest(
	command cqrs.Command[Guest], context context.Context,
) (bool, error) {
	if !id.CanBeGuest() {
		return false, nil
	}

	return true, command(id.Guest(), context)
}
