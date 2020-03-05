package timedqueue

import (
	"container/heap"
	"time"
)

type TimedQueue struct {
	queue priorityQueue
}

func NewTimedQueue() *TimedQueue {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	return &TimedQueue{
		queue: pq,
	}
}

func (pq *TimedQueue) Push(x interface{}) {
	heap.Push(&pq.queue, x)
	return
}

func (pq *TimedQueue) Pop() interface{} {
	return heap.Pop(&pq.queue)
}

func (pq *TimedQueue) Peek() (interface{}, bool) {
	if len(pq.queue) > 0 {
		return pq.queue[0], true
	}
	return nil, false
}

type Item struct {
	Value    interface{}
	Priority time.Time
	index    int
}

type priorityQueue []*Item

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority.Before(pq[j].Priority)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	val := x.(*Item)
	val.index = n
	*pq = append(*pq, val)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	val := old[n-1]
	old[n-1] = nil // avoid memory leak
	val.index = -1 // for safety
	*pq = old[0 : n-1]
	return val
}
