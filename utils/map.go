package utils

// Keys 获取map的所有键
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取map的所有值
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Entries 获取map的所有键值对
func Entries[K comparable, V any](m map[K]V) [][2]interface{} {
	entries := make([][2]interface{}, 0, len(m))
	for k, v := range m {
		entries = append(entries, [2]interface{}{k, v})
	}
	return entries
}

// FromEntries 从键值对数组创建map
func FromEntries[K comparable, V any](entries [][2]interface{}) map[K]V {
	m := make(map[K]V)
	for _, entry := range entries {
		if k, ok := entry[0].(K); ok {
			if v, ok := entry[1].(V); ok {
				m[k] = v
			}
		}
	}
	return m
}

// Merge 合并多个map
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Filter 过滤map
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Map 映射map
func MapMap[K comparable, V any, U any](m map[K]V, mapper func(K, V) U) map[K]U {
	result := make(map[K]U)
	for k, v := range m {
		result[k] = mapper(k, v)
	}
	return result
}

// ForEach 遍历map
func ForEach[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Find 查找满足条件的第一个键值对
func FindInMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, V, bool) {
	for k, v := range m {
		if predicate(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Every 检查是否所有键值对都满足条件
func EveryMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if !predicate(k, v) {
			return false
		}
	}
	return true
}

// Some 检查是否存在键值对满足条件
func SomeMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// Count 统计满足条件的键值对个数
func CountMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) int {
	count := 0
	for k, v := range m {
		if predicate(k, v) {
			count++
		}
	}
	return count
}

// GroupBy 按条件对map进行分组
func GroupByMap[K comparable, V any, G comparable](m map[K]V, keySelector func(K, V) G) map[G]map[K]V {
	result := make(map[G]map[K]V)
	for k, v := range m {
		group := keySelector(k, v)
		if result[group] == nil {
			result[group] = make(map[K]V)
		}
		result[group][k] = v
	}
	return result
}

// Invert 反转map的键值
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K)
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Pick 从map中选择指定的键
func Pick[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Omit 从map中排除指定的键
func Omit[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		exclude := false
		for _, key := range keys {
			if k == key {
				exclude = true
				break
			}
		}
		if !exclude {
			result[k] = v
		}
	}
	return result
}

// HasKey 检查map是否包含指定的键
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// HasValue 检查map是否包含指定的值
func HasValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

// Clear 清空map
func Clear[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

// Clone 克隆map
func Clone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Update 更新map中的值
func Update[K comparable, V any](m map[K]V, key K, updater func(V) V) bool {
	if v, ok := m[key]; ok {
		m[key] = updater(v)
		return true
	}
	return false
}

// GetOrDefault 获取map中的值，如果不存在则返回默认值
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// GetOrSet 获取map中的值，如果不存在则设置并返回默认值
func GetOrSet[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = defaultValue
	return defaultValue
}

// GetOrCompute 获取map中的值，如果不存在则计算、设置并返回新值
func GetOrCompute[K comparable, V any](m map[K]V, key K, computer func() V) V {
	if v, ok := m[key]; ok {
		return v
	}
	v := computer()
	m[key] = v
	return v
}
