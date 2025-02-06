package list

type BasicList[T any] interface {
	Push(x T)
	PushSlice(values []T)
	Pop() *T
	Len() int
	IsEmpty() bool
}
