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

// CloneHandler handles Git clone requests
func CloneHandler(c *fiber.Ctx) error {
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

	err := GetService().Clone(req.URL, req.Path, req.Branch, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Repository cloned successfully",
	})
}

// PullHandler handles Git pull requests
func PullHandler(c *fiber.Ctx) error {
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

	err := GetService().Pull(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pull successful",
	})
}

// PushHandler handles Git push requests
func PushHandler(c *fiber.Ctx) error {
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

	err := GetService().Push(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Push successful",
	})
}

// StatusHandler handles Git status requests
func StatusHandler(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	status, err := GetService().Status(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

// LogHandler handles Git log requests
func LogHandler(c *fiber.Ctx) error {
	path := c.Query("path")
	limit := c.QueryInt("limit", 10) // Default to 10 entries if not specified

	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Path is required",
		})
	}

	log, err := GetService().Log(path, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"log": log,
	})
}

// CommitHandler handles Git commit requests
func CommitHandler(c *fiber.Ctx) error {
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

	err := GetService().Commit(req.Path, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Changes committed successfully",
	})
}

// CheckoutHandler handles Git checkout requests
func CheckoutHandler(c *fiber.Ctx) error {
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

	err := GetService().Checkout(req.Path, req.Branch, req.Create)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
	})
}

// BranchHandler handles Git branch requests
func BranchHandler(c *fiber.Ctx) error {
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

	result, err := GetService().Branch(req.Path, req.Create, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
	})
}

// MergeHandler handles Git merge requests
func MergeHandler(c *fiber.Ctx) error {
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

	err := GetService().Merge(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Merge successful",
	})
}

// ResetHandler handles Git reset requests
func ResetHandler(c *fiber.Ctx) error {
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

	err := GetService().Reset(req.Path, req.Hard)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reset successful",
	})
}

// StashHandler handles Git stash requests
func StashHandler(c *fiber.Ctx) error {
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

	err := GetService().Stash(req.Path, req.Pop)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stash operation successful",
	})
}
