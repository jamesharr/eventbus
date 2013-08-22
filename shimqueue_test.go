package eventbus_test

import (
	"testing"
	"math/rand"
	"github.com/jamesharr/eventbus"
)

func ShimFillAndDrain(t *testing.T, items []interface{}, input chan eventbus.Message, output chan eventbus.Message) {

	// Send all
	for _,item := range items {
		t.Logf("sending %#v", item)
		input <- item
	}
	close(input)

	// Receive all
	items_received := make([]interface{},len(items))[0:0]
	for {
		item, ok := <-output
		if ok {
			t.Logf("received %#v", item)
			items_received = append(items_received, item)
		} else {
			break
		}
	}

	AssertEq(t, items, items_received, "Item Equivalency & Ordering test")
}

func createShim() (chan eventbus.Message, chan eventbus.Message) {
	input := make(chan eventbus.Message)
	output := make(chan eventbus.Message)
	eventbus.CreateShimQueue(input, output)
	return input, output
}

func TestShimEmpty(t *testing.T) {
	input, output := createShim()
	items := []interface{}{}
	ShimFillAndDrain(t, items, input, output)
}

func TestShimSmall(t *testing.T) {
	input, output := createShim()
	items := []interface{}{5,6,7,8,9}
	ShimFillAndDrain(t, items, input, output)
}

func TestShimLarge(t *testing.T) {
	input, output := createShim()
	items := make([]interface{}, 1000)
	for i:= 0; i<len(items); i++ {
		items[i] = rand.Int63()
	}
	ShimFillAndDrain(t, items, input, output)
}
