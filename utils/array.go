package utils

import (
	"math/rand"
	"sort"
	"time"
)

// Contains 检查数组是否包含某个元素
func Contains[T comparable](arr []T, target T) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

// IndexOf 获取元素在数组中的索引
func IndexOf[T comparable](arr []T, target T) int {
	for i, item := range arr {
		if item == target {
			return i
		}
	}
	return -1
}

// LastIndexOf 获取元素在数组中最后出现的索引
func LastIndexOf[T comparable](arr []T, target T) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == target {
			return i
		}
	}
	return -1
}

// Remove 从数组中移除指定元素
func Remove[T comparable](arr []T, target T) []T {
	result := make([]T, 0)
	for _, item := range arr {
		if item != target {
			result = append(result, item)
		}
	}
	return result
}

// RemoveAt 从数组中移除指定索引的元素
func RemoveAt[T any](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		return arr
	}
	return append(arr[:index], arr[index+1:]...)
}

// RemoveAll 从数组中移除所有指定元素
func RemoveAll[T comparable](arr []T, target T) []T {
	result := make([]T, 0)
	for _, item := range arr {
		if item != target {
			result = append(result, item)
		}
	}
	return result
}

// Unique 获取数组的唯一值
func Unique[T comparable](arr []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, item := range arr {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// Reverse 反转数组
func Reverse[T any](arr []T) []T {
	result := make([]T, len(arr))
	for i, j := 0, len(arr)-1; i < len(arr); i, j = i+1, j-1 {
		result[i] = arr[j]
	}
	return result
}

// Shuffle 打乱数组
func Shuffle[T any](arr []T) []T {
	result := make([]T, len(arr))
	copy(result, arr)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

// Sort 排序数组
func Sort[T any](arr []T, less func(i, j T) bool) []T {
	result := make([]T, len(arr))
	copy(result, arr)
	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})
	return result
}

// Filter 过滤数组
func Filter[T any](arr []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range arr {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map 映射数组
func Map[T any, U any](arr []T, mapper func(T) U) []U {
	result := make([]U, len(arr))
	for i, item := range arr {
		result[i] = mapper(item)
	}
	return result
}

// Reduce 归约数组
func Reduce[T any, U any](arr []T, reducer func(U, T) U, initial U) U {
	result := initial
	for _, item := range arr {
		result = reducer(result, item)
	}
	return result
}

// Find 查找满足条件的第一个元素
func Find[T any](arr []T, predicate func(T) bool) (T, bool) {
	for _, item := range arr {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex 查找满足条件的第一个元素的索引
func FindIndex[T any](arr []T, predicate func(T) bool) int {
	for i, item := range arr {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// Every 检查是否所有元素都满足条件
func Every[T any](arr []T, predicate func(T) bool) bool {
	for _, item := range arr {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Some 检查是否存在元素满足条件
func Some[T any](arr []T, predicate func(T) bool) bool {
	for _, item := range arr {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Count 统计满足条件的元素个数
func Count[T any](arr []T, predicate func(T) bool) int {
	count := 0
	for _, item := range arr {
		if predicate(item) {
			count++
		}
	}
	return count
}

// GroupBy 按条件对数组元素进行分组
func GroupBy[T any, K comparable](arr []T, keySelector func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range arr {
		key := keySelector(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Chunk 将数组分割成指定大小的块
func Chunk[T any](arr []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	result := make([][]T, 0)
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}
	return result
}

// Flatten 展平嵌套数组
func Flatten[T any](arr [][]T) []T {
	result := make([]T, 0)
	for _, subArr := range arr {
		result = append(result, subArr...)
	}
	return result
}

// Zip 将多个数组按索引组合
func Zip[T any](arrs ...[]T) [][]T {
	if len(arrs) == 0 {
		return [][]T{}
	}

	minLen := len(arrs[0])
	for _, arr := range arrs[1:] {
		if len(arr) < minLen {
			minLen = len(arr)
		}
	}

	result := make([][]T, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = make([]T, len(arrs))
		for j, arr := range arrs {
			result[i][j] = arr[i]
		}
	}
	return result
}

// Intersection 获取多个数组的交集
func Intersection[T comparable](arrs ...[]T) []T {
	if len(arrs) == 0 {
		return []T{}
	}

	seen := make(map[T]int)
	for _, item := range arrs[0] {
		seen[item] = 1
	}

	for _, arr := range arrs[1:] {
		currSeen := make(map[T]bool)
		for _, item := range arr {
			if count := seen[item]; count == len(seen) && !currSeen[item] {
				seen[item]++
				currSeen[item] = true
			}
		}
	}

	result := make([]T, 0)
	for item, count := range seen {
		if count == len(arrs) {
			result = append(result, item)
		}
	}
	return result
}

// Union 获取多个数组的并集
func Union[T comparable](arrs ...[]T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)

	for _, arr := range arrs {
		for _, item := range arr {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
	}
	return result
}

// Difference 获取数组的差集
func Difference[T comparable](arr1 []T, arr2 []T) []T {
	seen := make(map[T]bool)
	for _, item := range arr2 {
		seen[item] = true
	}

	result := make([]T, 0)
	for _, item := range arr1 {
		if !seen[item] {
			result = append(result, item)
		}
	}
	return result
}
