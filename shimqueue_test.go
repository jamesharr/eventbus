/**
 * Created with IntelliJ IDEA.
 * User: james
 * Date: 8/21/13
 * Time: 5:25 PM
 * To change this template use File | Settings | File Templates.
 */
package eventbus_test

import (
	"github.com/jamesharr/eventbus"
	"reflect"
	"testing"
	"math/rand"
)

func assertSame(t *testing.T, a, b interface{}, message string) {
	if a != b {
		t.Errorf("%#v == %#v assert failed. %s", a, b, message)
	}
}

func assertEq(t *testing.T, a, b interface{}, message string) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%#v == %#v assert failed. %s", a, b, message)
	}
}

func FillAndDrain(t *testing.T, items []interface{}, input chan eventbus.Message, output chan eventbus.Message) {

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

	assertEq(t, items, items_received, "Item Equivalency & Ordering test")
}

func createShim() (chan eventbus.Message, chan eventbus.Message) {
	input := make(chan eventbus.Message)
	output := make(chan eventbus.Message)
	eventbus.CreateShimQueue(input, output)
	return input, output
}

func TestEmpty(t *testing.T) {
	input, output := createShim()
	items := []interface{}{}
	FillAndDrain(t, items, input, output)
}

func TestSmall(t *testing.T) {
	input, output := createShim()
	items := []interface{}{5,6,7,8,9}
	FillAndDrain(t, items, input, output)
}

func TestLarge(t *testing.T) {
	input, output := createShim()
	items := make([]interface{}, 1000)
	for i:= 0; i<len(items); i++ {
		items[i] = rand.Int63()
	}
	FillAndDrain(t, items, input, output)
}
