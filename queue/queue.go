package queue

import "time"

type Item struct {
	value    interface{}
	priority time.Time
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority.Before(pq[j].priority)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	val := x.(*Item)
	val.index = n
	*pq = append(*pq, val)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	val := old[n-1]
	old[n-1] = nil // avoid memory leak
	val.index = -1 // for safety
	*pq = old[0 : n-1]
	return val
}

func (pq *PriorityQueue) Peek() interface{} {
	return (*pq)[0]
}
