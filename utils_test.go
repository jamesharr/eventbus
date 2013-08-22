package eventbus_test

import (
	"testing"
	"reflect"
)

func AssertSame(t *testing.T, a, b interface{}, message string) {
	if a != b {
		t.Errorf("%#v == %#v assert failed. %s", a, b, message)
	}
}

func AssertEq(t *testing.T, a, b interface{}, message string) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%#v == %#v assert failed. %s", a, b, message)
	}
}
