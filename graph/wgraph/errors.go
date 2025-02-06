package wgraph

import (
	"errors"
	"fmt"
)

var (
	ErrConnExists                 = errors.New("node connection already exists")
	ErrDijkstraNegativeWeightEdge = errors.New("a Dijkstra edge's weight must not be negative")
)

func wrapErrConnExists[W Weight](n1, n2 NodeWeighted[W]) error {
	return fmt.Errorf("existing edge between node '%s' to '%s': %w", n1.GetKey(), n2.GetKey(), ErrConnExists)
}
