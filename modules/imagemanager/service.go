package imagemanager

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

var imagemanagerService *ImagemanagerService

type ImageInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Format    string    `json:"format"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

type ImagemanagerService struct {
	rootPath     string
	thumbnailDir string
}

func InitService() {
	imagemanagerService = &ImagemanagerService{
		rootPath:     "./uploads/images",
		thumbnailDir: "./uploads/thumbnails",
	}
	// Create directories if they don't exist
	os.MkdirAll(imagemanagerService.rootPath, 0755)
	os.MkdirAll(imagemanagerService.thumbnailDir, 0755)
}

func GetService() *ImagemanagerService {
	return imagemanagerService
}

// isValidPath checks if the path is safe for file operations
func (s *ImagemanagerService) isValidPath(path string) bool {
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

// List returns a list of images in the specified path
func (s *ImagemanagerService) List(path string) ([]ImageInfo, error) {
	fullPath := filepath.Join(s.rootPath, path)
	if !s.isValidPath(fullPath) {
		return nil, errors.New("invalid path")
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var images []ImageInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			continue // Skip non-image files
		}

		imgInfo, err := s.GetInfo(filepath.Join(path, info.Name()))
		if err != nil {
			continue
		}

		images = append(images, *imgInfo)
	}

	return images, nil
}

// Upload handles image upload to the specified path
func (s *ImagemanagerService) Upload(path string, file io.Reader, filename string) error {
	fullPath := filepath.Join(s.rootPath, path, filename)
	if !s.isValidPath(fullPath) {
		return errors.New("invalid path")
	}

	// Read the entire file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Verify that it's a valid image
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return errors.New("invalid image format: " + format)
	}

	// Create the target file
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the data to the file
	_, err = io.Copy(f, bytes.NewReader(data))
	if err != nil {
		return err
	}

	// Generate thumbnail
	return s.generateThumbnail(fullPath)
}

// Delete removes an image and its thumbnail
func (s *ImagemanagerService) Delete(path string) error {
	fullPath := filepath.Join(s.rootPath, path)
	if !s.isValidPath(fullPath) {
		return errors.New("invalid path")
	}

	// Delete thumbnail first
	s.deleteThumbnail(fullPath)

	// Delete the image file
	return os.Remove(fullPath)
}

// Rename renames an image and its thumbnail
func (s *ImagemanagerService) Rename(oldPath, newPath string) error {
	oldFullPath := filepath.Join(s.rootPath, oldPath)
	newFullPath := filepath.Join(s.rootPath, newPath)

	if !s.isValidPath(oldFullPath) || !s.isValidPath(newFullPath) {
		return errors.New("invalid path")
	}

	// Rename thumbnail first
	s.renameThumbnail(oldFullPath, newFullPath)

	// Rename the image file
	return os.Rename(oldFullPath, newFullPath)
}

// Move moves an image and its thumbnail to a new location
func (s *ImagemanagerService) Move(sourcePath, destPath string) error {
	sourceFullPath := filepath.Join(s.rootPath, sourcePath)
	destFullPath := filepath.Join(s.rootPath, destPath)

	if !s.isValidPath(sourceFullPath) || !s.isValidPath(destFullPath) {
		return errors.New("invalid path")
	}

	// Move thumbnail first
	s.moveThumbnail(sourceFullPath, destFullPath)

	// Move the image file
	return os.Rename(sourceFullPath, destFullPath)
}

// Copy copies an image and its thumbnail to a new location
func (s *ImagemanagerService) Copy(sourcePath, destPath string) error {
	sourceFullPath := filepath.Join(s.rootPath, sourcePath)
	destFullPath := filepath.Join(s.rootPath, destPath)

	if !s.isValidPath(sourceFullPath) || !s.isValidPath(destFullPath) {
		return errors.New("invalid path")
	}

	// Copy the image file
	source, err := os.Open(sourceFullPath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	// Generate thumbnail for the copied image
	return s.generateThumbnail(destFullPath)
}

// GetInfo returns detailed information about an image
func (s *ImagemanagerService) GetInfo(path string) (*ImageInfo, error) {
	fullPath := filepath.Join(s.rootPath, path)
	if !s.isValidPath(fullPath) {
		return nil, errors.New("invalid path")
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	// Open the image file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image config
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		Name:      info.Name(),
		Path:      path,
		Size:      info.Size(),
		Width:     config.Width,
		Height:    config.Height,
		Format:    format,
		ModTime:   info.ModTime(),
		Extension: filepath.Ext(info.Name()),
	}, nil
}

// GetThumbnail returns a thumbnail of the image
func (s *ImagemanagerService) GetThumbnail(path string) (string, error) {
	fullPath := filepath.Join(s.rootPath, path)
	if !s.isValidPath(fullPath) {
		return "", errors.New("invalid path")
	}

	thumbnailPath := s.getThumbnailPath(fullPath)
	if _, err := os.Stat(thumbnailPath); err != nil {
		// Generate thumbnail if it doesn't exist
		err = s.generateThumbnail(fullPath)
		if err != nil {
			return "", err
		}
	}

	return thumbnailPath, nil
}

// generateThumbnail creates a thumbnail for the given image
func (s *ImagemanagerService) generateThumbnail(imagePath string) error {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image to create a thumbnail
	thumbnail := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	// Create the thumbnail file
	thumbnailPath := s.getThumbnailPath(imagePath)
	os.MkdirAll(filepath.Dir(thumbnailPath), 0755)
	thumbFile, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()

	// Encode the thumbnail
	switch format {
	case "jpeg":
		return jpeg.Encode(thumbFile, thumbnail, &jpeg.Options{Quality: 85})
	case "png":
		return png.Encode(thumbFile, thumbnail)
	default:
		return jpeg.Encode(thumbFile, thumbnail, &jpeg.Options{Quality: 85})
	}
}

// getThumbnailPath returns the path where the thumbnail should be stored
func (s *ImagemanagerService) getThumbnailPath(imagePath string) string {
	relPath, _ := filepath.Rel(s.rootPath, imagePath)
	return filepath.Join(s.thumbnailDir, relPath)
}

// deleteThumbnail removes the thumbnail for the given image
func (s *ImagemanagerService) deleteThumbnail(imagePath string) {
	thumbnailPath := s.getThumbnailPath(imagePath)
	os.Remove(thumbnailPath)
}

// renameThumbnail renames the thumbnail for the given image
func (s *ImagemanagerService) renameThumbnail(oldPath, newPath string) {
	oldThumbnail := s.getThumbnailPath(oldPath)
	newThumbnail := s.getThumbnailPath(newPath)
	os.Rename(oldThumbnail, newThumbnail)
}

// moveThumbnail moves the thumbnail for the given image
func (s *ImagemanagerService) moveThumbnail(sourcePath, destPath string) {
	sourceThumbnail := s.getThumbnailPath(sourcePath)
	destThumbnail := s.getThumbnailPath(destPath)
	os.Rename(sourceThumbnail, destThumbnail)
}
