package bugtracker

import (
	"errors"

	"github.com/andycai/goapi/models"
)

// createProject Create project methods
func createProject(project *models.Project) error {
	return app.DB.Create(project).Error
}

func updateProject(project *models.Project) error {
	return app.DB.Save(project).Error
}

func getProject(id int64) (*models.Project, error) {
	var project models.Project
	if err := app.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func getProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := app.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// createIteration Create iteration methods
func createIteration(iteration *models.Iteration) error {
	return app.DB.Create(iteration).Error
}

func updateIteration(iteration *models.Iteration) error {
	return app.DB.Save(iteration).Error
}

func getIteration(id int64) (*models.Iteration, error) {
	var iteration models.Iteration
	if err := app.DB.First(&iteration, id).Error; err != nil {
		return nil, err
	}
	return &iteration, nil
}

func getProjectIterations(projectID int64) ([]models.Iteration, error) {
	var iterations []models.Iteration
	if err := app.DB.Where("project_id = ?", projectID).Find(&iterations).Error; err != nil {
		return nil, err
	}
	return iterations, nil
}

// Issue methods
func createIssue(issue *models.Issue) error {
	if issue.ProjectID == 0 {
		return errors.New("project ID is required")
	}
	return app.DB.Create(issue).Error
}

func updateIssue(issue *models.Issue) error {
	return app.DB.Save(issue).Error
}

func getIssue(id int64) (*models.Issue, error) {
	var issue models.Issue
	if err := app.DB.First(&issue, id).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func getProjectIssues(projectID int64) ([]models.Issue, error) {
	var issues []models.Issue
	if err := app.DB.Where("project_id = ?", projectID).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

func getIterationIssues(iterationID int64) ([]models.Issue, error) {
	var issues []models.Issue
	if err := app.DB.Where("iteration_id = ?", iterationID).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

// createComment Create comment methods
func createComment(comment *models.Comment) error {
	if comment.IssueID == 0 {
		return errors.New("issue ID is required")
	}
	return app.DB.Create(comment).Error
}

func getIssueComments(issueID int64) ([]models.Comment, error) {
	var comments []models.Comment
	if err := app.DB.Where("issue_id = ?", issueID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
