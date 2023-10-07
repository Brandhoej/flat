package signin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
	"github.com/brandhoej/flat/pkg/jwt"
	"github.com/labstack/echo/v4"
)

type httpController struct {
	signIn useCase
}

func registerHttp(
	echo *echo.Echo,
	bus event.DeferredBus,
) {
	controller := httpController{
		useCase{
			bus,
		},
	}
	echo.POST("auth/signin", controller.endpoint)
}

func (controller *httpController) endpoint(c echo.Context) error {
	// Try identifying the user from the access token.
	id, err := identity.FromAccessCookie(c)
	if err != nil && !errors.Is(err, identity.ErrNoAccessCookie) {
		c.Response().WriteHeader(http.StatusInternalServerError)
		fmt.Println("signin", err)
		return err
	}

	// The user is unidentifiable - create a guest account.
	// This is done to preserve identity of the actions that is take by every user.
	if errors.Is(err, identity.ErrNoAccessCookie) {
		id = identity.NewGuest()
	}

	// Try signing in the user as a guest.
	context := context.Background()
	if isGuest, err := id.AsGuest(
		controller.signIn.guest(), context,
	); isGuest && err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		fmt.Println("signin", err)
		return err
	}

	// We could sign in so we create a JWT access token to identify the user in the other use cases.
	cookie, err := jwt.CreateAccessCookie(
		jwt.PublicClaims{
			Role: id.Role,
			PII:  id.Pii,
		},
		jwt.ForActor(id.Subject),
		jwt.WithExpiration(time.Now().Add(time.Hour)),
	)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		fmt.Println("signin", err)
		return err
	}
	c.SetCookie(&cookie)

	return nil
}
