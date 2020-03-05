/*
Package workqueue provides a task queue where tasks are executed not before a specific time in the future.

Tasks are sorted into a priority queue by time of execution. A dispatcher routine peeks at the root of the
priority queue and feeds items ready for execution into a work channel. Multiple consumers may read work off
the queue. Any given value in the queue may  only be under processing by one worker at a time.
*/
package workqueue

import (
	"context"
	"time"

	"github.com/alexeldeib/funchan/priorityqueue"
	"github.com/alexeldeib/funchan/queue"
	"github.com/alexeldeib/funchan/timedqueue"
)

type WorkQueue struct {
	waiting *priorityqueue.PriorityQueue
	ready   *queue.Queue
	out     chan interface{}
}

func NewWorkQueue(ctx context.Context) *WorkQueue {
	wq := &WorkQueue{
		waiting: priorityqueue.NewPriorityQueue(),
		ready:   queue.NewQueue(context.Background()),
		out:     make(chan interface{}, 1),
	}

	routines := []func(ctx context.Context) error{
		wq.feed,
		wq.drain,
	}

	for _, fn := range routines {
		do(ctx, fn)
	}

	return wq
}

func do(ctx context.Context, f func(ctx context.Context) error) {
	go func() {
		f(ctx)
	}()
}

func (wq *WorkQueue) Push(item interface{}, when time.Time) {
	wq.waiting.Push(timedqueue.NewItem(item, when))
}

func (wq *WorkQueue) Pop() interface{} {
	return <-wq.out
}

func (wq *WorkQueue) feed(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			now := time.Now()

			// TODO(ace): should probably do something smarter here
			head := wq.waiting.Peek()
			if head == nil {
				continue
			}

			item := head.(*timedqueue.Item)
			if now.After(item.Priority) {
				head := wq.waiting.Pop()
				val := head.(*timedqueue.Item).Value
				wq.ready.Push(val)
			}
		}
	}
}

func (wq *WorkQueue) drain(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if wq.ready.Peek() != nil {
				wq.out <- wq.ready.Pop()
			}
		}
	}
}
