package filemanager

import (
	"github.com/gofiber/fiber/v2"
)

// listFilesHandler handles the request to list files and directories
func listFilesHandler(c *fiber.Ctx) error {
	path := c.Query("path", "./")
	files, err := listFiles(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}

// uploadFileHandler handles file upload requests
func uploadFileHandler(c *fiber.Ctx) error {
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

	err = uploadFile(path, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
	})
}

// createHandler handles the creation of new files and directories
func createHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")
	isDir := c.FormValue("is_dir", "false") == "true"

	err := create(path, isDir)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Created successfully",
	})
}

// deleteHandler handles file and directory deletion
func deleteHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")

	err := delete(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

// renameHandler handles file and directory renaming
func renameHandler(c *fiber.Ctx) error {
	oldPath := c.FormValue("old_path")
	newPath := c.FormValue("new_path")

	err := rename(oldPath, newPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Renamed successfully",
	})
}

// moveHandler handles moving files and directories
func moveHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := move(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Moved successfully",
	})
}

// copyHandler handles copying files and directories
func copyHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := copy(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Copied successfully",
	})
}

// downloadHandler handles file downloads
func downloadHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	info, err := getInfo(path)
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

// infoHandler returns detailed information about a file or directory
func infoHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	info, err := getInfo(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"info": info,
	})
}
