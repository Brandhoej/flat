package message

import "context"

type messenger interface {
	message(command Command, context context.Context) error
}

type emailMessenger struct{}

func (emailMessenger *emailMessenger) message(command Command, context context.Context) error {
	return nil
}

type mockMessenger struct{}

func (emailMessenger *mockMessenger) message(command Command, context context.Context) error {
	return nil
}
