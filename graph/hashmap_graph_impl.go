package graph

// HashMapGraphImpl is a basic implementation of Graph, not safe for concurrent code.
type HashMapGraphImpl[T any] struct {
	Directed bool
	Nodes    []Node[T]
	Edges    map[Node[T]][]Node[T]
}

// NewHashMapGraphUnsafe[T] returns the default implementation of unweighted graph (*HashMapGraphImpl[T])
// without the mutex field. If your code is not concurrent, use this type, otherwise,
// consider calling NewGraph[T] instead.
func NewHashMapGraphUnsafe[T any](directed bool) HashMapGraphV1[T] {
	return &HashMapGraphImpl[T]{
		Directed: directed,
		Edges:    make(map[Node[T]][]Node[T]),
	}
}

// NewGraph[T] returns the default implementation of unweighted graph (*HashMapGraphImpl[T])
// wrapped inside a of SafeGraph[N any, E any, R any, W any]
func NewHashMapGraph[T any](directed bool) HashMapGraphV1[T] {
	return WrapSafeGenericGraph[Node[T], Node[T], any](
		NewHashMapGraphUnsafe[T](directed),
	)
}

func (g *HashMapGraphImpl[T]) SetDirection(directed bool) { g.Directed = directed }

func (g *HashMapGraphImpl[T]) IsDirected() bool { return g.Directed }

func (g *HashMapGraphImpl[T]) AddNode(node Node[T]) {
	g.Nodes = append(g.Nodes, node)
}

func (g *HashMapGraphImpl[T]) AddEdgeWeightOrDistance(n1, n2 Node[T], weight any) error {
	if weight != nil {
		return ErrEdgeWeightNotNull
	}

	// Add and edge from n1 leading to n2
	g.Edges[n1] = append(g.Edges[n1], n2)

	if g.Directed {
		return nil
	}
	// If it's not directed, then both nodes have links from and to each other
	g.Edges[n2] = append(g.Edges[n2], n1)
	return nil
}

// HashMapGraphImpl doesn't actually have an edge
func (g *HashMapGraphImpl[T]) AddEdge(n1, n2 Node[T], edge Node[T]) error {
	return g.AddEdgeWeightOrDistance(n1, n2, edge.GetValue())
}

func (g *HashMapGraphImpl[T]) GetNodes() []Node[T] { return g.Nodes }

func (g *HashMapGraphImpl[T]) GetEdges() []Node[T] {
	var edges []Node[T]

	for _, nodeEdges := range g.Edges {
		edges = append(edges, nodeEdges...)
	}

	return edges
}

func (g *HashMapGraphImpl[T]) GetNodeNeighbors(node Node[T]) []Node[T] { return g.Edges[node] }

func (g *HashMapGraphImpl[T]) GetNodeEdges(node Node[T]) []Node[T] { return g.Edges[node] }
