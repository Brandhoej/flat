package signup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
	"github.com/labstack/echo/v4"
)

type httpController struct {
	signUp            useCase
	emailReceptionist emailReceptionist
	emailRepository   emailRepository
}

func registerHttp(
	echo *echo.Echo,
	bus event.DeferredBus,
	emailReceptionist emailReceptionist,
	emailRepository emailRepository,
) {
	controller := httpController{
		useCase{
			bus,
		},
		emailReceptionist,
		emailRepository,
	}
	echo.POST("users", controller.endpoint)
}

type requestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseBody struct {
	UserID string `json:"user_id"`
}

func (controller *httpController) endpoint(c echo.Context) error {
	var requestBody requestBody
	if err := json.NewDecoder(c.Request().Body).Decode(&requestBody); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		fmt.Println("signup", err)
		return err
	}

	id, err := identity.FromAccessCookie(c)
	if err != nil {
		if errors.Is(err, identity.ErrNoAccessCookie) {
			c.Response().WriteHeader(http.StatusUnauthorized)
		} else {
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("signup", err)
		return err
	}

	context := context.Background()
	if canBeGuest, err := id.AsGuest(controller.signUp.guest(
		requestBody.Password,
		requestBody.Email,
		controller.emailReceptionist,
		controller.emailRepository,
	), context); canBeGuest && err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		fmt.Println("signup", err)
		return err
	}

	responseBody := responseBody{
		UserID: id.Subject,
	}

	c.Response().WriteHeader(http.StatusCreated)
	json.NewEncoder(c.Response()).Encode(responseBody)

	return nil
}
