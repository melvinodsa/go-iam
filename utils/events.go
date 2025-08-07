package utils

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type Event[T any] interface {
	Name() goiamuniverse.Event
	Payload() T
	Metadata() sdk.Metadata
	Context() context.Context
}

type Subscriber[T Event[V], V any] interface {
	HandleEvent(event T)
}

type Emitter[T Event[V], V any] interface {
	Emit(event T)
	Subscribe(eventName goiamuniverse.Event, subscriber Subscriber[T, V])
}

type emitter[T Event[V], V any] struct {
	subscribers map[goiamuniverse.Event][]Subscriber[T, V]
}

func (e *emitter[T, V]) Emit(event T) {
	if handlers, ok := e.subscribers[event.Name()]; ok {
		for _, handler := range handlers {
			handler.HandleEvent(event)
		}
	}
}

func (e *emitter[T, V]) Subscribe(eventName goiamuniverse.Event, subscriber Subscriber[T, V]) {
	e.subscribers[eventName] = append(e.subscribers[eventName], subscriber)
}

// NewEmitter creates a new Emitter instance for the specified event type.
func NewEmitter[T Event[V], V any]() Emitter[T, V] {
	return &emitter[T, V]{
		subscribers: make(map[goiamuniverse.Event][]Subscriber[T, V]),
	}
}
