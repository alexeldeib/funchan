/*
Package workqueue provides a task queue where tasks are executed not before a specific time in the future.

Tasks are sorted into a priority queue by time of execution. A dispatcher routine peeks at the root of the
priority queue and feeds items ready for execution into a work channel. Multiple consumers may read work off
the queue. Any given value in the queue may  only be under processing by one worker at a time.
*/
package workqueue

import (
	"context"
	"sync"
	"time"

	"github.com/alexeldeib/funchan/queue"
	"github.com/alexeldeib/funchan/timedqueue"
)

type WorkQueue struct {
	waiting *timedqueue.TimedQueue
	ready   queue.Queue
	out     chan interface{}
}

func NewWorkQueue() *WorkQueue {
	return &WorkQueue{
		waiting: timedqueue.NewTimedQueue(),
		ready:   queue.NewQueue(),
		out:     make(chan interface{}, 1),
	}
}

func (w *WorkQueue) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		w.feed(ctx)
		wg.Done()
	}()
	go func() {
		w.drain(ctx)
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}
	}

}

func (w *WorkQueue) Push(item interface{}, when time.Time) {
	w.waiting.Push(&timedqueue.Item{
		Value:    item,
		Priority: when,
	})
}

func (w *WorkQueue) Pop() interface{} {
	return <-w.out
}

func (w *WorkQueue) feed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			now := time.Now()

			// TODO(ace): should probably do something smarter here
			head, ok := w.waiting.Peek()
			if !ok {
				continue
			}

			item := head.(*timedqueue.Item)
			if now.After(item.Priority) {
				head := w.waiting.Pop()
				val := head.(*timedqueue.Item).Value
				w.ready.Push(val)
			}
		}
	}
}

func (w *WorkQueue) drain(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if w.ready.Peek() != nil {
				w.out <- w.ready.Pop()
			}
		}
	}
}
