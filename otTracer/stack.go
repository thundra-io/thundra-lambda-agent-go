package otTracer

import (
	"sync"
	"errors"
)

type spanTreeStack struct {
	lock sync.Mutex
	s    []*RawSpanTree
}

func NewStack() *spanTreeStack {
	return &spanTreeStack{sync.Mutex{}, make([]*RawSpanTree, 0),}
}

func (s *spanTreeStack) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *spanTreeStack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}
