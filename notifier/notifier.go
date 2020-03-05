package notifier

import "context"

type state struct {
	stream chan state
}

// Listener understands how to take a channel of streams, wait for a state, put it back, and
type Listener chan chan state

// Notifier permits signaling multiple waiters an event has occured,.
// It does not allow sending additional data with the message.
type Notifier struct {
	listenc chan Listener
	sendc   chan<- interface{}
}

type Waiter struct {
	stream chan state
}

// create a new broadcaster object.
func NewNotifier(ctx context.Context) Notifier {
	listenc := make(chan Listener)
	sendc := make(chan interface{})
	go func() {
		currc := make(chan state, 1)
		for {
			select {
			case <-ctx.Done():
				return
			case <-sendc:
				c := make(chan state, 1)
				bcast := state{stream: c}
				currc <- bcast
				currc = c
			case recv := <-listenc:
				recv <- currc
			}
		}
	}()
	return Notifier{
		listenc: listenc,
		sendc:   sendc,
	}
}

func (b Notifier) NewWaiter() Waiter {
	c := make(Listener, 0)
	b.listenc <- c
	return Waiter{stream: <-c}
}

func (b Notifier) Notify() {
	b.sendc <- struct{}{}
}

// Wait for a value to be broadcast, discarding it.
// This works because each reader has access to the same initial channel.
// When the first message is read, each receiver will read it and discard the previously used channel.
// Since each reader only takes an item once, we terminate when all readers have seen a value.
func (w *Waiter) Wait() {
	bcast := <-w.stream
	w.stream <- bcast
	w.stream = bcast.stream
	return
}
