package queue

import (
	"context"
)

type Queue struct {
	queue
	pushCh chan interface{}
	popCh  chan chan interface{}
}

func NewQueue(ctx context.Context) *Queue {
	q := &Queue{
		pushCh: make(chan interface{}),
		popCh:  make(chan chan interface{}),
	}
	go q.loop(ctx)
	return q
}

func (q *Queue) loop(ctx context.Context) error {
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

func (q *Queue) Push(x interface{}) {
	q.pushCh <- x
	return
}

// Pop takes a value off the heap.
func (q *Queue) Pop() interface{} {
	out := make(chan interface{})
	q.popCh <- out
	return <-out
}

func (q *Queue) Peek() interface{} {
	if q.queue.Len() < 1 {
		return nil
	}
	return q.queue.items[0]
}

func (q *Queue) Len() int {
	return q.queue.Len()
}

type queue struct {
	items []interface{}
}

func (q *queue) Push(val interface{}) {
	q.items = append(q.items, val)
}

func (q *queue) Pop() interface{} {
	var head interface{}
	if len(q.items) == 1 {
		head, q.items = q.items[0], nil
	}
	if len(q.items) > 1 {
		head, q.items = q.items[0], q.items[1:]
	}
	return head
}

func (q *queue) Peek() interface{} {
	if q.Len() < 1 {
		return nil
	}
	return q.items[0]
}

func (q *queue) Len() int {
	return len(q.items)
}
