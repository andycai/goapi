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

// cloneHandler handles Git cloneHandler requests
func cloneHandler(c *fiber.Ctx) error {
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

	err := Clone(req.URL, req.Path, req.Branch, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Repository cloned successfully",
	})
}

// pullHandler handles Git pullHandler requests
func pullHandler(c *fiber.Ctx) error {
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

	err := Pull(req.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pull successful",
	})
}

// pushHandler handles Git pushHandler requests
func pushHandler(c *fiber.Ctx) error {
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

	err := Push(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Push successful",
	})
}

// statusHandler handles Git statusHandler requests
func statusHandler(c *fiber.Ctx) error {
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

// logHandler handles Git logHandler requests
func logHandler(c *fiber.Ctx) error {
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

// commitHandler handles Git commitHandler requests
func commitHandler(c *fiber.Ctx) error {
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

	err := Commit(req.Path, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Changes committed successfully",
	})
}

// checkoutHandler handles Git checkoutHandler requests
func checkoutHandler(c *fiber.Ctx) error {
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

	err := Checkout(req.Path, req.Branch, req.Create)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
	})
}

// branchHandler handles Git branchHandler requests
func branchHandler(c *fiber.Ctx) error {
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

	result, err := Branch(req.Path, req.Create, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
	})
}

// mergeHandler handles Git mergeHandler requests
func mergeHandler(c *fiber.Ctx) error {
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

	err := Merge(req.Path, req.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Merge successful",
	})
}

// resetHandler handles Git resetHandler requests
func resetHandler(c *fiber.Ctx) error {
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

	err := Reset(req.Path, req.Hard)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reset successful",
	})
}

// stashHandler handles Git stashHandler requests
func stashHandler(c *fiber.Ctx) error {
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

	err := Stash(req.Path, req.Pop)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stash operation successful",
	})
}
