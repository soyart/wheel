package wgraph

import (
	"fmt"

	"github.com/soyart/wheel"
	"github.com/soyart/wheel/tree"
)

// GraphDijkstraImpl wraps [GraphWeighted], where T is generic type numeric types and S is ~string.
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

// dijkstraQueueItem is a priority queue entry carrying a node together with
// the distance it was queued at, frozen at push time. The heap orders on
// this frozen dist, never on node.GetValue() directly, so a node's queue
// position can never go stale relative to its own (mutable) field: a pop
// whose dist no longer matches node.GetValue() has been superseded by a
// cheaper relaxation pushed later, and is simply discarded.
type dijkstraQueueItem[T WeightDijkstra] struct {
	node NodeDijkstra[T]
	dist T
}

// DijkstraShortestPathFrom takes a [NodeDijkstra] startNode, and finds the shortest path from startNode to all other nodes.
// This implementation uses PriorityQueue[T], so the nodes' values must satisfy constraints.Ordered.
func (g *GraphDijkstraImpl[T]) DijkstraShortestPathFrom(startNode NodeDijkstra[T]) *DijstraShortestPath[T] {
	startNode.SetValueOrCost(0)
	startNode.SetPrevious(nil)

	visited := make(map[NodeDijkstra[T]]bool)
	parents := make(map[NodeDijkstra[T]]NodeDijkstra[T])

	pq := tree.NewHeapCustom(wheel.LessFuncBy(
		wheel.Ascending,
		func(item dijkstraQueueItem[T]) T { return item.dist },
	))
	pq.Push(dijkstraQueueItem[T]{node: startNode, dist: 0})

	for !pq.IsEmpty() {
		item, ok := pq.Pop()
		if !ok {
			panic("popped from empty heap - should not happen")
		}

		// Stale entry: a cheaper relaxation already superseded it, or this
		// node was already finalized via another entry. Either way, its
		// distance was already used to relax its neighbors; discard it.
		if item.dist != item.node.GetValue() || visited[item.node] {
			continue
		}
		visited[item.node] = true

		for _, edge := range g.GetNodeEdges(item.node) {
			edgeNode := edge.ToNode()
			if visited[edgeNode] {
				continue
			}

			// If getting to edge from current is cheaper that the edge current cost state,
			// update it to pass via current instead
			if newCost := item.dist + edge.GetWeight(); newCost < edgeNode.GetValue() {
				edgeNode.SetValueOrCost(newCost)
				edgeNode.SetPrevious(item.node)

				// Save (best) path answer to parents
				parents[edgeNode] = item.node

				pq.Push(dijkstraQueueItem[T]{node: edgeNode, dist: newCost})
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
