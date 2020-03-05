package queue

import (
	"context"
)

type Queue struct {
	queue
	pushCh chan interface{}
	popCh  chan chan interface{}
}

func NewQueue() *Queue {
	return &Queue{
		pushCh: make(chan interface{}),
		popCh:  make(chan chan interface{}),
	}
}

func (q *Queue) loop(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-q.pushCh:
				q.queue.Push(val)
			case outCh := <-q.popCh:
				outCh <- q.queue.Pop()
			}
		}
	}()
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
