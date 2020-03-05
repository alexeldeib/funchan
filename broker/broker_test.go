package broker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexeldeib/funchan/broker"
)

func Test_SinglePubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := broker.NewBroker(ctx, 0)

	want := "hello, world!"

	ch := b.Subscribe()
	b.Publish(want)

	timeout := 1 * time.Second
	select {
	case <-time.After(timeout):
		t.Errorf("failed to receive message within timeout %02d", timeout/time.Second)
	case got := <-ch:
		if got != want {
			t.Errorf("expected %#+v, but got %#+v\n", want, got)
		}
	}
}

func drain(id int, ch <-chan interface{}) {
	for datum := range ch {
		fmt.Printf("worker %d received message: %#+v\n", id, datum)
	}
}

func generate(min, max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := min; i < max; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
