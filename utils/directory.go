package utils

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// CreateDirectory 创建目录
func CreateDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RemoveDirectory 删除目录
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

// CopyDirectory 复制目录
func CopyDirectory(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 遍历源目录
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		// 如果是目录，创建对应的目标目录
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 复制文件
		return copyFile(path, dstPath)
	})
}

// copyFile 复制文件（内部使用）
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 获取源文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 设置目标文件权限
	return os.Chmod(dst, srcInfo.Mode())
}

// ListDirectory 列出目录内容
func ListDirectory(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

// ListDirectoryWithFilter 列出目录内容（带过滤器）
func ListDirectoryWithFilter(path string, filter func(string) bool) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// GetDirectorySize 获取目录大小
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// IsEmptyDirectory 检查目录是否为空
func IsEmptyDirectory(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// GetSubdirectories 获取所有子目录
func GetSubdirectories(path string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

// GetFiles 获取所有文件
func GetFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// GetFilesWithExtension 获取指定扩展名的文件
func GetFilesWithExtension(path string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// CleanDirectory 清空目录
func CleanDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// MoveDirectory 移动目录
func MoveDirectory(src, dst string) error {
	if err := CopyDirectory(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// GetDirectoryPermissions 获取目录权限
func GetDirectoryPermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// SetDirectoryPermissions 设置目录权限
func SetDirectoryPermissions(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// GetDirectoryOwner 获取目录所有者
func GetDirectoryOwner(path string) (uint32, uint32, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}
	stat := info.Sys().(*syscall.Stat_t)
	return stat.Uid, stat.Gid, nil
}

// SetDirectoryOwner 设置目录所有者
func SetDirectoryOwner(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}
