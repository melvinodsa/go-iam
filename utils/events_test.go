package utils

import (
	"context"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of Event interface
type MockEvent struct {
	mock.Mock
}

func (m *MockEvent) Name() goiamuniverse.Event {
	args := m.Called()
	return args.Get(0).(goiamuniverse.Event)
}

func (m *MockEvent) Payload() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockEvent) Metadata() sdk.Metadata {
	args := m.Called()
	return args.Get(0).(sdk.Metadata)
}

func (m *MockEvent) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// Simple concrete Event implementation for testing
type TestEvent struct {
	name     goiamuniverse.Event
	payload  string
	metadata sdk.Metadata
	ctx      context.Context
}

func (e *TestEvent) Name() goiamuniverse.Event {
	return e.name
}

func (e *TestEvent) Payload() string {
	return e.payload
}

func (e *TestEvent) Metadata() sdk.Metadata {
	return e.metadata
}

func (e *TestEvent) Context() context.Context {
	return e.ctx
}

// Mock Subscriber implementation
type MockSubscriber struct {
	mock.Mock
	events []Event[string]
}

func (m *MockSubscriber) HandleEvent(event Event[string]) {
	m.Called(event)
	m.events = append(m.events, event)
}

func TestNewEmitter(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()

	assert.NotNil(t, emitter)

	// Test that it implements the interface correctly by calling interface methods
	assert.NotPanics(t, func() {
		subscriber := &MockSubscriber{}
		emitter.Subscribe("test", subscriber)
	})
}

func TestEmitter_Subscribe(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")

	// Subscribe to event - should not panic
	assert.NotPanics(t, func() {
		emitter.Subscribe(eventName, subscriber)
	})
}

func TestEmitter_Subscribe_MultipleSubscribers(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber1 := &MockSubscriber{}
	subscriber2 := &MockSubscriber{}
	subscriber3 := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")

	// Subscribe multiple subscribers to same event - should not panic
	assert.NotPanics(t, func() {
		emitter.Subscribe(eventName, subscriber1)
		emitter.Subscribe(eventName, subscriber2)
		emitter.Subscribe(eventName, subscriber3)
	})
}

func TestEmitter_Subscribe_DifferentEvents(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber1 := &MockSubscriber{}
	subscriber2 := &MockSubscriber{}

	event1 := goiamuniverse.Event("event-1")
	event2 := goiamuniverse.Event("event-2")

	// Subscribe to different events - should not panic
	assert.NotPanics(t, func() {
		emitter.Subscribe(event1, subscriber1)
		emitter.Subscribe(event2, subscriber2)
	})
}

func TestEmitter_Emit_WithSubscribers(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit
	emitter.Subscribe(eventName, subscriber)
	emitter.Emit(event)

	// Verify subscriber was called
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
	assert.Equal(t, event, subscriber.events[0])
}

func TestEmitter_Emit_NoSubscribers(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      context.Background(),
	}

	// Emit event with no subscribers - should not panic
	assert.NotPanics(t, func() {
		emitter.Emit(event)
	})
}

func TestEmitter_Emit_MultipleSubscribers(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber1 := &MockSubscriber{}
	subscriber2 := &MockSubscriber{}
	subscriber3 := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber1.On("HandleEvent", event).Return()
	subscriber2.On("HandleEvent", event).Return()
	subscriber3.On("HandleEvent", event).Return()

	// Subscribe all and emit
	emitter.Subscribe(eventName, subscriber1)
	emitter.Subscribe(eventName, subscriber2)
	emitter.Subscribe(eventName, subscriber3)
	emitter.Emit(event)

	// Verify all subscribers were called
	subscriber1.AssertExpectations(t)
	subscriber2.AssertExpectations(t)
	subscriber3.AssertExpectations(t)

	assert.Len(t, subscriber1.events, 1)
	assert.Len(t, subscriber2.events, 1)
	assert.Len(t, subscriber3.events, 1)
}

func TestEmitter_Emit_WrongEventName(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	subscribeEventName := goiamuniverse.Event("subscribe-event")
	emitEventName := goiamuniverse.Event("emit-event")

	event := &TestEvent{
		name:     emitEventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      context.Background(),
	}

	// Subscribe to one event but emit different event
	emitter.Subscribe(subscribeEventName, subscriber)
	emitter.Emit(event)

	// Subscriber should not be called
	subscriber.AssertNotCalled(t, "HandleEvent", mock.Anything)
	assert.Len(t, subscriber.events, 0)
}

// Test context key type to avoid lint warning
type testContextKey string

func TestEvent_Interface_Methods(t *testing.T) {
	eventName := goiamuniverse.Event("test-event")
	payload := "test-payload"
	metadata := sdk.Metadata{
		ProjectIds: []string{"project1", "project2"},
	}
	ctx := context.WithValue(context.Background(), testContextKey("test-key"), "test-value")

	event := &TestEvent{
		name:     eventName,
		payload:  payload,
		metadata: metadata,
		ctx:      ctx,
	}

	// Test all interface methods
	assert.Equal(t, eventName, event.Name())
	assert.Equal(t, payload, event.Payload())
	assert.Equal(t, metadata, event.Metadata())
	assert.Equal(t, ctx, event.Context())
}

