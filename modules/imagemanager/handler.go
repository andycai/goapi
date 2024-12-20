package imagemanager

import (
	"github.com/gofiber/fiber/v2"
)

// ListHandler handles the request to list images
func ListHandler(c *fiber.Ctx) error {
	path := c.Query("path", "./")
	images, err := GetService().List(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"images": images,
	})
}

// UploadHandler handles image upload requests
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
		"message": "Image uploaded successfully",
	})
}

// DeleteHandler handles image deletion
func DeleteHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")

	err := GetService().Delete(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}

// RenameHandler handles image renaming
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
		"message": "Image renamed successfully",
	})
}

// MoveHandler handles moving images
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
		"message": "Image moved successfully",
	})
}

// CopyHandler handles copying images
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
		"message": "Image copied successfully",
	})
}

// InfoHandler returns detailed information about an image
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

// ThumbnailHandler serves image thumbnails
func ThumbnailHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	thumbnailPath, err := GetService().GetThumbnail(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendFile(thumbnailPath)
}

// ViewHandler serves the original image
func ViewHandler(c *fiber.Ctx) error {
	path := c.Query("path")
	fullPath := GetService().rootPath + "/" + path

	return c.SendFile(fullPath)
}
