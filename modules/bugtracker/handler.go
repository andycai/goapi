package bugtracker

import (
	"strconv"

	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// Project handlers
func ListProjectsHandler(c *fiber.Ctx) error {
	projects, err := GetService().ListProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"projects": projects,
	})
}

func CreateProjectHandler(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := GetService().CreateProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

func UpdateProjectHandler(c *fiber.Ctx) error {
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

	if err := GetService().UpdateProject(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

func GetProjectHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	project, err := GetService().GetProject(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"project": project,
	})
}

// Iteration handlers
func ListIterationsHandler(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Query("project_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid project ID",
		})
	}

	iterations, err := GetService().ListProjectIterations(projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iterations": iterations,
	})
}

func CreateIterationHandler(c *fiber.Ctx) error {
	var iteration models.Iteration
	if err := c.BodyParser(&iteration); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := GetService().CreateIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

func UpdateIterationHandler(c *fiber.Ctx) error {
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

	if err := GetService().UpdateIteration(&iteration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

func GetIterationHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid iteration ID",
		})
	}

	iteration, err := GetService().GetIteration(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"iteration": iteration,
	})
}

// Issue handlers
func ListIssuesHandler(c *fiber.Ctx) error {
	projectID, _ := strconv.ParseInt(c.Query("project_id"), 10, 64)
	iterationID, _ := strconv.ParseInt(c.Query("iteration_id"), 10, 64)

	var issues []models.Issue
	var err error

	if iterationID > 0 {
		issues, err = GetService().ListIterationIssues(iterationID)
	} else if projectID > 0 {
		issues, err = GetService().ListProjectIssues(projectID)
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

func CreateIssueHandler(c *fiber.Ctx) error {
	var issue models.Issue
	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := GetService().CreateIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

func UpdateIssueHandler(c *fiber.Ctx) error {
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

	if err := GetService().UpdateIssue(&issue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

func GetIssueHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	issue, err := GetService().GetIssue(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"issue": issue,
	})
}

// Comment handlers
func ListCommentsHandler(c *fiber.Ctx) error {
	issueID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue ID",
		})
	}

	comments, err := GetService().ListIssueComments(issueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comments": comments,
	})
}

func CreateCommentHandler(c *fiber.Ctx) error {
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

	if err := GetService().CreateComment(&comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comment": comment,
	})
}
