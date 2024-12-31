package collections

import (
	"container/heap"
	"errors"
)

// 最小优先队列（优先级低的在前）
// minPQ := NewMinPriorityQueue[string]()
// minPQ.Enqueue("task1", 1)
// minPQ.Enqueue("task2", 2)
// minPQ.Enqueue("task3", 3)
// 出队顺序：task1(1) -> task2(2) -> task3(3)

// MinPriorityQueue 最小优先队列实现（优先级值越小越优先）
type MinPriorityQueue[T any] struct {
	items []*PriorityItem[T]
}

// NewMinPriorityQueue 创建新的最小优先队列
func NewMinPriorityQueue[T any]() *MinPriorityQueue[T] {
	pq := &MinPriorityQueue[T]{
		items: make([]*PriorityItem[T], 0),
	}
	heap.Init(pq)
	return pq
}

// Len 实现 heap.Interface
func (pq *MinPriorityQueue[T]) Len() int {
	return len(pq.items)
}

// Less 实现 heap.Interface，优先级低的在前
func (pq *MinPriorityQueue[T]) Less(i, j int) bool {
	return pq.items[i].Priority < pq.items[j].Priority
}

// Swap 实现 heap.Interface
func (pq *MinPriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// Push 实现 heap.Interface
func (pq *MinPriorityQueue[T]) Push(x any) {
	n := len(pq.items)
	item := x.(*PriorityItem[T])
	item.index = n
	pq.items = append(pq.items, item)
}

// Pop 实现 heap.Interface
func (pq *MinPriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.index = -1 // 标记为已移除
	pq.items = old[0 : n-1]
	return item
}

// Enqueue 入队
func (pq *MinPriorityQueue[T]) Enqueue(value T, priority int) *PriorityItem[T] {
	item := &PriorityItem[T]{
		Value:    value,
		Priority: priority,
	}
	heap.Push(pq, item)
	return item
}

// Dequeue 出队
func (pq *MinPriorityQueue[T]) Dequeue() (T, int, error) {
	var zero T
	if len(pq.items) == 0 {
		return zero, 0, errors.New("priority queue is empty")
	}
	item := heap.Pop(pq).(*PriorityItem[T])
	return item.Value, item.Priority, nil
}

// Peek 查看队首元素
func (pq *MinPriorityQueue[T]) Peek() (T, int, error) {
	var zero T
	if len(pq.items) == 0 {
		return zero, 0, errors.New("priority queue is empty")
	}
	item := pq.items[0]
	return item.Value, item.Priority, nil
}

// IsEmpty 检查队列是否为空
func (pq *MinPriorityQueue[T]) IsEmpty() bool {
	return len(pq.items) == 0
}

// Size 获取队列大小
func (pq *MinPriorityQueue[T]) Size() int {
	return len(pq.items)
}

// Clear 清空队列
func (pq *MinPriorityQueue[T]) Clear() {
	pq.items = make([]*PriorityItem[T], 0)
	heap.Init(pq)
}

// UpdatePriority 更新元素优先级
func (pq *MinPriorityQueue[T]) UpdatePriority(item *PriorityItem[T], priority int) {
	item.Priority = priority
	heap.Fix(pq, item.index)
}

// ForEach 遍历优先队列（按优先级顺序）
func (pq *MinPriorityQueue[T]) ForEach(fn func(value T, priority int) bool) {
	// 创建临时队列，不影响原队列
	tempQueue := &MinPriorityQueue[T]{
		items: make([]*PriorityItem[T], len(pq.items)),
	}
	copy(tempQueue.items, pq.items)
	heap.Init(tempQueue)

	// 按优先级顺序遍历
	for !tempQueue.IsEmpty() {
		value, priority, _ := tempQueue.Dequeue()
		if !fn(value, priority) {
			break
		}
	}
}

// ToSlice 转换为切片（按优先级顺序）
func (pq *MinPriorityQueue[T]) ToSlice() []T {
	result := make([]T, 0, pq.Len())
	pq.ForEach(func(value T, _ int) bool {
		result = append(result, value)
		return true
	})
	return result
}

// FromSlice 从切片创建优先队列
func (pq *MinPriorityQueue[T]) FromSlice(items []T, priorities []int) error {
	if len(items) != len(priorities) {
		return errors.New("items and priorities length mismatch")
	}

	pq.Clear()
	for i, item := range items {
		pq.Enqueue(item, priorities[i])
	}
	return nil
}
