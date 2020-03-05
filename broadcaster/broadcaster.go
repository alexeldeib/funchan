// Package broadcaster is derived from https://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/
// Names and types modified for readability.
//
// Note that this package allows for non-blocking writes and blocking reads, which is rarely desirable in real world situations
package broadcaster

type message struct {
	stream chan message
	data   interface{}
}

// Listener understands how to take a channel of streams, wait for a message, put it back, and
type Listener chan chan message

type Broadcaster struct {
	listenc chan Listener
	sendc   chan<- interface{}
}

type Receiver struct {
	stream chan message
}

// create a new broadcaster object.
func NewBroadcaster() Broadcaster {
	listenc := make(chan Listener)
	sendc := make(chan interface{})
	go func() {
		currc := make(chan message, 1)
		for {
			select {
			case sent := <-sendc:
				if sent == nil {
					currc <- message{}
					return
				}
				c := make(chan message, 1)
				bcast := message{stream: c, data: sent}
				currc <- bcast
				currc = c
			case recv := <-listenc:
				recv <- currc
			}
		}
	}()
	return Broadcaster{
		listenc: listenc,
		sendc:   sendc,
	}
}

// start listening to the broadcasts.
func (b Broadcaster) Listen() Receiver {
	c := make(Listener, 0)
	b.listenc <- c
	return Receiver{stream: <-c}
}

// broadcast a value to all listeners.
func (b Broadcaster) Write(data interface{}) {
	b.sendc <- data
}

// Read a value that has been broadcast, waiting until one is available if necessary.
// This works because each reader has access to the same initial channel.
// When the first message is read, each receiver will read it and discard the previously used channel.
// Since each reader only takes an item once, we terminate when all readers have seen a value.
func (r *Receiver) Read() interface{} {
	bcast := <-r.stream
	data := bcast.data
	r.stream <- bcast
	r.stream = bcast.stream
	return data
}
