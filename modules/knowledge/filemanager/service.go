package filemanager

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var rootPath string

func initService() {
	rootPath = "./"
}

// isValidPath checks if the path is safe for file operations
func isValidPath(path string) bool {
	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Check for suspicious patterns
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

	// Ensure the path is within the root directory
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return false
	}

	rootAbs, err := filepath.Abs(rootPath)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, rootAbs)
}

// listFiles returns a list of files and directories in the specified path
func listFiles(path string) ([]FileInfo, error) {
	if !isValidPath(path) {
		return nil, errors.New("invalid path")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:      info.Name(),
			Path:      filepath.Join(path, info.Name()),
			Size:      info.Size(),
			IsDir:     info.IsDir(),
			Mode:      info.Mode().String(),
			ModTime:   info.ModTime(),
			Extension: filepath.Ext(info.Name()),
		})
	}

	// Sort files: directories first, then files
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

// uploadFile handles file upload to the specified path
func uploadFile(path string, file io.Reader, filename string) error {
	if !isValidPath(path) {
		return errors.New("invalid path")
	}

	targetPath := filepath.Join(path, filename)
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	return err
}

// create creates a new directory or file
func create(path string, isDir bool) error {
	if !isValidPath(path) {
		return errors.New("invalid path")
	}

	if isDir {
		return os.MkdirAll(path, 0755)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// delete removes a file or directory
func delete(path string) error {
	if !isValidPath(path) {
		return errors.New("invalid path")
	}
	return os.RemoveAll(path)
}

// rename renames a file or directory
func rename(oldPath, newPath string) error {
	if !isValidPath(oldPath) || !isValidPath(newPath) {
		return errors.New("invalid path")
	}
	return os.Rename(oldPath, newPath)
}

// move moves a file or directory to a new location
func move(sourcePath, destPath string) error {
	if !isValidPath(sourcePath) || !isValidPath(destPath) {
		return errors.New("invalid path")
	}
	return os.Rename(sourcePath, destPath)
}

// copy copies a file or directory to a new location
func copy(sourcePath, destPath string) error {
	if !isValidPath(sourcePath) || !isValidPath(destPath) {
		return errors.New("invalid path")
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return copyDir(sourcePath, destPath)
	}
	return copyFile(sourcePath, destPath)
}

// copyFile copies a single file
func copyFile(sourcePath, destPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

// copyDir copies a directory recursively
func copyDir(sourcePath, destPath string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destPath, sourceInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceName := filepath.Join(sourcePath, entry.Name())
		destName := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			err = copyDir(sourceName, destName)
		} else {
			err = copyFile(sourceName, destName)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// getInfo returns detailed information about a file or directory
func getInfo(path string) (*FileInfo, error) {
	if !isValidPath(path) {
		return nil, errors.New("invalid path")
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:      info.Name(),
		Path:      path,
		Size:      info.Size(),
		IsDir:     info.IsDir(),
		Mode:      info.Mode().String(),
		ModTime:   info.ModTime(),
		Extension: filepath.Ext(info.Name()),
	}, nil
}
