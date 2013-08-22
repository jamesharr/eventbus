package eventbus

type Message interface{}
type Handler chan Message

type EventBus struct {
	// Messages/Events are emitted through this
	input chan Message

	// Avenue for registering handlers
	register   chan Handler
	unregister chan Handler

	// Channel never sends, will be closed when bus is being shut down
	close chan Message
}

func CreateEventBus() *EventBus {
	bus := &EventBus{}
	bus.input = make(Handler)
	bus.register = make(chan Handler)
	bus.unregister = make(chan Handler)
	bus.close = make(chan Message)
	go bus.run()
	return bus
}

func (bus *EventBus) Emit(msg Message) {
	bus.input <- msg
}

func (bus *EventBus) Close() {
	close(bus.close)
}

func (bus *EventBus) Register(h Handler) {
	bus.register <- h
}

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
