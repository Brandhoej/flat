package main

import (
	"time"

	"github.com/brandhoej/flat/internal/message"
	"github.com/brandhoej/flat/internal/signin"
	"github.com/brandhoej/flat/internal/signup"
	"github.com/brandhoej/flat/internal/simulation"
	"github.com/brandhoej/flat/pkg/event/memory"
	"github.com/labstack/echo/v4"
)

func main() {
	bus := memory.CreateBus()
	echo := echo.New()

	signin.Register(
		signin.WithEventBus(bus),
		signin.WithHttp(echo),
	)
	signup.Register(
		signup.WithEventBus(bus),
		signup.WithHttp(echo),
	)
	message.Register(
		message.WithEventBus(bus),
		message.ListenToSignUp(),
	)

	go echo.Start(":8081")

	time.Sleep(5 * time.Second)

	simulation := simulation.Start(bus)
	simulation.Step()
}
