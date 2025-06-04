package string

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank 检查字符串是否为空白（空字符串或仅含空白字符）
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Truncate 截断字符串，如果超过最大长度则添加后缀
func Truncate(s string, maxLen int, suffix string) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	return string([]rune(s)[:maxLen]) + suffix
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Capitalize 将字符串首字母大写
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// PadLeft 左侧填充字符串到指定长度
func PadLeft(s string, padChar rune, totalLength int) string {
	count := totalLength - utf8.RuneCountInString(s)
	if count <= 0 {
		return s
	}
	return strings.Repeat(string(padChar), count) + s
}

// PadRight 右侧填充字符串到指定长度
func PadRight(s string, padChar rune, totalLength int) string {
	count := totalLength - utf8.RuneCountInString(s)
	if count <= 0 {
		return s
	}
	return s + strings.Repeat(string(padChar), count)
}

// IsNumeric 检查字符串是否只包含数字
func IsNumeric(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 检查字符串是否只包含字母
func IsAlpha(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查字符串是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// ToSnakeCase 将驼峰命名转换为蛇形命名
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// ToCamelCase 将蛇形命名转换为驼峰命名
func ToCamelCase(s string) string {
	var result strings.Builder
	upper := false
	for _, r := range s {
		if r == '_' {
			upper = true
		} else {
			if upper {
				result.WriteRune(unicode.ToUpper(r))
				upper = false
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}

// ToBase64 将字符串编码为Base64
func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// FromBase64 将Base64字符串解码
func FromBase64(s string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseInt 将字符串解析为整数，可设置默认值
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseFloat 将字符串解析为浮点数，可设置默认值
func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// IsValidEmail 验证电子邮件格式
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ContainsAny 检查字符串是否包含任一指定字符串
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// EscapeHTML 转义HTML特殊字符
func EscapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(s)
}

// JoinInt 将整数切片以指定分隔符连接为字符串
func JoinInt(ints []int, separator string) string {
	if len(ints) == 0 {
		return ""
	}
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, separator)
}

// RemoveNonPrintable 移除所有不可打印字符
func RemoveNonPrintable(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToJSON 将对象转换为JSON字符串
func ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 将JSON字符串转换为对象
func FromJSON(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
