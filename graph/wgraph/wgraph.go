package wgraph

import (
	"cmp"

	"github.com/soyart/wheel/graph"
)

type Weight cmp.Ordered

// GraphWeighted has T as node values and edge weight.
type GraphWeighted[N NodeWeighted[T], E EdgeWeighted[T, N], T Weight] graph.Graph[N, E, T]
