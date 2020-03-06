package queue

import (
	"context"
)

type Interface interface {
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
	Len() int
}

type Queue struct {
	items []interface{}
}

func NewQueue(_ context.Context) *Queue {
	return &Queue{}
}

func (q *Queue) Push(val interface{}) {
	q.items = append(q.items, val)
}

func (q *Queue) Pop() interface{} {
	var head interface{}
	if len(q.items) == 1 {
		head, q.items = q.items[0], nil
	}
	if len(q.items) > 1 {
		head, q.items = q.items[0], q.items[1:]
	}
	return head
}

func (q *Queue) Peek() interface{} {
	if q.Len() < 1 {
		return nil
	}
	return q.items[0]
}

func (q *Queue) Len() int {
	return len(q.items)
}