func TestEmitter_ComplexScenario(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()

	// Create multiple subscribers for different events
	userSubscriber := &MockSubscriber{}
	roleSubscriber := &MockSubscriber{}
	allEventsSubscriber := &MockSubscriber{}

	userEvent := goiamuniverse.Event("user:created")
	roleEvent := goiamuniverse.Event("role:updated")

	userEventData := &TestEvent{
		name:     userEvent,
		payload:  "user-data",
		metadata: sdk.Metadata{ProjectIds: []string{"user-project"}},
		ctx:      context.Background(),
	}

	roleEventData := &TestEvent{
		name:     roleEvent,
		payload:  "role-data",
		metadata: sdk.Metadata{ProjectIds: []string{"role-project"}},
		ctx:      context.Background(),
	}

	// Set up expectations
	userSubscriber.On("HandleEvent", userEventData).Return()
	roleSubscriber.On("HandleEvent", roleEventData).Return()
	allEventsSubscriber.On("HandleEvent", userEventData).Return()
	allEventsSubscriber.On("HandleEvent", roleEventData).Return()

	// Subscribe
	emitter.Subscribe(userEvent, userSubscriber)
	emitter.Subscribe(roleEvent, roleSubscriber)
	emitter.Subscribe(userEvent, allEventsSubscriber)
	emitter.Subscribe(roleEvent, allEventsSubscriber)

	// Emit events
	emitter.Emit(userEventData)
	emitter.Emit(roleEventData)

	// Verify expectations
	userSubscriber.AssertExpectations(t)
	roleSubscriber.AssertExpectations(t)
	allEventsSubscriber.AssertExpectations(t)

	// Verify event counts
	assert.Len(t, userSubscriber.events, 1)
	assert.Len(t, roleSubscriber.events, 1)
	assert.Len(t, allEventsSubscriber.events, 2)
}

func TestEmitter_EmptyEventName(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	emptyEventName := goiamuniverse.Event("")
	event := &TestEvent{
		name:     emptyEventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit with empty event name
	emitter.Subscribe(emptyEventName, subscriber)
	emitter.Emit(event)

	// Should work normally
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
}

func TestEmitter_NilContext(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{},
		ctx:      nil, // nil context
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit
	emitter.Subscribe(eventName, subscriber)
	emitter.Emit(event)

	// Should work normally
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
	assert.Nil(t, subscriber.events[0].Context())
}

func TestEmitter_EmptyMetadata(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{}, // empty metadata
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit
	emitter.Subscribe(eventName, subscriber)
	emitter.Emit(event)

	// Should work normally
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
	assert.Empty(t, subscriber.events[0].Metadata())
}

func TestEmitter_NilMetadata(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("test-event")
	event := &TestEvent{
		name:     eventName,
		payload:  "test-payload",
		metadata: sdk.Metadata{}, // empty metadata instead of nil
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit
	emitter.Subscribe(eventName, subscriber)
	emitter.Emit(event)

	// Should work normally
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
	assert.Empty(t, subscriber.events[0].Metadata().ProjectIds)
}

// Test with different payload types
type ComplexPayload struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

type ComplexEvent struct {
	name     goiamuniverse.Event
	payload  ComplexPayload
	metadata sdk.Metadata
	ctx      context.Context
}

func (e *ComplexEvent) Name() goiamuniverse.Event {
	return e.name
}

func (e *ComplexEvent) Payload() ComplexPayload {
	return e.payload
}

func (e *ComplexEvent) Metadata() sdk.Metadata {
	return e.metadata
}

func (e *ComplexEvent) Context() context.Context {
	return e.ctx
}

type ComplexSubscriber struct {
	mock.Mock
	events []Event[ComplexPayload]
}

func (s *ComplexSubscriber) HandleEvent(event Event[ComplexPayload]) {
	s.Called(event)
	s.events = append(s.events, event)
}

func TestEmitter_ComplexPayload(t *testing.T) {
	emitter := NewEmitter[Event[ComplexPayload], ComplexPayload]()
	subscriber := &ComplexSubscriber{}

	eventName := goiamuniverse.Event("complex-event")
	payload := ComplexPayload{
		ID: "test-id",
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		},
	}

	event := &ComplexEvent{
		name:     eventName,
		payload:  payload,
		metadata: sdk.Metadata{ProjectIds: []string{"complex-project"}},
		ctx:      context.Background(),
	}

	// Set up mock expectations
	subscriber.On("HandleEvent", event).Return()

	// Subscribe and emit
	emitter.Subscribe(eventName, subscriber)
	emitter.Emit(event)

	// Verify
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 1)
	assert.Equal(t, payload.ID, subscriber.events[0].Payload().ID)
	assert.Equal(t, payload.Data, subscriber.events[0].Payload().Data)
}

func TestEmitter_ConcurrentAccess(t *testing.T) {
	emitter := NewEmitter[Event[string], string]()
	subscriber := &MockSubscriber{}

	eventName := goiamuniverse.Event("concurrent-event")

	// Set up mock expectations for multiple events
	subscriber.On("HandleEvent", mock.AnythingOfType("*utils.TestEvent")).Return().Times(10)

	// Subscribe
	emitter.Subscribe(eventName, subscriber)

	// Emit multiple events concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			event := &TestEvent{
				name:     eventName,
				payload:  "payload-" + string(rune(index)),
				metadata: sdk.Metadata{},
				ctx:      context.Background(),
			}
			emitter.Emit(event)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all events were handled
	subscriber.AssertExpectations(t)
	assert.Len(t, subscriber.events, 10)
}
