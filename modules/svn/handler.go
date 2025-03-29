package svn

import (
	"github.com/gofiber/fiber/v2"
)

type SVNRequest struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
	Limit    int    `json:"limit"`
}

// checkout handles SVN checkout requests
func checkout(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.URL == "" || req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL and path are required",
		})
	}

	err := Checkout(req.URL, req.Path, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
	})
}

// update handles SVN update requests
func update(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	err := Update(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Update successful",
	})
}

// commit handles SVN commit requests
func commit(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path and message are required",
		})
	}

	err := Commit(req.Path, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Commit successful",
	})
}

// status handles SVN status requests
func status(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	status, err := Status(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

// info handles SVN info requests
func info(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	info, err := Info(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"info": info,
	})
}

// log handles SVN log requests
func log(c *fiber.Ctx) error {
	path := c.Query("path")
	limit := c.QueryInt("limit", 10) // Default to 10 entries if not specified

	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	log, err := Log(path, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"log": log,
	})
}

// revert handles SVN revert requests
func revert(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	err := Revert(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Revert successful",
	})
}

// add handles SVN add requests
func add(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	err := Add(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Add successful",
	})
}

// delete handles SVN delete requests
func delete(c *fiber.Ctx) error {
	var req SVNRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	err := Delete(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Delete successful",
	})
}
