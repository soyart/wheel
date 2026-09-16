package wgraph

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

// This file is the regression suite for a non-determinism investigation
// into DijkstraShortestPathFrom. Across that investigation, four different
// implementations were tried, in order:
//
//  1. push-before-relax (the original bug): a neighbor is pushed onto the
//     heap using its stale, pre-relaxation value, so the heap's internal
//     tree can end up ordered inconsistently with the nodes' real values.
//  2. push-after-relax (commit 60b7555): fixes the entry it just wrote, but
//     an older, already-queued entry for the same node is never
//     repositioned, so it can still go stale.
//  3. two-pass relax-then-push, marking a node visited at push time
//     (commit 62a8067): fixes same-node double-pushes, but a node is only
//     ever pushed once, so an improvement to its value that arrives after
//     that single push has no way to propagate onward.
//  4. no visited tracking at all (commit f8bbba1): fixes the above, but
//     never terminates on any graph with a cycle, and blows up
//     exponentially even on some acyclic graphs (a "diamond cascade"
//     replays every distinct path independently).
//
// The fix that actually closes all of the above pushes an immutable
// (node, dist) snapshot per relaxation and discards a popped entry whose
// dist no longer matches the node's live value - see dijkstraQueueItem in
// dijkstra_impl.go. The tests below hold that fix in place: the oracle
// fuzz test catches wrong-answer regressions (variants 1-3's failure
// mode), and the two dedicated tests catch termination/performance
// regressions (variant 4's failure mode).

// fuzzEdge is a directed edge used by the tests in this file.
type fuzzEdge struct {
	from, to int
	weight   int
}

// randomGraph builds a random directed graph over 6-15 nodes. With
// allowCycles false, edges only go from a lower index to a higher one,
// guaranteeing a DAG; with it true, edges are drawn from any ordered pair,
// so cycles (including 2-cycles from a pair of opposite edges) are
// possible.
func randomGraph(seed int64, allowCycles bool) (n int, edges []fuzzEdge) {
	r := rand.New(rand.NewSource(seed))
	n = 6 + r.Intn(10)

	for i := range n {
		for j := range n {
			if i == j {
				continue
			}
			if !allowCycles && i > j {
				continue
			}
			if r.Intn(100) < 35 {
				edges = append(edges, fuzzEdge{from: i, to: j, weight: 1 + r.Intn(9)})
			}
		}
	}

	return n, edges
}

// referenceShortestPaths computes shortest distances from node 0 with a
// plain O(V^2) selection-based Dijkstra and no priority queue, as an
// oracle independent of the real implementation under test. It is correct
// for any non-negative-weight graph, cyclic or not.
func referenceShortestPaths(n int, edges []fuzzEdge) []int {
	const inf = 1 << 30

	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = 0

	adj := make([][]fuzzEdge, n)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], e)
	}

	done := make([]bool, n)
	for range dist {
		u := -1
		for i := range n {
			if !done[i] && (u == -1 || dist[i] < dist[u]) {
				u = i
			}
		}
		if u == -1 || dist[u] == inf {
			break
		}
		done[u] = true
		for _, e := range adj[u] {
			if nd := dist[u] + e.weight; nd < dist[e.to] {
				dist[e.to] = nd
			}
		}
	}

	return dist
}

func fuzzNodeName(i int) string {
	return string(rune('A'+i%26)) + string(rune('0'+i/26))
}

func buildFuzzGraph(n int, edges []fuzzEdge) ([]*NodeDijkstraImpl[int], GraphDijkstra[int]) {
	const inf = 1 << 30

	nodes := make([]*NodeDijkstraImpl[int], n)
	for i := range nodes {
		nodes[i] = &NodeDijkstraImpl[int]{
			NodeWeightedImpl: NodeWeightedImpl[int]{Name: fuzzNodeName(i), ValueOrCost: inf},
		}
	}

	g := NewDijkstraGraph[int](true)
	for _, node := range nodes {
		g.AddNode(node)
	}
	for _, e := range edges {
		_ = g.AddEdgeWeightOrDistance(nodes[e.from], nodes[e.to], e.weight)
	}

	return nodes, g
}

// TestDijkstraAgainstReferenceOracle fuzzes DijkstraShortestPathFrom
// against an independent O(V^2) oracle across thousands of random small
// graphs, both acyclic and with cycles. Every fix attempt that came before
// dijkstraQueueItem's snapshot approach fails a small fraction of these
// cases; this is the test that would have caught each of them.
func TestDijkstraAgainstReferenceOracle(t *testing.T) {
	const trialsPerMode = 20000

	for _, allowCycles := range []bool{false, true} {
		for seed := int64(0); seed < trialsPerMode; seed++ {
			n, edges := randomGraph(seed, allowCycles)
			want := referenceShortestPaths(n, edges)

			nodes, g := buildFuzzGraph(n, edges)
			result := g.DijkstraShortestPathFrom(nodes[0])

			for i, node := range nodes {
				if got := node.GetValue(); got != want[i] {
					t.Fatalf(
						"cycles=%v seed=%d: node %s: got %d, want %d\nedges=%v",
						allowCycles, seed, node.GetKey(), got, want[i], edges,
					)
				}
			}

			assertPathsConsistent(t, allowCycles, seed, nodes, edges, result)
		}
	}
}

