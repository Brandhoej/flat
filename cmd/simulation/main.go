package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/brandhoej/flat/internal/message"
	"github.com/brandhoej/flat/internal/signin"
	"github.com/brandhoej/flat/internal/signup"
	"github.com/brandhoej/flat/pkg/event/memory"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/publicsuffix"
)

// Outputs: User commands.
// Input: Use case events.

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

	// TODO: Simulation should be various listens to events.

	client := newClient()
	client.signInAsGuest()
	client.signUp("andreasbrandhoej@hotmail.com", "P@ssW0rd")

	time.Sleep(time.Hour)
}

type client struct {
	httpClient http.Client
}

func newClient() client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	cookies, err := cookiejar.New(&options)
	if err != nil {
		panic(err)
	}

	return client{
		httpClient: http.Client{
			Jar: cookies,
		},
	}
}

func (client *client) signInAsGuest() {
	client.post("/auth/signin")
}

type signUpRequestBody struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type signUpResponseBody struct {
	UserID string `json:"user_id,omitempty"`
}

func (client *client) signUp(email string, password string) {
	request := signUpRequestBody{
		Email:    email,
		Password: password,
	}

	var response signUpResponseBody
	client.postJson("/users", request, &response)

	fmt.Println(response)
}

func (client *client) post(url string) {
	client.httpClient.Post("http://localhost:8081"+url, "application/json", nil)
}

func (client *client) postJson(url string, requestBody any, responseBody any) {
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	response, err := client.httpClient.Post(
		"http://localhost:8081"+url, "application/json", bytes.NewBuffer(jsonRequestBody),
	)
	if err != nil {
		panic(err)
	}

	json.NewDecoder(response.Body).Decode(&responseBody)
}
