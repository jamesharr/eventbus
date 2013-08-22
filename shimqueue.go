package eventbus

type Message interface{}

type shimQueue struct {
	// Listener channel
	input <-chan Message

	// Emitter channel
	output chan<- Message

	// Queue
	queue []Message
}

// Create a queue shim
func CreateShimQueueFromInput(input <-chan Message) (output chan<- Message) {
	output = make(chan<- Message)
	CreateShimQueue(input, output)
	return
}

// Create a queue shim
func CreateShimQueue(input <-chan Message, output chan<- Message) {
	var shim shimQueue
	shim.input = input
	shim.output = output
	shim.queue = make([]Message, 0)
	go shim.run()
}

func (shim *shimQueue) run() {
	// Main loop -- recv+queue messages, send+deque messages
	done := false
	for !done {
		if len(shim.queue) > 0 {
			// Event loop with items in queue
			select {

			case msg, ok := <-shim.input:
				// Activity on input queue
				if ok {
					shim.queue = append(shim.queue, msg)
				} else {
					done = true
				}

			case shim.output <- shim.queue[0]:
				// Item sent on receive queue
				shim.queue = shim.queue[1:]
			}
		} else {
			// Event loop without queue processing
			select {

			case msg, ok := <-shim.input:
				// Activity on input queue
				if ok {
					shim.queue = append(shim.queue, msg)
				} else {
					done = true
				}
			}
		}
	}

	// Drain queue
	for _, msg := range shim.queue {
		shim.output <- msg
	}

	// Tear down some things
	close(shim.output)
	shim.queue = []Message{}
}
