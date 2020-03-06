package chanqueue_test

import (
	"context"
	"testing"

	"github.com/alexeldeib/funchan/chanqueue"
)

func Test_Queue(t *testing.T) {
	q := chanqueue.NewChanQueue(context.Background())
	want := []interface{}{"foo", "bar", "baz"}

	for i := range want {
		q.Push(want[i])
	}

	for i := range want {
		got := q.Pop()
		if want[i] != got {
			t.Errorf("expected: %#+v, got %#+v\n", want, got)
		}
	}

	if q.Len() != 0 {
		t.Errorf("expected empty queue after removing only item")
	}
}
