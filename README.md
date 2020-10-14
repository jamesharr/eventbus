# Project Status - Archived

I'm considering this project archived for a variety of reasons 1) I haven't touched it in a long time 2) there are much better event bus options out there 3) I think since the original authoring of this package that go idioms have become more clear and it's likely there's better things out there.

Never the less, this code is still available for reference for anyone who wants to look at it.

# A concurrent EventBus for Go

[![build status](https://secure.travis-ci.org/jamesharr/eventbus.png)](http://travis-ci.org/jamesharr/eventbus)

EventBus is a simple publish and subscribe tool for go. It contains no built-in filtering mechanism,
but it is concurrent in nearly every respect, and won't block any process emitting/publishing a
message.

The main intention of this project is to support a [Go Expect](https://github.com/jamesharr/expect) library.

## Example
```go
package main

import "github.com/jamesharr/eventbus"

func main() {
	// Create the bus
	bus := eventbus.CreateEventBus()

	// Register a couple handlers
	handler1 = make(chan eventbus.Message)
	handler2 = make(chan eventbus.Message)
	bus.Register(handler1) // Register/Unregister are GoRoutine safe.
	bus.Register(handler2)

	// Emit a couple events onto the bus (will never block0
	bus.Emit("Hello")
	bus.Emit("World")

	// Normally done in another GoRoutine, but this demonstrates the buffering.
	msg := <-handler1 // msg = "Hello"
	msg = <-handler1  // msg = "World"

	// Close the bus, and flush out any messages.
	bus.Close()

	// More buffering
	msg = <-handler2   // msg = "Hello"
	msg = <-handler2   // msg = "World"
	_, ok = <-handler2 // OK = false, bus is closed, all events have been flushed.

	// Same here
	_, ok := <-handler1 // OK = false, bus is closed.

	// This will cause a panic, since the bus is closed.
	bus.Emit("DOH!")

}
```

## API Documentation
[API Documentation](http://godoc.org/github.com/jamesharr/eventbus)
