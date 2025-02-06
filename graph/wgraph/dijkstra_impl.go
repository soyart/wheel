package wgraph

import (
	"fmt"
	"reflect"

	"github.com/soyart/wheel"
	"github.com/soyart/wheel/tree"
)

// GraphDijkstraImpl[T] wraps GraphWeightedImpl[T], where T is generic type numeric types and S is ~string.
// It uses HashMapGraphWeighted as the underlying wheel structure.
type GraphDijkstraImpl[T WeightDijkstra] struct {
	graph GraphWeighted[NodeDijkstra[T], EdgeWeighted[T, NodeDijkstra[T]], T]
}

func (g *GraphDijkstraImpl[T]) SetDirection(directed bool) {
	g.graph.SetDirection(directed)
}

func (g *GraphDijkstraImpl[T]) IsDirected() bool {
	return g.graph.IsDirected()
}

func (g *GraphDijkstraImpl[T]) AddNode(node NodeDijkstra[T]) {
	g.graph.AddNode(node)
}

func (g *GraphDijkstraImpl[T]) AddEdgeWeightOrDistance(n1, n2 NodeDijkstra[T], weight T) error {
	if weight < 0 {
		return fmt.Errorf("negative edge weight %v: %w", weight, ErrDijkstraNegativeWeightEdge)
	}

	//nolint:wrapcheck
	return g.graph.AddEdgeWeightOrDistance(n1, n2, weight)
}

func (g *GraphDijkstraImpl[T]) AddEdge(n1, n2 NodeDijkstra[T], edge EdgeWeighted[T, NodeDijkstra[T]]) error {
	weight := edge.GetWeight()
	if weight < 0 {
		return fmt.Errorf("negative edge weight %v: %w", weight, ErrDijkstraNegativeWeightEdge)
	}

	//nolint:wrapcheck
	return g.graph.AddEdge(n1, n2, edge)
}

func (g *GraphDijkstraImpl[T]) GetNodes() []NodeDijkstra[T] {
	nodes := any(g.graph.GetNodes())
	return nodes.([]NodeDijkstra[T])
}

func (g *GraphDijkstraImpl[T]) GetEdges() []EdgeWeighted[T, NodeDijkstra[T]] {
	return g.graph.GetEdges()
}

func (g *GraphDijkstraImpl[T]) GetNodeNeighbors(node NodeDijkstra[T]) []NodeDijkstra[T] {
	return g.graph.GetNodeNeighbors(node)
}

func (g *GraphDijkstraImpl[T]) GetNodeEdges(node NodeDijkstra[T]) []EdgeWeighted[T, NodeDijkstra[T]] {
	return g.graph.GetNodeEdges(node)
}

// DjisktraFrom takes a *NodeImpl[T] startNode, and finds the shortest path from startNode to all other nodes.
// This implementation uses PriorityQueue[T], so the nodes' values must satisfy constraints.Ordered.
func (g *GraphDijkstraImpl[T]) DijkstraShortestPathFrom(startNode NodeDijkstra[T]) *DijstraShortestPath[T] {
	startNode.SetValueOrCost(0)
	startNode.SetPrevious(nil)

	visited := make(map[NodeDijkstra[T]]bool)
	parents := make(map[NodeDijkstra[T]]NodeDijkstra[T])

	// pq := list.NewPriorityQueue[T](wheel.Ascending)
	// heap.Push(pq, startNode)
	pq := tree.NewHeap[T](wheel.Ascending)
	pq.PushGetter(startNode)

	for !pq.IsEmpty() {
		// Pop the top of pq and mark it as visited
		popped := pq.PopGetter()
		if popped == nil {
			panic("popped nil - should not happen")
		}

		current, ok := popped.(NodeDijkstra[T])
		if !ok {
			typeOfCurrent := reflect.TypeOf(current)
			panic(fmt.Sprintf("current is %s, not *Node[T]", typeOfCurrent))
		}

		visited[current] = true
		edges := g.GetNodeEdges(current)

		for _, edge := range edges {
			edgeNode := edge.ToNode()

			// Skip visited
			if visited[edgeNode] {
				continue
			}

			pq.PushGetter(edgeNode)

			// If getting to edge from current is cheaper that the edge current cost state,
			// update it to pass via current instead
			if newCost := current.GetValue() + edge.GetWeight(); newCost < edgeNode.GetValue() {
				edgeNode.SetValueOrCost(newCost)
				edgeNode.SetPrevious(current)

				// Save (best) path answer to parents
				parents[edgeNode] = current
			}
		}
	}

	return &DijstraShortestPath[T]{
		From:  startNode,
		Paths: parents,
	}
}

// ReconstructPathTo reconstructs shortest path as slice of nodes, from `d.from` to `to`
// using DijkstraShortestPathReconstruct
func (d *DijstraShortestPath[T]) ReconstructPathTo(to NodeDijkstra[T]) []NodeDijkstra[T] {
	return DijkstraShortestPathReconstruct(d.Paths, d.From, to)
}
