package timedqueue

import (
	"time"
)

type Item struct {
	Value    interface{}
	Priority time.Time
	index    int
}

func NewItem(value interface{}, when time.Time) *Item {
	return &Item{
		Value:    value,
		Priority: when,
	}
}

func (it *Item) Less(x interface{}) bool {
	other := x.(*Item)
	return it.Priority.Before(other.Priority)
}

func (it *Item) SetIndex(index int) {
	it.index = index
}
