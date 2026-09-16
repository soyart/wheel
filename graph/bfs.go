package graph

import (
	"github.com/soyart/wheel"
	"github.com/soyart/wheel/list"
)

// BFS calls BFSSearch and uses its output to call BFSShortestPathReconstruct.
// It then returns the shortest path (slice of nodes), the number of hops it takes from `src` to `dst`,
// and a true boolean value if there's a path from `src` to `dst`.
// Otherwise, a nil slice, -1, and false is returned if there's no such path.
func BFS[T, E any](
	g Graph[Node[T], E, any],
	src Node[T],
	dst Node[T],
) (
	[]Node[T],
	int,
	bool,
) {
	rawPath, found := BFSSearchGeneric(g, src, dst)
	if !found {
		return nil, -1, false
	}
	shortestPath, hops := BFSShortestPathReconstruct(rawPath, src, dst)
	return shortestPath, hops, found
}

// BFSSearchGeneric traverses the graph in BFS manner, and collecting VFS traversal information in a map `prev`.
// It returns the map, and a boolean value denoting if dst was reachable from src
func BFSSearchGeneric[T, E any](
	g Graph[Node[T], E, any],
	src Node[T],
	dst Node[T],
) (
	map[Node[T]]Node[T],
	bool,
) {
	q := list.NewQueueSafe[Node[T]]()
	q.Push(src)

	visited := make(map[Node[T]]bool)
	prev := make(map[Node[T]]Node[T])

	var found bool
	for !q.IsEmpty() {
		popped := q.Pop()
		if popped == nil {
			panic("popped nil - should not happen")
		}

		current := *popped
		neighbors := g.GetNodeNeighbors(current)
		for _, neighbor := range neighbors {
			if visited[neighbor] {
				continue
			}
			visited[neighbor] = true

			if neighbor == dst {
				found = true
			}

			q.Push(neighbor)
			prev[neighbor] = current
		}
	}
	return prev, found
}

// BFSShortestPathReconstruct returns inclusive path
// from src to dst, with number of hops
func BFSShortestPathReconstruct[T any](
	backwardPath map[Node[T]]Node[T],
	src Node[T],
	dst Node[T],
) (
	[]Node[T],
	int,
) {
	current := dst
	shortestPath := []Node[T]{current}

	var hops int
	if current == src {
		return shortestPath, hops
	}

	for current != src {
		next, found := backwardPath[current]
		if !found {
			break
		}
		shortestPath = append(shortestPath, next)
		current = next
		hops++
	}

	wheel.ReverseInPlace(shortestPath)
	return shortestPath, hops
}

// BFSHashMapGraphV1 wraps BFS to make it easier for users to call BFS with HashMapGraphV1.
func BFSHashMapGraphV1[T any](
	g HashMapGraphV1[T],
	src Node[T],
	dst Node[T],
) (
	[]Node[T],
	int,
	bool,
) {
	return BFS(g.(Graph[Node[T], Node[T], any]), src, dst)
}
