package imagemanager

import (
	"github.com/gofiber/fiber/v2"
)

// listFilesHandler handles the request to list images
func listFilesHandler(c *fiber.Ctx) error {
	path := c.Query("path", "./")
	images, err := srv.List(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"images": images,
	})
}

// uploadHandler handles image upload requests
func uploadHandler(c *fiber.Ctx) error {
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

	err = srv.Upload(path, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image uploaded successfully",
	})
}

// deleteHandler handles image deletion
func deleteHandler(c *fiber.Ctx) error {
	path := c.FormValue("path")

	err := srv.Delete(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}

// renameHandler handles image renaming
func renameHandler(c *fiber.Ctx) error {
	oldPath := c.FormValue("old_path")
	newPath := c.FormValue("new_path")

	err := srv.Rename(oldPath, newPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image renamed successfully",
	})
}

// moveHandler handles moving images
func moveHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := srv.Move(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image moved successfully",
	})
}

// copyHandler handles copying images
func copyHandler(c *fiber.Ctx) error {
	sourcePath := c.FormValue("source_path")
	destPath := c.FormValue("dest_path")

	err := srv.Copy(sourcePath, destPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image copied successfully",
	})
}

// infoHandler returns detailed information about an image
func infoHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	info, err := srv.GetInfo(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"info": info,
	})
}

// thumbnailHandler serves image thumbnails
func thumbnailHandler(c *fiber.Ctx) error {
	path := c.Query("path")

	thumbnailPath, err := srv.GetThumbnail(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendFile(thumbnailPath)
}

// viewHandler serves the original image
func viewHandler(c *fiber.Ctx) error {
	path := c.Query("path")
	fullPath := srv.rootPath + "/" + path

	return c.SendFile(fullPath)
}
