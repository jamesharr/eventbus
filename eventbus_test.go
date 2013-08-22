package eventbus_test

import (
	"testing"
	"github.com/jamesharr/eventbus"
)



func TestEmpty(t *testing.T) {
	eb := eventbus.CreateEventBus()
	eb.Close()
}
