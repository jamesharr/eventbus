package eventbus_test

import (
	"github.com/bmizerany/assert"
	"github.com/jamesharr/eventbus"
	"testing"
)

type receiver struct {
	input   chan eventbus.Message
	results chan []eventbus.Message
}

func makeReceiver() receiver {
	input := make(chan eventbus.Message)
	received_messages := make(chan []eventbus.Message)
	go func() {
		messages := make([]eventbus.Message, 0)
		for {
			msg, ok := <-input
			if ok {
				messages = append(messages, msg)
			} else {
				break
			}
		}
		received_messages <- messages
		close(received_messages)
	}()
	return receiver{input, received_messages}
}

func TestBusEmpty(t *testing.T) {
	eb := eventbus.CreateEventBus()
	eb.Close()
}

func TestBusShort(t *testing.T) {
	message_list := []eventbus.Message{1, 5, 7, 3, 5, 8, 9, 90, -1, "HAI"}

	eb := eventbus.CreateEventBus()
	h1 := makeReceiver()
	eb.Register(h1.input)
	for _, msg := range message_list {
		eb.Emit(msg)
	}
	eb.Close()

	assert.Equal(t, message_list, <-h1.results, "Handler 1 received items")
}

func TestMultipleListeners(t *testing.T) {
	message_list := []eventbus.Message{1, 5, 7, 3, 5, 8, 9, 90, -1, "HAI"}

	eb := eventbus.CreateEventBus()
	h1 := makeReceiver()
	eb.Register(h1.input)
	h2 := makeReceiver()
	eb.Register(h2.input)
	for _, msg := range message_list {
		eb.Emit(msg)
	}
	eb.Close()

	assert.Equal(t, message_list, <-h1.results, "Handler 1 received items")
	assert.Equal(t, message_list, <-h2.results, "Handler 2 received items")
}
