package git

import (
	"github.com/gofiber/fiber/v2"
)

type GitRequest struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Branch   string `json:"branch"`
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
	Create   bool   `json:"create"`
	Hard     bool   `json:"hard"`
	Pop      bool   `json:"pop"`
	Limit    int    `json:"limit"`
}

// clone handles Git clone requests
func clone(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Clone(req.URL, req.Path, req.Branch, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Repository cloned successfully",
	})
}

// pull handles Git pull requests
func pull(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Pull(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pull successful",
	})
}

// push handles Git push requests
func push(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Push(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Push successful",
	})
}

// status handles Git status requests
func status(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	status, err := srv.Status(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

// log handles Git log requests
func log(c *fiber.Ctx) error {
	path := c.Query("path")
	limit := c.QueryInt("limit", 10) // Default to 10 entries if not specified

	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	log, err := srv.Log(path, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"log": log,
	})
}

// commit handles Git commit requests
func commit(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Commit(req.Path, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Changes committed successfully",
	})
}

// checkout handles Git checkout requests
func checkout(c *fiber.Ctx) error {
	var req GitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" || req.Branch == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path and branch are required",
		})
	}

	err := srv.Checkout(req.Path, req.Branch, req.Create)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
	})
}

// branch handles Git branch requests
func branch(c *fiber.Ctx) error {
	var req GitRequest
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

	result, err := srv.Branch(req.Path, req.Create, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
	})
}

// merge handles Git merge requests
func merge(c *fiber.Ctx) error {
	var req GitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Path == "" || req.Branch == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path and branch are required",
		})
	}

	err := srv.Merge(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Merge successful",
	})
}

// reset handles Git reset requests
func reset(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Reset(req.Path, req.Hard)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reset successful",
	})
}

// stash handles Git stash requests
func stash(c *fiber.Ctx) error {
	var req GitRequest
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

	err := srv.Stash(req.Path, req.Pop)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stash operation successful",
	})
}
