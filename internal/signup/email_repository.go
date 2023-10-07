package signup

import "context"

type emailRepository interface {
	store(context context.Context, email string) error
	contains(context context.Context, email string) (bool, error)
	remove(context context.Context, email string) error
}
