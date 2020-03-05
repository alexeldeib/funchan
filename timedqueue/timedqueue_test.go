package timedqueue_test

import (
	"testing"
	"time"

	"github.com/alexeldeib/funchan/timedqueue"
	"github.com/google/go-cmp/cmp"
)

func TestPriorityQueue(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		inputs []*timedqueue.Item
		want   []string
	}{
		"simple": {
			inputs: []*timedqueue.Item{
				{
					Value:    "apple",
					Priority: now,
				},
				{
					Value:    "banana",
					Priority: now.Add(time.Second * 30),
				},
				{
					Value:    "caramel",
					Priority: now.Add(time.Second * -30),
				},
			},
			want: []string{"caramel", "apple", "banana"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := timedqueue.NewTimedQueue()

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
