package eventbus

// Simple expanding queue structure that can be thrown between a sender/receiver.
// This can effectively emulate a channel with unrestricted buffering.
type shimQueue struct {
	// Listener channel
	input <-chan Message

	// Emitter channel
	output chan<- Message
}

// Create a queue shim
func CreateShimQueue(input <-chan Message, output chan<- Message) {
	var shim shimQueue
	shim.input = input
	shim.output = output
	go shim.run()
}

func (shim *shimQueue) run() {
	// Main loop -- recv+queue messages, send+deque messages
	queue := make([]Message, 0)
	done := false

	var outMsg Message
	var outChan chan<- Message

	// Event loop with items in queue
	for !done {

		if len(queue) > 0 {
			outMsg = queue[0]
			outChan = shim.output
		} else {
			outChan = nil
		}

		select {

		case msg, ok := <-shim.input:
			// Activity on input queue
			if ok {
				queue = append(queue, msg)
			} else {
				done = true
			}

		case outChan <- outMsg:
			// Item sent on receive queue
			queue = queue[1:]
		}
	}

	// Drain queue
	for _, msg := range queue {
		shim.output <- msg
	}

	// Tear down some things
	close(shim.output)
	queue = []Message{}
}
