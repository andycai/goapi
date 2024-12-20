package filemanager

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var filemanagerService *FilemanagerService

type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"is_dir"`
	Mode      string    `json:"mode"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

type FilemanagerService struct {
	rootPath string
}

func InitService() {
	filemanagerService = &FilemanagerService{
		rootPath: "./",
	}
}

func GetService() *FilemanagerService {
	return filemanagerService
}

// isValidPath checks if the path is safe for file operations
func (s *FilemanagerService) isValidPath(path string) bool {
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

	rootAbs, err := filepath.Abs(s.rootPath)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, rootAbs)
}

// List returns a list of files and directories in the specified path
func (s *FilemanagerService) List(path string) ([]FileInfo, error) {
	if !s.isValidPath(path) {
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

// Upload handles file upload to the specified path
func (s *FilemanagerService) Upload(path string, file io.Reader, filename string) error {
	if !s.isValidPath(path) {
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

// Create creates a new directory or file
func (s *FilemanagerService) Create(path string, isDir bool) error {
	if !s.isValidPath(path) {
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

// Delete removes a file or directory
func (s *FilemanagerService) Delete(path string) error {
	if !s.isValidPath(path) {
		return errors.New("invalid path")
	}
	return os.RemoveAll(path)
}

// Rename renames a file or directory
func (s *FilemanagerService) Rename(oldPath, newPath string) error {
	if !s.isValidPath(oldPath) || !s.isValidPath(newPath) {
		return errors.New("invalid path")
	}
	return os.Rename(oldPath, newPath)
}

// Move moves a file or directory to a new location
func (s *FilemanagerService) Move(sourcePath, destPath string) error {
	if !s.isValidPath(sourcePath) || !s.isValidPath(destPath) {
		return errors.New("invalid path")
	}
	return os.Rename(sourcePath, destPath)
}

// Copy copies a file or directory to a new location
func (s *FilemanagerService) Copy(sourcePath, destPath string) error {
	if !s.isValidPath(sourcePath) || !s.isValidPath(destPath) {
		return errors.New("invalid path")
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return s.copyDir(sourcePath, destPath)
	}
	return s.copyFile(sourcePath, destPath)
}

// copyFile copies a single file
func (s *FilemanagerService) copyFile(sourcePath, destPath string) error {
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
func (s *FilemanagerService) copyDir(sourcePath, destPath string) error {
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
			err = s.copyDir(sourceName, destName)
		} else {
			err = s.copyFile(sourceName, destName)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// GetInfo returns detailed information about a file or directory
func (s *FilemanagerService) GetInfo(path string) (*FileInfo, error) {
	if !s.isValidPath(path) {
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