// assertPathsConsistent checks that every (node -> via) entry in the
// returned Paths map corresponds to a real edge whose weight exactly
// explains the recorded distance, catching a correct-looking distance
// table paired with a broken/inconsistent path reconstruction.
func assertPathsConsistent(
	t *testing.T,
	allowCycles bool,
	seed int64,
	nodes []*NodeDijkstraImpl[int],
	edges []fuzzEdge,
	result *DijstraShortestPath[int],
) {
	t.Helper()

	indexOf := make(map[NodeDijkstra[int]]int, len(nodes))
	for i, node := range nodes {
		indexOf[node] = i
	}

	weight := make(map[[2]int]int, len(edges))
	for _, e := range edges {
		weight[[2]int{e.from, e.to}] = e.weight
	}

	for node, via := range result.Paths {
		ni, vi := indexOf[node], indexOf[via]
		w, ok := weight[[2]int{vi, ni}]
		if !ok {
			t.Fatalf("cycles=%v seed=%d: Paths records %s via %s, but no such edge exists",
				allowCycles, seed, node.GetKey(), via.GetKey())
		}
		if node.GetValue() != via.GetValue()+w {
			t.Fatalf("cycles=%v seed=%d: %s recorded via %s, but dist %d != %d+%d",
				allowCycles, seed, node.GetKey(), via.GetKey(), node.GetValue(), via.GetValue(), w)
		}
	}
}

// TestDijkstraTerminatesOnCycle is a narrow regression test for commit
// f8bbba1 ("dont track visited"), which never terminated on any graph
// containing a cycle: every pop unconditionally re-pushed all of its
// neighbors, with nothing to ever stop it. The simplest possible cycle - a
// single edge modeled in both directions - hangs it forever.
func TestDijkstraTerminatesOnCycle(t *testing.T) {
	a := &NodeDijkstraImpl[int]{NodeWeightedImpl: NodeWeightedImpl[int]{Name: "A"}}
	b := &NodeDijkstraImpl[int]{NodeWeightedImpl: NodeWeightedImpl[int]{Name: "B", ValueOrCost: 1 << 30}}

	g := NewDijkstraGraph[int](true)
	g.AddNode(a)
	g.AddNode(b)
	if err := g.AddEdgeWeightOrDistance(a, b, 1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdgeWeightOrDistance(b, a, 1); err != nil {
		t.Fatal(err)
	}

	done := make(chan *DijstraShortestPath[int], 1)
	go func() { done <- g.DijkstraShortestPathFrom(a) }()

	select {
	case result := <-done:
		if len(result.Paths) != 1 {
			t.Fatalf("expected exactly 1 path entry (B via A), got %d", len(result.Paths))
		}
		if b.GetValue() != 1 {
			t.Fatalf("expected B's distance to be 1, got %d", b.GetValue())
		}
	case <-time.After(3 * time.Second):
		t.Fatal("DijkstraShortestPathFrom did not terminate within 3s on a 2-node cycle")
	}
}

// TestDijkstraStaysLinearOnDeepCascade is a narrow regression test for the
// other failure mode of commit f8bbba1: with nothing to stop a node from
// being reprocessed, a "diamond cascade" (each layer of 2 nodes fully
// connected to the next layer's 2 nodes) makes the number of heap pushes
// double per layer, because every distinct path is replayed independently.
// At depth 50 that is over 10^15 pushes; a correct implementation, which
// finalizes each node exactly once, handles it in microseconds regardless
// of depth.
func TestDijkstraStaysLinearOnDeepCascade(t *testing.T) {
	const depth = 50

	start := &NodeDijkstraImpl[int]{NodeWeightedImpl: NodeWeightedImpl[int]{Name: "start"}}
	g := NewDijkstraGraph[int](true)
	g.AddNode(start)

	prevLayer := []NodeDijkstra[int]{start}
	nodeCount := 1
	for range depth {
		layer := make([]NodeDijkstra[int], 2)
		for k := range layer {
			node := &NodeDijkstraImpl[int]{
				NodeWeightedImpl: NodeWeightedImpl[int]{Name: fuzzNodeName(nodeCount), ValueOrCost: 1 << 30},
			}
			nodeCount++
			g.AddNode(node)
			layer[k] = node
		}
		for _, p := range prevLayer {
			for _, ln := range layer {
				if err := g.AddEdgeWeightOrDistance(p, ln, 1); err != nil {
					t.Fatal(err)
				}
			}
		}
		prevLayer = layer
	}

	done := make(chan *DijstraShortestPath[int], 1)
	go func() { done <- g.DijkstraShortestPathFrom(start) }()

	select {
	case result := <-done:
		if len(result.Paths) != nodeCount-1 {
			t.Fatalf("expected %d path entries, got %d", nodeCount-1, len(result.Paths))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DijkstraShortestPathFrom did not finish within 2s on a 50-layer diamond cascade (101 nodes) - looks like an exponential blowup regressed")
	}
}

// TestDijkstraRejectsNegativeWeight documents and locks in that
// AddEdgeWeightOrDistance refuses negative edge weights.
func TestDijkstraRejectsNegativeWeight(t *testing.T) {
	a := &NodeDijkstraImpl[int]{NodeWeightedImpl: NodeWeightedImpl[int]{Name: "A"}}
	b := &NodeDijkstraImpl[int]{NodeWeightedImpl: NodeWeightedImpl[int]{Name: "B"}}

	g := NewDijkstraGraph[int](true)
	g.AddNode(a)
	g.AddNode(b)

	err := g.AddEdgeWeightOrDistance(a, b, -1)
	if !errors.Is(err, ErrDijkstraNegativeWeightEdge) {
		t.Fatalf("expected ErrDijkstraNegativeWeightEdge, got %v", err)
	}
}
