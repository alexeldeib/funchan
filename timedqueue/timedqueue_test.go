package timedqueue_test

import (
	"testing"
	"time"

	"github.com/alexeldeib/funchan/priorityqueue"
	"github.com/alexeldeib/funchan/timedqueue"
	"github.com/google/go-cmp/cmp"
)

func TestPriorityQueue(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		inputs []priorityqueue.Heapable
		want   []string
	}{
		"simple": {
			inputs: []priorityqueue.Heapable{
				timedqueue.NewItem("apple", now),
				timedqueue.NewItem("banana", now.Add(time.Second*30)),
				timedqueue.NewItem("caramel", now.Add(time.Second*-30)),
			},
			want: []string{"caramel", "apple", "banana"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := priorityqueue.NewPriorityQueue()

			for _, input := range tc.inputs {
				q.Push(input)
			}

			got := []string{}
			for range tc.inputs {
				val := q.Pop().(*timedqueue.Item)
				got = append(got, val.Value.(string))
			}
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Errorf(diff)
			}
		})
	}

}
