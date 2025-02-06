package graph

import (
	"sync"
)

// SafeGraph[N, E, R, W] wraps any BasicGraph[N, E, R, W] graphs.
type SafeGraph[N, E, W any] struct {
	sync.RWMutex
	Graph Graph[N, E, W]
}

// WrapSafeGenericGraph[N, E, R, W] wraps BasicGraph[N, E, R, W]
// with *SafeGraph[N, E, R, W] to use mutex to avoid data races
func WrapSafeGenericGraph[N, E, W any](g Graph[N, E, W]) *SafeGraph[N, E, W] {
	return &SafeGraph[N, E, W]{
		Graph: g,
	}
}

func (g *SafeGraph[N, E, W]) SetDirection(value bool) {
	g.Lock()
	defer g.Unlock()

	g.Graph.SetDirection(value)
}

func (g *SafeGraph[N, E, W]) IsDirected() bool {
	g.RLock()
	defer g.RUnlock()

	return g.Graph.IsDirected()
}

func (g *SafeGraph[N, E, W]) AddNode(node N) {
	g.Lock()
	defer g.Unlock()

	g.Graph.AddNode(node)
}

func (g *SafeGraph[N, E, W]) AddEdgeWeightOrDistance(n1, n2 N, weight W) error {
	g.Lock()
	defer g.Unlock()

	//nolint:wrapcheck
	return g.Graph.AddEdgeWeightOrDistance(n1, n2, weight)
}

func (g *SafeGraph[N, E, W]) AddEdge(n1, n2 N, edge E) error {
	g.Lock()
	defer g.Unlock()

	//nolint:wrapcheck
	return g.Graph.AddEdge(n1, n2, edge)
}

func (g *SafeGraph[N, E, W]) GetNodes() []N {
	g.RLock()
	defer g.RUnlock()

	return g.Graph.GetNodes()
}

func (g *SafeGraph[N, E, W]) GetEdges() []E {
	g.RLock()
	defer g.RUnlock()

	return g.Graph.GetEdges()
}

func (g *SafeGraph[N, E, W]) GetNodeNeighbors(node N) []N {
	g.RLock()
	defer g.RUnlock()

	return g.Graph.GetNodeNeighbors(node)
}

func (g *SafeGraph[N, E, W]) GetNodeEdges(node N) []E {
	g.RLock()
	defer g.RUnlock()

	return g.Graph.GetNodeEdges(node)
}
