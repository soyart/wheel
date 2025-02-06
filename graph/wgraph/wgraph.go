package wgraph

import (
	"golang.org/x/exp/constraints"

	"github.com/soyart/wheel/graph"
)

type Weight constraints.Ordered

// GraphWeighted has T as node values and edge weight.
type GraphWeighted[N NodeWeighted[T], E EdgeWeighted[T, N], T Weight] graph.Graph[N, E, T]
