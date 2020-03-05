package workqueue_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexeldeib/funchan/workqueue"
)

func Test_WorkQueue(t *testing.T) {
	wq := workqueue.NewWorkQueue()
	go func(t *testing.T) {
		if err := wq.Start(context.Background()); err != nil {
			t.Errorf("failed to start workqueue continuously: %#+v", err)
		}
	}(t)

	now := time.Now()
	wq.Push("foo", now.Add(50*time.Millisecond))
	wq.Push("bar", now.Add(150*time.Millisecond))
	wq.Push("baz", now.Add(100*time.Millisecond))

	expected := []interface{}{"foo", "baz", "bar"}
	times := []time.Duration{}
	// start := time.Now()
	for i := 0; i < 3; i++ {
		got := wq.Pop()
		delay := time.Now().Sub(now)
		times = append(times, delay)
		if got != expected[i] {
			t.Errorf("expected %#+v, but got %#+v", expected[i], got)
		}
	}
	for i := range times {
		prev := i - 1
		if prev < 0 {
			prev = 0
		}

		wiggle := 10 * time.Millisecond
		delay := times[i]

		low, high := delay-wiggle, delay+wiggle
		if delay < low || delay > high {
			t.Errorf("expected delay between %02d and %02d, got %02d", low, high, delay)
		}
	}
}

func Test_WorkQueueShutdown(t *testing.T) {
	wq := workqueue.NewWorkQueue()

	timeout, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	ctx, cancel := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())

	go func(t *testing.T) {
		if err := wq.Start(ctx); err != nil {
			done()
		}
	}(t)

	cancel()
	select {
	case <-timeout.Done():
		t.Errorf("failed to shut down work queue gracefully")
	case <-doneCtx.Done():
	}
}