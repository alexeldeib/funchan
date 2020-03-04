package queue

import (
	"container/heap"
	"context"
)

type workqueue struct {
	todo PriorityQueue
	work chan interface{}
}

func newWorkqueue() *workqueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &workqueue{
		todo: pq,
		work: make(chan interface{}),
	}
}

func (w *workqueue) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		}
	}
}
