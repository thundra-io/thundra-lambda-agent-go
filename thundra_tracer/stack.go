package ttracer

import (
	"errors"
	"sync"
)

// spanTreeStack holds rawspantrees in a stack.
// stack is used to hold activity records of functions, but uses references of spantrees instead of spans
// in order to further manipulate them if a new child span is added
type spanStack struct {
	lock sync.Mutex
	slice    []*RawSpan
}

// NewSpanStack returns a new stack to store spans
func NewSpanStack() *spanTreeStack {
	return &spanStack{sync.Mutex{}, make([]*RawSpan, 0)}
}

// Push pushes on top of stack.
func (s *spanStack) Push(v *RawSpan) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.slice = append(s.slice, v)
}

// Pop pops top of stack and returns it.
func (s *spanStack) Pop() (*RawSpan, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.slice)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.slice[l-1]
	s.slice = s.slice[:l-1]
	return res, nil
}

// Top returns top of stack.
func (s *spanStack) Top() (*RawSpan, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.slice)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.slice[l-1]
	return res, nil
}

// Clear clears the stack
func (s *spanStack) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.slice = nil
}
