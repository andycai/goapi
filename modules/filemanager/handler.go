package filemanager

import (
	"github.com/gofiber/fiber/v2"
)

// ListHandler handles the request to list files and directories
func ListHandler(c *fiber.Ctx) error {
	path := c.Query("path", "./")
	files, err := GetService().List(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}

// UploadHandler handles file upload requests
func UploadHandler(c *fiber.Ctx) error {
	path := c.FormValue("path", "./")
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer src.Close()

	err = GetService().Upload(path, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
	})
}

// CreateHandler handles the creation of new files and directories
func CreateHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")
	isDir := c.FormValue("is_dir", "false") == "true"

	err := GetService().Create(path, isDir)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Created successfully",
	})
}

// DeleteHandler handles file and directory deletion
func DeleteHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")

	err := GetService().Delete(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

// RenameHandler handles file and directory renaming
func RenameHandler(c *fiber.Ctx) error {
	oldPath := c.FormValue("old_path")
	newPath := c.FormValue("new_path")

	err := GetService().Rename(oldPath, newPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Renamed successfully",
	})
}

// MoveHandler handles moving files and directories
func MoveHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := GetService().Move(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Moved successfully",
	})
}

// CopyHandler handles copying files and directories
func CopyHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := GetService().Copy(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Copied successfully",
	})
}

// DownloadHandler handles file downloads
func DownloadHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	info, err := GetService().GetInfo(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if info.IsDir {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot download directories",
		})
	}

	return c.Download(path)
}

// InfoHandler returns detailed information about a file or directory
func InfoHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	info, err := GetService().GetInfo(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"info": info,
	})
}
