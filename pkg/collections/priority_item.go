package collections

// PriorityItem 优先队列项
type PriorityItem[T any] struct {
	Value    T
	Priority int
	index    int
}
