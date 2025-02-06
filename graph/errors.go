package graph

import "errors"

var ErrEdgeWeightNotNull = errors.New("found edge weight in unweighted graph")
