package otTracer

type RawSpanTree struct {
	Value    *spanImpl
	Children []*RawSpanTree
}

func (t *RawSpanTree) addChild(child *RawSpanTree) {
	t.Children = append(t.Children, child)
}

func newRawSpanTree(span *spanImpl) *RawSpanTree {
	tree := &RawSpanTree{}
	tree.Value = span
	tree.Children = make([]*RawSpanTree, 0)
	return tree
}