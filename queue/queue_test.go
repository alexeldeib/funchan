package queue

import (
	"container/heap"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPriorityQueue(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		inputs []*Item
		want   []string
	}{
		"simple": {
			inputs: []*Item{
				{
					value:    "apple",
					priority: now,
				},
				{
					value:    "banana",
					priority: now.Add(time.Second * 30),
				},
				{
					value:    "caramel",
					priority: now.Add(time.Second * -30),
				},
			},
			want: []string{"caramel", "apple", "banana"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := make(PriorityQueue, 0)
			heap.Init(&q)

			for _, input := range tc.inputs {
				heap.Push(&q, input)
			}

			got := []string{}
			for range tc.inputs {
				val := heap.Pop(&q).(*Item)
				got = append(got, val.value.(string))
			}
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Errorf(diff)
			}
		})
	}

}
