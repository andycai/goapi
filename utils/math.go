package utils

import (
	"math"
	"math/rand"
	"time"
)

// Round 四舍五入到指定小数位
func Round(x float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Round(x*p) / p
}

// Floor 向下取整到指定小数位
func Floor(x float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(x*p) / p
}

// Ceil 向上取整到指定小数位
func Ceil(x float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Ceil(x*p) / p
}

// Max 返回最大值
func Max[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Min 返回最小值
func Min[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Clamp 将值限制在指定范围内
func Clamp[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// RandomInt 生成指定范围内的随机整数
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// RandomFloat 生成指定范围内的随机浮点数
func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// Lerp 线性插值
func Lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}

// Distance 计算两点之间的距离
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// Factorial 计算阶乘
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// IsPrime 判断是否为素数
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// GCD 计算最大公约数
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM 计算最小公倍数
func LCM(a, b int) int {
	return a * b / GCD(a, b)
}
