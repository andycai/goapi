package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CreateDirIfNotExists 如果目录不存在则创建
func CreateDirIfNotExists(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

// GetFileName 获取文件名（不含扩展名）
func GetFileName(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}

// ReadFile 读取文件内容
func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile 写入文件内容
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// AppendToFile 追加内容到文件
func AppendToFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
	return os.Remove(path)
}

// GetFileMD5 计算文件的MD5值
func GetFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ReadLines 读取文件的所有行
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines 将多行内容写入文件
func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

// IsTextFile 判断是否为文本文件
func IsTextFile(path string) bool {
	ext := GetFileExtension(path)
	textExtensions := map[string]bool{
		".txt":  true,
		".log":  true,
		".json": true,
		".xml":  true,
		".yaml": true,
		".yml":  true,
		".md":   true,
		".csv":  true,
		".ini":  true,
		".conf": true,
		".go":   true,
		".js":   true,
		".css":  true,
		".html": true,
		".sql":  true,
	}
	return textExtensions[ext]
}

// IsBinaryFile 判断是否为二进制文件
func IsBinaryFile(path string) bool {
	ext := GetFileExtension(path)
	binaryExtensions := map[string]bool{
		".exe":   true,
		".dll":   true,
		".so":    true,
		".dylib": true,
		".bin":   true,
		".dat":   true,
		".zip":   true,
		".tar":   true,
		".gz":    true,
		".rar":   true,
		".7z":    true,
		".jpg":   true,
		".jpeg":  true,
		".png":   true,
		".gif":   true,
		".bmp":   true,
		".ico":   true,
		".mp3":   true,
		".mp4":   true,
		".avi":   true,
		".mov":   true,
		".pdf":   true,
		".doc":   true,
		".docx":  true,
		".xls":   true,
		".xlsx":  true,
		".ppt":   true,
		".pptx":  true,
	}
	return binaryExtensions[ext]
}

// GetFilePermissions 获取文件权限
func GetFilePermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// SetFilePermissions 设置文件权限
func SetFilePermissions(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// GetFileOwner 获取文件所有者
func GetFileOwner(path string) (uint32, uint32, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}
	stat := info.Sys().(*syscall.Stat_t)
	return stat.Uid, stat.Gid, nil
}

// SetFileOwner 设置文件所有者
func SetFileOwner(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}
