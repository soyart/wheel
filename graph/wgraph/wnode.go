package wgraph

import "github.com/soyart/wheel/graph"

// NodeWeighted should be able to put in a priority queue, just in case topological sort is needed.
type NodeWeighted[T Weight] interface {
	// Inherit some from unweighted graphs
	graph.Node[T]

	GetKey() string         // Get the node's key, names, IDs
	SetValueOrCost(value T) // Save cost or value to the node
}

type NodeWeightedImpl[T Weight] struct {
	Name        string
	ValueOrCost T
	// Previous is for Dijkstra shortest path algorithm, or other sort programs.
}

type NodeDijkstraImpl[T WeightDijkstra] struct {
	NodeWeightedImpl[T]
	Previous NodeDijkstra[T]
}

func (n *NodeWeightedImpl[T]) GetValue() T {
	return n.ValueOrCost
}

func (n *NodeWeightedImpl[T]) GetKey() string {
	return n.Name
}

func (n *NodeWeightedImpl[T]) SetValueOrCost(value T) {
	n.ValueOrCost = value
}

func (n *NodeDijkstraImpl[T]) GetPrevious() NodeDijkstra[T] {
	return n.Previous
}

func (n *NodeDijkstraImpl[T]) SetPrevious(prev NodeDijkstra[T]) {
	n.Previous = prev
}
