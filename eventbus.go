package eventbus

// Item sent on the eventbus
type Message interface{}

// Receiver of events
type Handler chan Message

// A simplistic pub/sub system.
//
// * All methods are thread/GoRoutine safe. See Close() for some related caveats
// * Slow/laggy handlers should not slow down routines Emit()ing events.
// * All messages will be received in the same order on all handlers.
type EventBus struct {
	// Messages/Events are emitted through this
	input chan Message

	// Avenue for registering handlers
	register   chan Handler
	unregister chan Handler

	// Channel never sends, will be closed when bus is being shut down
	close chan Message
}

// Create an EventBus and start all related GoRoutines.
func CreateEventBus() *EventBus {
	bus := &EventBus{}
	bus.input = make(Handler)
	bus.register = make(chan Handler)
	bus.unregister = make(chan Handler)
	bus.close = make(chan Message)
	go bus.run()
	return bus
}

// Emit a message onto the bus
func (bus *EventBus) Emit(msg Message) {
	bus.input <- msg
}

// Shutdown the event bus.
//
// This will drain all current messages from the queues before closing the handler channels passed to Register().
// After closing an EventBus, any attempt to Emit() a message will cause a panic.
func (bus *EventBus) Close() {
	close(bus.close)
}

// Register a message handler.
func (bus *EventBus) Register(h Handler) {
	bus.register <- h
}

// Unregister an event handler.
func (bus *EventBus) Unregister(h Handler) {
	bus.unregister <- h
}

func (bus *EventBus) run() {

	// A list of handlers and their corresponding shim queues.
	// The shim queue ensures that a slow handler won't bog down the event emitter.
	// Map structure is: handler -> shimQueue(handler)
	handlers := make(map[chan Message]chan<- Message)

	for done := false; !done; {
		select {
		case _, ok := <-bus.close:
			if !ok {
				done = true
			} else {
				// not sure how this happened.
			}

		case msg := <-bus.input:
			for _, shim := range handlers {
				shim <- msg
			}

		case handler := <-bus.register:
			// Register a handler (only if it hasn't been registered yet)
			_, ok := handlers[handler]
			if !ok {
				shim_input := make(chan Message)
				CreateShimQueue(shim_input, handler)
				handlers[handler] = shim_input
			}

		case handler := <-bus.unregister:
			// Unregister a handler
			shim, ok := handlers[handler]
			if ok {
				close(shim)
				delete(handlers, handler)
			}
		}
	}

	// Close all our shim queues
	for _, shim := range handlers {
		close(shim)
	}
}
