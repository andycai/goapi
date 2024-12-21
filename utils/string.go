package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// MD5 计算字符串的MD5值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// TruncateString 截断字符串到指定长度，并添加省略号
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// RemoveSpaces 移除字符串中的所有空白字符
func RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// IsEmail 验证邮箱地址格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsURL 验证URL格式
func IsURL(url string) bool {
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}[a-zA-Z0-9\-\._~:/?#\[\]@!\$&'\(\)\*\+,;=]*$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// Capitalize 首字母大写
func Capitalize(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// CamelCase 将字符串转换为驼峰命名
func CamelCase(str string) string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for i := 1; i < len(words); i++ {
		words[i] = Capitalize(words[i])
	}

	return strings.Join(words, "")
}

// SnakeCase 将字符串转换为蛇形命名
func SnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// KebabCase 将字符串转换为短横线命名
func KebabCase(str string) string {
	return strings.ReplaceAll(SnakeCase(str), "_", "-")
}

// ReverseString 反转字符串
func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CountWords 统计单词数量
func CountWords(str string) int {
	return len(strings.Fields(str))
}

// ExtractEmails 从文本中提取所有邮箱地址
func ExtractEmails(text string) []string {
	pattern := `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// ExtractURLs 从文本中提取所有URL
func ExtractURLs(text string) []string {
	pattern := `(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}[a-zA-Z0-9\-\._~:/?#\[\]@!\$&'\(\)\*\+,;=]*`
	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// Slugify 将字符串转换为URL友好的格式
func Slugify(str string) string {
	// 转换为小写
	str = strings.ToLower(str)

	// 替换非字母数字字符为短横线
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	str = reg.ReplaceAllString(str, "-")

	// 移除首尾的短横线
	str = strings.Trim(str, "-")

	return str
}

// WrapText 按指定宽度换行文本
func WrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	lineLength := 0

	for i, word := range words {
		wordLength := len(word)
		if lineLength+wordLength > width {
			result.WriteString("\n")
			lineLength = 0
		} else if i > 0 {
			result.WriteString(" ")
			lineLength++
		}

		result.WriteString(word)
		lineLength += wordLength
	}

	return result.String()
}

// FormatTemplate 简单的模板格式化
func FormatTemplate(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}
	return result
}
