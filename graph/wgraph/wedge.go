package wgraph

// EdgeWeighted represents what a weighted edge should be able to do.
type EdgeWeighted[T Weight, N NodeWeighted[T]] interface {
	ToNode() (to N)
	GetWeight() T
}

type EdgeWeightedImpl[T Weight, N NodeWeighted[T]] struct {
	toNode N
	weight T
}

// If E is an edge from nodes A to B, then E.GetToNode() returns B.
func (e *EdgeWeightedImpl[T, N]) ToNode() N {
	return e.toNode
}

func (e *EdgeWeightedImpl[T, N]) GetWeight() T {
	return e.weight
}
