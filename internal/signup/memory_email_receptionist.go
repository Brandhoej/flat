package signup

import (
	"errors"
	"time"
)

var errEmailAlreadyReserved = errors.New("email already reserved")

type memoryEmailReceptionist struct {
	reservations map[string]time.Time
}

func (receptionist *memoryEmailReceptionist) reserve(email string, expiration time.Time) error {
	expires_at, exists := receptionist.reservations[email]

	if exists && expires_at.After(time.Now()) {
		return errEmailAlreadyReserved
	}

	receptionist.reservations[email] = expiration

	return nil
}

func createMemoryEmailReceptionist() emailReceptionist {
	return &memoryEmailReceptionist{
		make(map[string]time.Time),
	}
}
