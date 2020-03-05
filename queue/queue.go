package queue

type Queue []interface{}

func NewQueue() Queue {
	return nil
}

func (q *Queue) Push(val interface{}) {
	*q = append(*q, val)
}

func (q *Queue) Pop() interface{} {
	var head interface{}
	if len(*q) < 1 {
		return nil
	}
	if len(*q) == 1 {
		head, *q = (*q)[0], nil
	}
	if len(*q) > 1 {
		head, *q = (*q)[0], (*q)[1:]
	}
	return head
}

func (q *Queue) Peek() interface{} {
	if len(*q) < 1 {
		return nil
	}
	return (*q)[0]
}

func (q *Queue) Len() int {
	return len(*q)
}
