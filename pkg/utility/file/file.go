package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Exists 检查文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir 检查指定路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 检查指定路径是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// CreateDir 创建目录，包括不存在的父目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ReadFile 读取整个文件内容为字节数组
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileString 读取整个文件内容为字符串
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile 将字节数组写入文件
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if !Exists(dir) {
		if err := CreateDir(dir); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

// WriteFileString 将字符串写入文件
func WriteFileString(path string, content string) error {
	return WriteFile(path, []byte(content))
}

// AppendFile 追加内容到文件
func AppendFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// AppendFileString 追加字符串到文件
func AppendFileString(path string, content string) error {
	return AppendFile(path, []byte(content))
}

// ReadLines 读取文件所有行
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

// ReadJSON 从文件读取JSON数据到结构体
func ReadJSON(path string, v interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSON 将结构体以JSON格式写入文件
func WriteJSON(path string, v interface{}, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	return WriteFile(path, data)
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
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
	if !IsFile(path) {
		return errors.New("指定路径不是文件")
	}
	return os.Remove(path)
}

// SafeDelete 安全删除文件或目录（检查存在性）
func SafeDelete(path string) error {
	if !Exists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetModTime 获取文件修改时间
func GetModTime(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}
