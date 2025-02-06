package wgraph

import (
	"reflect"
	"testing"

	"github.com/soyart/wheel"
)

type dijkstraTestUtils[T WeightDijkstra] struct {
	inititalValue      T
	expectedFinalValue T
	expectedPathHops   int
	expectedPathway    []*NodeDijkstraImpl[T]

	edges []*struct {
		to     *NodeDijkstraImpl[T]
		weight T
	}
}

// TODO: Add tests for other types

func TestDijkstra(t *testing.T) {
	const (
		nameStart  = "start"
		nameFinish = "finish"
	)

	testDijkstra[uint](t, nameStart, nameFinish)
	testDijkstra[uint8](t, nameStart, nameFinish)
	testDijkstra[int32](t, nameStart, nameFinish)
	testDijkstra[float64](t, nameStart, nameFinish)
}

func constructDijkstraTestGraph[T WeightDijkstra](nameStart, nameFinish string) map[NodeDijkstra[T]]*dijkstraTestUtils[T] {
	// TODO: infinity is way too low, because dijkstraWeight also has uint8
	infinity := T(100)
	nodeStart := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        nameStart,
			ValueOrCost: T(0),
		},
	}
	nodeA := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        "A",
			ValueOrCost: infinity,
		},
	}
	nodeB := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        "B",
			ValueOrCost: infinity,
		},
	}
	nodeC := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        "C",
			ValueOrCost: infinity,
		},
	}
	nodeD := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        "D",
			ValueOrCost: infinity,
		},
	}
	nodeFinish := &NodeDijkstraImpl[T]{
		NodeWeightedImpl: NodeWeightedImpl[T]{
			Name:        nameFinish,
			ValueOrCost: infinity,
		},
	}
	m := map[NodeDijkstra[T]]*dijkstraTestUtils[T]{
		nodeStart: {
			inititalValue:      T(0),
			expectedFinalValue: T(0),
			expectedPathHops:   1,
			expectedPathway:    []*NodeDijkstraImpl[T]{},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{
				{
					to:     nodeA,
					weight: T(2),
				}, {
					to:     nodeB,
					weight: T(4),
				}, {
					to:     nodeD,
					weight: T(4),
				},
			},
		},
		nodeFinish: {
			inititalValue:      infinity,
			expectedFinalValue: T(7),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T]{nodeStart, nodeD, nodeFinish},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{},
		},
		nodeA: {
			inititalValue:      infinity,
			expectedFinalValue: T(2),
			expectedPathHops:   2,
			expectedPathway:    []*NodeDijkstraImpl[T]{nodeStart, nodeA},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{
				{
					to:     nodeB,
					weight: T(1),
				},
				{
					to:     nodeC,
					weight: T(2),
				},
			},
		},
		nodeB: {
			inititalValue:      infinity,
			expectedFinalValue: T(3),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T]{nodeStart, nodeA, nodeB},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{
				{
					to:     nodeFinish,
					weight: T(5),
				},
			},
		},
		nodeC: {
			inititalValue:      infinity,
			expectedFinalValue: T(4),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T]{nodeStart, nodeA, nodeC},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{
				{
					to:     nodeD,
					weight: T(2),
				},
			},
		},
		nodeD: {
			inititalValue:      infinity,
			expectedFinalValue: T(4),
			expectedPathHops:   2,
			expectedPathway:    []*NodeDijkstraImpl[T]{nodeStart, nodeD},
			edges: []*struct {
				to     *NodeDijkstraImpl[T]
				weight T
			}{
				{
					to:     nodeFinish,
					weight: T(3),
				},
			},
		},
	}
	return m
}

// The weighted graph used in this test can be viewed at assets/dijkstra_test_graph.png
func testDijkstra[T WeightDijkstra](t *testing.T, nameStart, nameFinish string) {
	nodesMap := constructDijkstraTestGraph[T](nameStart, nameFinish)

	// Prepare graph
	directed := true
	djikGraph := NewDijkstraGraph[T](directed)
	for node, util := range nodesMap {
		// Add node
		djikGraph.AddNode(node)
		// Add edges
		nodeEdges := util.edges
		for _, edge := range nodeEdges {
			if err := djikGraph.AddEdgeWeightOrDistance(node, edge.to, edge.weight); err != nil {
				t.Error(err.Error())
			}
		}
	}

	var startNode NodeDijkstra[T]
	for node := range nodesMap {
		if node.GetKey() == nameStart {
			startNode = node
		}
	}

	dijkShortestPaths := djikGraph.DijkstraShortestPathFrom(startNode)
	var dummyT T
	for src, dst := range dijkShortestPaths.Paths {
		t.Log(reflect.TypeOf(dummyT).String(), "from", dijkShortestPaths.From.GetKey(), "src", src.GetKey(), src.GetValue(), "dst", dst.GetKey(), dst.GetValue())
	}

	// The check is hard-coded
	for node, util := range nodesMap {
		// Test costs
		if node.GetValue() != util.expectedFinalValue {
			t.Fatalf("invalid answer for (%s->%s): %v, expecting %v", nameStart, node.GetKey(), node.GetValue(), util.expectedFinalValue)
		}

		if node == startNode {
			continue
		}
		nodeKey := node.GetKey()
		// t.Logf("dst node: %v\n", nodeKey)
		// Test paths
		pathToNode := dijkShortestPaths.ReconstructPathTo(node)
		wheel.ReverseInPlace(pathToNode)
		if hops := len(pathToNode); hops != util.expectedPathHops {
			t.Fatalf("invalid returned path length (%s->%s): %d, expecting %d", nameStart, nodeKey, hops, util.expectedPathHops)
		}
		for i, actual := range pathToNode {
			expected := util.expectedPathway[i]
			// t.Logf("prev[%d]: %v, expected: %v\n", i, actual.GetKey(), expected.GetKey())
			if expected != actual {
				t.Fatalf("invalid via path (%s->%s)[%d]: %s, expecting %s", nameStart, nodeKey, i, actual.GetKey(), expected.GetKey())
			}
		}
	}
}
