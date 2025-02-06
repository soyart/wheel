package graph

// Graph represents what an wheel graph should look like.
// It is not intended to be used by production code, but more like an internal building block for wheel graphs.
// It's not minimal, and was designed around flexibility and coverage.
// This interface can be used for both unweighted and weighted graph (see wgraph package).
type Graph[
	N any, // Type for graph node
	E any, // Type for graph edge
	W any, // Type for graph edge weight or node values
] interface {
	// SetDirection sets the directionality of the graph.
	SetDirection(value bool)

	// IsDirected
	IsDirected() bool

	// Add a node to a graph
	AddNode(node N)

	// AddEdge adds a real edge from n1 to n2.
	AddEdge(n1, n2 N, edge E) error

	// AddEdgeWeightOrDistance adds default weighted edge to the graph. If the graph is directional, then AddEdgeWeightOrDistance will only adds edge from n1 to n2. If the graph is undirectional, then both connections (n1 -> n2 and n2 -> n1) should be added.
	AddEdgeWeightOrDistance(n1, n2 N, weight W) error

	// Get all nodes in the graph
	GetNodes() []N

	// Get all edges in the graph
	GetEdges() []E

	// GetNodeNeighbors returns a slice of N connected to node
	GetNodeNeighbors(node N) []N

	// GetNodeEdge takes in a node, and returns the connections from that node.
	GetNodeEdges(node N) []E
}
