package set

type set map[interface{}]struct{}

func (s set) has(item interface{}) bool {
	_, ok := s[item]
	return ok
}

func (s set) insert(item interface{}) {
	s[item] = struct{}{}
}

func (s set) delete(item interface{}) {
	delete(s, item)
}
