package wgraph

// TODO: more relaxed generic types

import (
	"github.com/soyart/wheel/graph"
)

// HashMapGraphWeighted is the default implementation of [GraphWeighted].
type HashMapGraphWeighted[
	N NodeWeighted[W],
	E EdgeWeighted[W, N],
	W Weight,
] struct {
	Directed bool
	Nodes    []N
	Edges    map[NodeWeighted[W]]map[NodeWeighted[W]]EdgeWeighted[W, N]
}

// NewGraphWeightedUnsafe returns the default implementation of [GraphWeighted] without the concurrency wrapper.
func NewGraphWeightedUnsafe[N NodeWeighted[W], E EdgeWeighted[W, N], W Weight](directed bool) GraphWeighted[N, E, W] {
	return &HashMapGraphWeighted[N, E, W]{
		Directed: directed,
		Nodes:    make([]N, 0),
		Edges:    make(map[NodeWeighted[W]]map[NodeWeighted[W]]EdgeWeighted[W, N]),
	}
}

// WrapSafeGraphWeighted wraps any graph g that implements [GraphWeighted] with [graph.SafeGraph].
func WrapSafeGraphWeighted[N NodeWeighted[W], E EdgeWeighted[W, N], W Weight](g GraphWeighted[N, E, W]) GraphWeighted[N, E, W] {
	return graph.WrapSafeGenericGraph(g)
}

// NewGraphWeighted returns the default implementation of [GraphWeighted] with the concurrency safety wrapper.
func NewGraphWeighted[N NodeWeighted[W], E EdgeWeighted[W, N], W Weight](directed bool) GraphWeighted[N, E, W] {
	return WrapSafeGraphWeighted(NewGraphWeightedUnsafe[N, E, W](directed))
}

func (g *HashMapGraphWeighted[N, E, W]) SetDirection(directed bool) {
	g.Directed = directed
}

func (g *HashMapGraphWeighted[N, E, W]) IsDirected() bool {
	return g.Directed
}

func (g *HashMapGraphWeighted[N, E, W]) AddNode(node N) {
	g.Nodes = append(g.Nodes, node)
}

func (g *HashMapGraphWeighted[N, E, W]) AddEdgeWeightOrDistance(n1, n2 N, weight W) error {
	// Overwrite existing edge from n1 to n2, if there is any
	if m := g.Edges[n1]; m == nil {
		g.Edges[n1] = make(map[NodeWeighted[W]]EdgeWeighted[W, N])
	} else if m[n2] != nil {
		return wrapErrConnExists(n1, n2)
	}

	g.Edges[n1][n2] = &EdgeWeightedImpl[W, N]{
		toNode: n2,
		weight: weight,
	}

	if g.Directed {
		return nil
	}

	if m := g.Edges[n2]; m == nil {
		g.Edges[n2] = make(map[NodeWeighted[W]]EdgeWeighted[W, N])
	} else if m[n1] != nil {
		return wrapErrConnExists(n2, n1)
	}

	g.Edges[n2][n1] = &EdgeWeightedImpl[W, N]{
		toNode: n1,
		weight: weight,
	}

	return nil
}

// AddEdge adds edge from n1 to n2
func (g *HashMapGraphWeighted[N, E, W]) AddEdge(n1, n2 N, edge E) error {
	// Overwrite existing edge from n1 to n2, if there is any
	if m := g.Edges[n1]; m == nil {
		g.Edges[n1] = make(map[NodeWeighted[W]]EdgeWeighted[W, N])
	} else if m[n2] != nil {
		return wrapErrConnExists(n1, n2)
	}

	g.Edges[n1][n2] = edge
	return nil
}

func (g *HashMapGraphWeighted[N, E, W]) GetNodes() []N {
	return g.Nodes
}

func (g *HashMapGraphWeighted[N, E, W]) GetEdges() []E {
	var edges []E
	for _, nodeEdges := range g.Edges {
		for _, edge := range nodeEdges {
			edges = append(edges, edge.(E))
		}
	}

	return edges
}

func (g *HashMapGraphWeighted[N, E, W]) GetNodeNeighbors(node N) []N {
	edges := g.Edges[node]
	neighbors := make([]N, len(edges))
	var c int
	for _, edge := range edges {
		neighbors[c] = edge.ToNode()
		c++
	}
	return neighbors
}

func (g *HashMapGraphWeighted[N, E, W]) GetNodeEdges(node N) []E {
	edgesMap := g.Edges[node]
	edges := make([]E, len(edgesMap))
	c := 0
	for _, edge := range edgesMap {
		edges[c] = edge.(E)
		c++
	}
	return edges
}
