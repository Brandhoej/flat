package signup

import (
	"context"
	"errors"
)

var (
	errEmailAlreadyStored             = errors.New("email is already stored")
	errCouldNotRemoveNonExistingEmail = errors.New("could not remove email as it is not stored")
)

type memoryEmailRepository struct {
	emails map[string]*any
}

func (repository *memoryEmailRepository) store(context context.Context, email string) error {
	if exists, _ := repository.contains(context, email); exists {
		return errEmailAlreadyStored
	}

	repository.emails[email] = nil

	return nil
}

func (repository *memoryEmailRepository) contains(context context.Context, email string) (bool, error) {
	_, exists := repository.emails[email]

	return exists, nil
}

func (repository *memoryEmailRepository) remove(context context.Context, email string) error {
	if exists, _ := repository.contains(context, email); exists {
		return errCouldNotRemoveNonExistingEmail
	}

	delete(repository.emails, email)

	return nil
}

func createMemoryEmailRepository() emailRepository {
	return &memoryEmailRepository{
		map[string]*any{},
	}
}
