package simulation

import (
	"errors"

	"github.com/brandhoej/flat/pkg/event"
	"github.com/brandhoej/flat/pkg/identity"
)

var (
	ErrUseCaseTestFailed  = errors.New("a use case test failed")
	ErrUnknowPositiveTest = errors.New("the input action does not have a correspoding positive test")
	ErrUnknowNegativeTest = errors.New("the input action does not have a correspoding negative test")
)

type (
	UseCaseTest func() error
	EventTest   func(event *event.Event) error
)

type Tester interface {
	Identity() (identity.Identity, error)
	Test(useCaseID event.Index) error
	Handle(event *event.Event) error
	ReceivedUseCaseEvent() (bool, error)
}

type HTTPTester struct {
	tests            HTTPTests
	receviedEvent    bool
	pendingUseCaseID event.Index
	eventError       error
}

func NewHTTPTester() *HTTPTester {
	tests := *newHTTPTests()
	return &HTTPTester{
		tests: tests,
	}
}

func (tester *HTTPTester) Identity() (identity.Identity, error) {
	return identity.Identity{}, nil
}

func (tester *HTTPTester) Test(useCaseID event.Index) error {
	tester.pendingUseCaseID = useCaseID
	if err := tester.tests.useCase(useCaseID); err != nil {
		tester.receviedEvent = false
		return errors.Join(ErrUseCaseTestFailed, err)
	}

	return nil
}

func (tester *HTTPTester) Handle(event *event.Event) error {
	// We received an event but we are not awaiting one.
	if !tester.receviedEvent {
		return nil
	}

	// We are awaiting an event but it was not the use case event.
	if tester.pendingUseCaseID != event.UseCase {
		return nil
	}

	// we received the correct use case event so we are no longer awaiting one.
	tester.receviedEvent = true

	if err := tester.tests.event(event); err != nil {
		tester.eventError = err
	}

	return nil
}

func (tester *HTTPTester) ReceivedUseCaseEvent() (bool, error) {
	return tester.receviedEvent, tester.eventError
}
