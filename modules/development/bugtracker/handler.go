package bugtracker

import (
	"strconv"

	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// Project handlers
func listProjectsHandler(c *fiber.Ctx) error {
	projects, err := getProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"projects": projects,
	})
}

func createProjectHandler(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := createProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

func updateProjectHandler(c *fiber.Ctx) error {
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

	if err := updateProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

// getProject get project by ID
func getProjectHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	project, err := getProject(id)
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
func listIterationsHandler(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Query("project_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	iterations, err := getProjectIterations(projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iterations": iterations,
	})
}

// createIterationHandler create iteration
func createIterationHandler(c *fiber.Ctx) error {
	var iteration models.Iteration
	if err := c.BodyParser(&iteration); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := createIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// updateIterationHandler update iteration
func updateIterationHandler(c *fiber.Ctx) error {
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

	if err := updateIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// getIterationHandler get iteration by ID
func getIterationHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid iteration ID",
		})
	}

	iteration, err := getIteration(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// listIssuesHandler list issues by project_id or iteration_id
// if both are provided, iteration_id will take precedence
// if neither are provided, return 400 Bad Request
// if project_id is provided, return all issues in the project
func listIssuesHandler(c *fiber.Ctx) error {
	projectID, _ := strconv.ParseInt(c.Query("project_id"), 10, 64)
	iterationID, _ := strconv.ParseInt(c.Query("iteration_id"), 10, 64)

	var issues []models.Issue
	var err error

	if iterationID > 0 {
		issues, err = getIterationIssues(iterationID)
	} else if projectID > 0 {
		issues, err = getProjectIssues(projectID)
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

// createIssueHandler create issue
func createIssueHandler(c *fiber.Ctx) error {
	var issue models.Issue
	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := createIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// updateIssueHandler update issue
func updateIssueHandler(c *fiber.Ctx) error {
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

	if err := updateIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// getIssueHandler get issue by ID
func getIssueHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	issue, err := getIssue(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// listCommentsHandler list comments by issue ID
func listCommentsHandler(c *fiber.Ctx) error {
	issueID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	comments, err := getIssueComments(issueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comments": comments,
	})
}

// createCommentHandler create comment
func createCommentHandler(c *fiber.Ctx) error {
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

	if err := createComment(&comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comment": comment,
	})
}
