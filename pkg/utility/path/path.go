package path

import (
	"os"
	"path/filepath"
	"strings"
)

// IsVallid 是否有效路径
func IsValid(path string) bool {
	// 清理和规范化路径
	cleanPath := filepath.Clean(path)

	// 检查可疑模式
	suspicious := []string{
		"..",
		"~",
		"$",
		"|",
		">",
		"<",
		"&",
		"`",
	}

	for _, pattern := range suspicious {
		if strings.Contains(cleanPath, pattern) {
			return false
		}
	}

	return true
}

// Join 连接路径片段
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Base 获取路径的最后一个元素
func Base(path string) string {
	return filepath.Base(path)
}

// Dir 获取路径的目录部分
func Dir(path string) string {
	return filepath.Dir(path)
}

// Ext 获取路径的扩展名
func Ext(path string) string {
	return filepath.Ext(path)
}

// Clean 清理路径，消除多余的分隔符和上级引用
func Clean(path string) string {
	return filepath.Clean(path)
}

// Abs 获取绝对路径
func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// Rel 获取相对路径
func Rel(basepath, targetpath string) (string, error) {
	return filepath.Rel(basepath, targetpath)
}

// IsAbs 判断路径是否为绝对路径
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Split 分割路径为目录和文件名
func Split(path string) (dir, file string) {
	return filepath.Split(path)
}

// VolumeName 获取路径的卷名
func VolumeName(path string) string {
	return filepath.VolumeName(path)
}

// Match 路径通配符匹配
func Match(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

// ChangeExt 更改文件扩展名
func ChangeExt(path, newExt string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path + newExt
	}
	return path[:len(path)-len(ext)] + newExt
}

// RemoveExt 移除文件扩展名
func RemoveExt(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path
	}
	return path[:len(path)-len(ext)]
}

// GetTempDir 获取系统临时目录
func GetTempDir() string {
	return os.TempDir()
}

// GetWorkingDir 获取当前工作目录
func GetWorkingDir() (string, error) {
	return os.Getwd()
}

// GetHomeDir 获取用户主目录
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// EnsureTrailingSlash 确保路径以斜杠结尾
func EnsureTrailingSlash(path string) string {
	if path == "" {
		return "/"
	}
	if path[len(path)-1] != '/' && path[len(path)-1] != '\\' {
		return path + string(os.PathSeparator)
	}
	return path
}

// RemoveTrailingSlash 移除路径末尾的斜杠
func RemoveTrailingSlash(path string) string {
	if path == "" || path == "/" || path == "\\" {
		return path
	}
	if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
		return path[:len(path)-1]
	}
	return path
}

// IsSubPath 检查childPath是否是parentPath的子路径
func IsSubPath(parentPath, childPath string) (bool, error) {
	// 获取绝对路径
	absParent, err := filepath.Abs(parentPath)
	if err != nil {
		return false, err
	}
	absChild, err := filepath.Abs(childPath)
	if err != nil {
		return false, err
	}

	// 规范化路径
	absParent = filepath.Clean(absParent)
	absChild = filepath.Clean(absChild)

	// 确保有斜杠
	absParent = EnsureTrailingSlash(absParent)

	// 检查子路径
	return strings.HasPrefix(absChild, absParent), nil
}
