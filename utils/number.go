package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatThousands 格式化数字为千分位
func FormatThousands(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in = in[1:]
		out[0] = '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			break
		}
		k++
		if k == 3 {
			j--
			out[j] = ','
			k = 0
		}
	}
	return string(out)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
	var units = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	fsize := float64(size)
	for fsize >= 1024 && i < len(units)-1 {
		fsize /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", fsize, units[i])
}

// ParseInt 安全的将字符串转换为整数
func ParseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultVal
}

// ParseFloat 安全的将字符串转换为浮点数
func ParseFloat(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return defaultVal
}

// ToFixed 将浮点数格式化为指定小数位数的字符串
func ToFixed(num float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, num)
}

// IsNumeric 判断字符串是否为数字
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	dotCount := 0
	for _, c := range s {
		if c == '.' {
			dotCount++
			if dotCount > 1 {
				return false
			}
		} else if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// FormatPercent 格式化百分比
func FormatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value*100)
}

// FormatCurrency 格式化货币
func FormatCurrency(amount float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, amount)
}

// ParseVersion 解析版本号
func ParseVersion(version string) (major, minor, patch int, err error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid version format")
	}

	major = ParseInt(parts[0], 0)
	minor = ParseInt(parts[1], 0)
	patch = ParseInt(parts[2], 0)
	return
}

// CompareVersion 比较两个版本号
func CompareVersion(v1, v2 string) int {
	major1, minor1, patch1, err1 := ParseVersion(v1)
	major2, minor2, patch2, err2 := ParseVersion(v2)

	if err1 != nil || err2 != nil {
		return 0
	}

	if major1 != major2 {
		return major1 - major2
	}
	if minor1 != minor2 {
		return minor1 - minor2
	}
	return patch1 - patch2
}
