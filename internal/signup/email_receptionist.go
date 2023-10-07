package signup

import "time"

type emailReceptionist interface {
	reserve(email string, expiration time.Time) error
}
