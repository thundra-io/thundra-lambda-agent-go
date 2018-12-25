package ttracer

//
type RawSpanTree struct {
	Value    *RawSpan
	Children []*RawSpanTree
}

func (t *RawSpanTree) addChild(child *RawSpanTree) {
	t.Children = append(t.Children, child)
}

func newRawSpanTree(span *RawSpan) *RawSpanTree {
	tree := &RawSpanTree{}
	tree.Value = span
	tree.Children = make([]*RawSpanTree, 0)
	return tree
}

// Walk traverses a tree depth-first
func (t *RawSpanTree) Walk(ch chan *RawSpan) {
	if t == nil {
		return
	}
	ch <- t.Value
	for _, child := range t.Children {
		child.Walk(ch)
	}
}
