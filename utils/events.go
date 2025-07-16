package utils

type Event[T any] interface {
	Name() string
	Payload() T
}

type Subscriber[T Event[V], V any] interface {
	Handle(event T)
}

type Emitter[T Event[V], V any] interface {
	Emit(event T)
	Subscribe(eventName string, subscriber Subscriber[T, V])
}

type emitter[T Event[V], V any] struct {
	subscribers map[string][]Subscriber[T, V]
}

func (e *emitter[T, V]) Emit(event T) {
	if handlers, ok := e.subscribers[event.Name()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
	}
}

func (e *emitter[T, V]) Subscribe(eventName string, subscriber Subscriber[T, V]) {
	e.subscribers[eventName] = append(e.subscribers[eventName], subscriber)
}

// NewEmitter creates a new Emitter instance for the specified event type.
func NewEmitter[T Event[V], V any]() Emitter[T, V] {
	return &emitter[T, V]{
		subscribers: make(map[string][]Subscriber[T, V]),
	}
}
