package broker

import (
	"context"
)

type Broker struct {
	msgCh   chan interface{}
	subCh   chan chan interface{}
	unsubCh chan chan interface{}
}

func NewBroker(ctx context.Context, buffer int) *Broker {
	b := &Broker{
		msgCh:   make(chan interface{}, buffer),
		subCh:   make(chan chan interface{}, buffer),
		unsubCh: make(chan chan interface{}, buffer),
	}

	go loop(ctx, b)

	return b
}

func loop(ctx context.Context, b *Broker) {
	subs := map[chan interface{}]struct{}{}
	defer func() {
		for sub := range subs {
			close(sub)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case sub := <-b.subCh:
			subs[sub] = struct{}{}
		case unsub := <-b.unsubCh:
			delete(subs, unsub)
			close(unsub)
		case datum := <-b.msgCh:
			undelivered := map[chan interface{}]struct{}{}
			for k, v := range subs {
				undelivered[k] = v
			}
			for len(undelivered) > 0 {
				for sub := range undelivered {
					select {
					case <-ctx.Done():
						return
					case sub <- datum:
						delete(undelivered, sub)
					default:
					}
				}
			}
		}
	}
}

func (b *Broker) Subscribe() chan interface{} {
	ch := make(chan interface{})
	b.subCh <- ch
	return ch
}

func (b *Broker) Unsubscribe(ch chan interface{}) {
	b.unsubCh <- ch
}

func (b *Broker) Publish(msg interface{}) {
	b.msgCh <- msg
}
