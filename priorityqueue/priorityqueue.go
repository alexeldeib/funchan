package priorityqueue

import (
	"container/heap"
)

type Heapable interface {
	Less(priority interface{}) bool
	SetIndex(i int)
}

type PriorityQueue struct {
	queue priorityQueue
}

type priorityQueue struct {
	items []Heapable
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

func (pq *PriorityQueue) Push(x interface{}) {
	heap.Push(&pq.queue, x)
	return
}

func (pq *PriorityQueue) Pop() interface{} {
	return heap.Pop(&pq.queue)
}

func (pq *PriorityQueue) Peek() (interface{}, bool) {
	if pq.Len() > 0 {
		return pq.queue.items[0], true
	}
	return nil, false
}

func (pq *PriorityQueue) Len() int {
	return pq.queue.Len()
}

func (pq *priorityQueue) Len() int {
	return len(pq.items)
}

func (pq *priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.items[i].Less(pq.items[j])
}

func (pq *priorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].SetIndex(i)
	pq.items[j].SetIndex(j)
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(pq.items)
	val := x.(Heapable)
	val.SetIndex(n)
	pq.items = append(pq.items, val)
}

func (pq *priorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	val := old[n-1]
	old[n-1] = nil   // avoid memory leak
	val.SetIndex(-1) // for safety
	pq.items = old[0 : n-1]
	return val
}
