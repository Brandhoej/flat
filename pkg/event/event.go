package event

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

type (
	Index uint
)

var (
	ErrNoData        = errors.New("event does not have any data")
	ErrMarshalling   = errors.New("event data could not be marshalled")
	ErrUnmarshalling = errors.New("event data could not be unmarshalled")
	ErrDataIsNil     = errors.New("cannot unmarshal data to a nil reference")
)

type Event struct {
	// The global identification of the actor which is responsible for producing this event.
	// For groups the actor itself is the group which constitutes the set of actors.
	// The point is that the actor is accountable for the production of the event.
	Actor string

	// The global identification of the use case was was executed by
	// the actor and led to the production of the event.
	// Example: Sign up, sign in, sign out.
	UseCase Index

	// Used for event versioning if it is strictly impossible for the event evolution to stay backward compatible.
	Version Index

	// The creation time of the event used for event ordering and in/out order determination and correction.
	Time time.Time

	// The event data.
	Data interface{}

	// Additional metadata which can be attached but is not strictly required.
	Metadata map[string]interface{}
}

type EventOption func(event *Event)

func ForVersion(version Index) EventOption {
	return func(event *Event) {
		event.Version = version
	}
}

func WithMetadata(key string, data interface{}) EventOption {
	return func(event *Event) {
		event.Metadata[key] = data
	}
}

func WithPII(pii []byte) EventOption {
	return WithMetadata("pii", pii)
}

func Create(
	actor string,
	useCase Index,
	data interface{},
	options ...EventOption,
) Event {
	event := Event{
		Actor:    actor,
		UseCase:  useCase,
		Version:  Index(0),
		Data:     data,
		Time:     time.Now(),
		Metadata: map[string]interface{}{},
	}

	for _, option := range options {
		option(&event)
	}

	return event
}

func (event Event) Reason() (string, error) {
	if event.Data == nil {
		return "", ErrNoData
	}

	return reflect.TypeOf(event.Data).Elem().Name(), nil
}

func EventData[T any](event Event) (T, error) {
	var data T
	err := event.DataAs(&data)
	return data, err
}

func (event Event) DataAs(data interface{}) error {
	if data == nil {
		return ErrDataIsNil
	}

	bytes, err := json.Marshal(event.Data)
	if err != nil {
		return errors.Join(ErrMarshalling, err)
	}

	if err := json.Unmarshal(bytes, data); err != nil {
		return errors.Join(ErrUnmarshalling, err)
	}

	return nil
}

func Metadata[T any](key string, event Event) (T, bool) {
	var value T

	entry, exists := event.Metadata[key]
	if !exists {
		return value, false
	}

	value, ok := entry.(T)
	if !ok {
		return value, false
	}

	return value, true
}
