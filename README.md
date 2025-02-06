# Wheel

Wheel is a minimal Go library providing frequently or prevalently used
functionalities like sorting algorithms, data structures, and utilities.

Its purpose is to prevent me (and you) from reinventing the wheel.

In general, wheel is divided into 4 modules:

- [`github.com/soyart/wheel`](./) - for basic stuff like sorting and finding max/min

- [`github.com/soyart/wheel/list`](./list/) - for basic containers like stacks, queues, and sets

- [`github.com/soyart/wheel/tree`](./tree/)- for trees, like BST and heap

- [`github.com/soyart/wheel/graph`](./graph/) - for graph types

  > Note that `graph` is very hard to use, and will need a rework.
  > It also has a flaky tests from undeterminism of hash maps backing the graph.