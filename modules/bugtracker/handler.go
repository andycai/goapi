package bugtracker

import (
	"strconv"

	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// Project handlers
func listProjects(c *fiber.Ctx) error {
	projects, err := srv.ListProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"projects": projects,
	})
}

func createProject(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := srv.CreateProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

func updateProject(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	project.ID = id

	if err := srv.UpdateProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

// getProject get project by ID
func getProject(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	project, err := srv.GetProject(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

// listIterations list iterations by project ID
func listIterations(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Query("project_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	iterations, err := srv.ListProjectIterations(projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iterations": iterations,
	})
}

// createIteration create iteration
func createIteration(c *fiber.Ctx) error {
	var iteration models.Iteration
	if err := c.BodyParser(&iteration); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := srv.CreateIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// updateIteration update iteration
func updateIteration(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid iteration ID",
		})
	}

	var iteration models.Iteration
	if err := c.BodyParser(&iteration); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	iteration.ID = id

	if err := srv.UpdateIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// getIteration get iteration by ID
func getIteration(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid iteration ID",
		})
	}

	iteration, err := srv.GetIteration(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// listIssues list issues by project_id or iteration_id
// if both are provided, iteration_id will take precedence
// if neither are provided, return 400 Bad Request
// if project_id is provided, return all issues in the project
func listIssues(c *fiber.Ctx) error {
	projectID, _ := strconv.ParseInt(c.Query("project_id"), 10, 64)
	iterationID, _ := strconv.ParseInt(c.Query("iteration_id"), 10, 64)

	var issues []models.Issue
	var err error

	if iterationID > 0 {
		issues, err = srv.ListIterationIssues(iterationID)
	} else if projectID > 0 {
		issues, err = srv.ListProjectIssues(projectID)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "project_id or iteration_id is required",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issues": issues,
	})
}

// createIssue create issue
func createIssue(c *fiber.Ctx) error {
	var issue models.Issue
	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := srv.CreateIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// updateIssue update issue
func updateIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	var issue models.Issue
	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	issue.ID = id

	if err := srv.UpdateIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// getIssue get issue by ID
func getIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	issue, err := srv.GetIssue(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// listComments list comments by issue ID
func listComments(c *fiber.Ctx) error {
	issueID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	comments, err := srv.ListIssueComments(issueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comments": comments,
	})
}

// createComment create comment
func createComment(c *fiber.Ctx) error {
	issueID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	comment.IssueID = issueID

	if err := srv.CreateComment(&comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comment": comment,
	})
}
