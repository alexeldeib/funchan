package chanqueue

import (
	"context"

	"github.com/alexeldeib/funchan/queue"
)

var _ queue.Interface = &ChanQueue{}

type ChanQueue struct {
	queue  *queue.Queue
	pushCh chan interface{}
	popCh  chan chan interface{}
}

func NewChanQueue(ctx context.Context) *ChanQueue {
	q := &ChanQueue{
		queue:  queue.NewQueue(ctx),
		pushCh: make(chan interface{}),
		popCh:  make(chan chan interface{}),
	}
	go q.loop(ctx)
	return q
}

func (q *ChanQueue) loop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val := <-q.pushCh:
			q.queue.Push(val)
		case outCh := <-q.popCh:
			outCh <- q.queue.Pop()
		}
	}
}

func (q *ChanQueue) Push(x interface{}) {
	q.pushCh <- x
	return
}

// Pop takes a value off the heap.
func (q *ChanQueue) Pop() interface{} {
	out := make(chan interface{})
	q.popCh <- out
	return <-out
}

func (q *ChanQueue) Peek() interface{} {
	if q.queue.Len() < 1 {
		return nil
	}
	return q.queue.Peek()
}

func (q *ChanQueue) Len() int {
	return q.queue.Len()
}
