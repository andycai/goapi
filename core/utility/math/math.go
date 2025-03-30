package math

import (
	"math"
	"math/rand"
)

// Min 返回两个整数中的较小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max 返回两个整数中的较大值
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt64 返回两个int64中的较小值
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MaxInt64 返回两个int64中的较大值
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 返回两个float64中的较小值
func MinFloat64(a, b float64) float64 {
	return math.Min(a, b)
}

// MaxFloat64 返回两个float64中的较大值
func MaxFloat64(a, b float64) float64 {
	return math.Max(a, b)
}

// Abs 返回整数的绝对值
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// AbsInt64 返回int64的绝对值
func AbsInt64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// AbsFloat64 返回float64的绝对值
func AbsFloat64(a float64) float64 {
	return math.Abs(a)
}

// Round 将浮点数四舍五入到最接近的整数
func Round(a float64) int {
	return int(math.Round(a))
}

// RoundToFloat 将浮点数四舍五入到指定小数位数
func RoundToFloat(a float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(a*shift) / shift
}

// Floor 向下取整
func Floor(a float64) int {
	return int(math.Floor(a))
}

// Ceil 向上取整
func Ceil(a float64) int {
	return int(math.Ceil(a))
}

// Pow 计算x的y次方
func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

// Sqrt 计算平方根
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// Cbrt 计算立方根
func Cbrt(x float64) float64 {
	return math.Cbrt(x)
}

// Clamp 将值限制在指定范围内
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat64 将浮点值限制在指定范围内
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsEven 判断整数是否为偶数
func IsEven(n int) bool {
	return n%2 == 0
}

// IsOdd 判断整数是否为奇数
func IsOdd(n int) bool {
	return n%2 != 0
}

// IsPowerOfTwo 判断整数是否为2的幂
func IsPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// Gcd 计算最大公约数
func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Lcm 计算最小公倍数
func Lcm(a, b int) int {
	return a / Gcd(a, b) * b
}

// Factorial 计算阶乘
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// Fibonacci 计算斐波那契数列的第n项
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciIterative 迭代方式计算斐波那契数列的第n项（避免递归导致的堆栈溢出）
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Sin 计算正弦值
func Sin(x float64) float64 {
	return math.Sin(x)
}

// Cos 计算余弦值
func Cos(x float64) float64 {
	return math.Cos(x)
}

// Tan 计算正切值
func Tan(x float64) float64 {
	return math.Tan(x)
}

// Asin 计算反正弦值
func Asin(x float64) float64 {
	return math.Asin(x)
}

// Acos 计算反余弦值
func Acos(x float64) float64 {
	return math.Acos(x)
}

// Atan 计算反正切值
func Atan(x float64) float64 {
	return math.Atan(x)
}

// Atan2 计算y/x的反正切值
func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// DegreesToRadians 角度转弧度
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// RadiansToDegrees 弧度转角度
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// Log 计算自然对数
func Log(x float64) float64 {
	return math.Log(x)
}

// Log10 计算以10为底的对数
func Log10(x float64) float64 {
	return math.Log10(x)
}

// Log2 计算以2为底的对数
func Log2(x float64) float64 {
	return math.Log2(x)
}

// RandomInt 生成范围内的随机整数 [min, max)
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomFloat64 生成范围内的随机浮点数 [min, max)
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomBool 生成随机布尔值
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// Lerp 线性插值
func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Distance 计算两点间距离
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceSquared 计算两点间距离的平方（避免开平方操作）
func DistanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

// Normalize 将向量归一化
func Normalize(x, y float64) (float64, float64) {
	length := math.Sqrt(x*x + y*y)
	if length < 1e-10 {
		return 0, 0
	}
	return x / length, y / length
}

// Sum 求整数切片的和
func Sum(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

// SumFloat64 求浮点数切片的和
func SumFloat64(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

// Average 求整数切片的平均值
func Average(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(Sum(values)) / float64(len(values))
}

// AverageFloat64 求浮点数切片的平均值
func AverageFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return SumFloat64(values) / float64(len(values))
}
