package cqrs

import (
	"context"
	"errors"
)

type (
	Command[A any]      func(actor A, context context.Context) error
	Query[A any, M any] func(actor A, context context.Context) (M, error)
	MiddleWare[A any]   func(actor A, next func(), context context.Context)
)

var (
	ErrGuardError = errors.New("a command guard encountered an error")
	ErrReadError  = errors.New("a read encountered an error")
)

func Guard[A any](
	command Command[A],
	guards ...Command[A],
) Command[A] {
	return func(actor A, context context.Context) error {
		for idx := range guards {
			if context.Err() != nil {
				return errors.Join(ErrGuardError, context.Err())
			}

			if err := guards[idx](actor, context); err != nil {
				return errors.Join(ErrGuardError, err)
			}
		}
		return command(actor, context)
	}
}

func Middle[A any](
	command Command[A],
	middlewares ...MiddleWare[A],
) Command[A] {
	var err error
	var next func(idx int, actor A, context context.Context) func()
	next = func(idx int, actor A, context context.Context) func() {
		return func() {
			if context.Err() != nil {
				return
			}

			lastMiddleWare := idx == len(middlewares)-1
			if lastMiddleWare {
				err = command(actor, context)
			} else {
				middlewares[idx](actor, next(idx+1, actor, context), context)
			}
		}
	}

	return func(actor A, context context.Context) error {
		for idx := range middlewares {
			middlewares[idx](actor, next(idx, actor, context), context)
		}
		return err
	}
}

func Read[A, M any](
	output *M, query Query[A, M],
) Command[A] {
	return func(guest A, context context.Context) error {
		if context.Err() != nil {
			return errors.Join(ErrReadError, context.Err())
		}

		result, err := query(guest, context)
		*output = result
		return err
	}
}
