package simulation

import (
	"errors"
	"fmt"

	"github.com/brandhoej/flat/pkg/event"
)

var (
	ErrUnknownHttpTest  = errors.New("unkonw use casetest")
	ErrUnknownEventTest = errors.New("unknown event test")
	ErrTestFailed       = errors.New("test failed")
)

type HTTPTests struct {
	useCases map[event.Index]UseCaseTest
	events   map[event.Index]EventTest
}

func newHTTPTests() *HTTPTests {
	tests := &HTTPTests{
		useCases: make(map[event.Index]UseCaseTest),
		events:   make(map[event.Index]EventTest),
	}

	tests.useCases[signInUseCaseID] = tests.signIn
	tests.events[signInUseCaseID] = tests.signedIn
	tests.useCases[signUpUseCaseID] = tests.signUp
	tests.events[signUpUseCaseID] = tests.signedUp
	tests.useCases[acceptSignUpUseCaseID] = tests.acceptSignUp
	tests.events[acceptSignUpUseCaseID] = tests.acceptedSignUp
	tests.useCases[signOutUseCaseID] = tests.signOut
	tests.events[signOutUseCaseID] = tests.signedOut

	return tests
}

func (tests *HTTPTests) useCase(useCaseID event.Index) error {
	testCase, exists := tests.useCases[useCaseID]
	if !exists {
		return ErrUnknownHttpTest
	}

	if err := testCase(); err != nil {
		return errors.Join(ErrTestFailed, err)
	}

	return nil
}

func (tests *HTTPTests) event(event *event.Event) error {
	testCase, exists := tests.events[event.UseCase]
	if !exists {
		return ErrUnknownEventTest
	}

	if err := testCase(event); err != nil {
		return errors.Join(ErrTestFailed, err)
	}

	return nil
}

func (tests *HTTPTests) signIn() error {
	fmt.Println("Use case test: signIn")
	return nil
}

func (tests *HTTPTests) signedIn(event *event.Event) error {
	fmt.Println("Event test: signedIn")
	return nil
}

func (tests *HTTPTests) signUp() error {
	fmt.Println("Use case test: signUp")
	return nil
}

func (tests *HTTPTests) signedUp(event *event.Event) error {
	fmt.Println("Event test: signedUp")
	return nil
}

func (tests *HTTPTests) acceptSignUp() error {
	fmt.Println("Use case test: acceptSignUp")
	return nil
}

func (tests *HTTPTests) acceptedSignUp(event *event.Event) error {
	fmt.Println("Event test: acceptedSignUp")
	return nil
}

func (tests *HTTPTests) signOut() error {
	fmt.Println("Use case test: signOut")
	return nil
}

func (tests *HTTPTests) signedOut(event *event.Event) error {
	fmt.Println("Event test: signedOut")
	return nil
}
