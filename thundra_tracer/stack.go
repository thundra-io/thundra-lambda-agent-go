package thundra_tracer

import (
	"sync"
	"errors"
)

// spanTreeStack holds rawspantrees in a stack.
// stack is used to hold activity records of functions, but uses references of spantrees instead of spans
// in order to further manipulate them if a new child span is added
type spanTreeStack struct {
	lock sync.Mutex
	t    []*RawSpanTree
}

// NewStack returns a new spanTreeStack
func NewStack() *spanTreeStack {
	return &spanTreeStack{sync.Mutex{}, make([]*RawSpanTree, 0)}
}

// Push Pushes on top of stack.
func (s *spanTreeStack) Push(v *RawSpanTree) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.t = append(s.t, v)
}

// Pop pops top of stack and returns it.
func (s *spanTreeStack) Pop() (*RawSpanTree, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.t)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.t[l-1]
	s.t = s.t[:l-1]
	return res, nil
}

// Top returns top of stack.
func (s *spanTreeStack) Top() (*RawSpanTree, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.t)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.t[l-1]
	return res, nil
}
