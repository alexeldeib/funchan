package priorityqueue

import (
	"container/heap"
	"context"

	"github.com/alexeldeib/funchan/queue"
)

type Heapable interface {
	Less(priority interface{}) bool
	SetIndex(i int)
}

var _ queue.Interface = &PriorityQueue{}

type PriorityQueue struct {
	queue  priorityQueue
	pushCh chan interface{}
	popCh  chan chan interface{}
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		pushCh: make(chan interface{}),
		popCh:  make(chan chan interface{}),
	}
	go pq.loop(context.Background())
	return pq
}

func (pq *PriorityQueue) loop(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-pq.pushCh:
				heap.Push(&pq.queue, val)
			case outCh := <-pq.popCh:
				outCh <- heap.Pop(&pq.queue)
			}
		}
	}()
}

func (pq *PriorityQueue) Push(x interface{}) {
	pq.pushCh <- x
	return
}

// Pop takes a value off the heap.
func (pq *PriorityQueue) Pop() interface{} {
	out := make(chan interface{})
	pq.popCh <- out
	return <-out
}

func (pq *PriorityQueue) Peek() interface{} {
	if pq.queue.Len() < 1 {
		return nil
	}
	return pq.queue.items[0]
}

func (pq *PriorityQueue) Len() int {
	return pq.queue.Len()
}

// priority queue is the underlying, non-thread safe slice implementing the heap interface.
type priorityQueue struct {
	items []Heapable
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
