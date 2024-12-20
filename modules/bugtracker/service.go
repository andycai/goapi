package bugtracker

import (
	"errors"

	"github.com/andycai/unitool/models"
)

var bugtrackerService *BugtrackerService

type BugtrackerService struct{}

func InitService() {
	bugtrackerService = &BugtrackerService{}
	initTables()
}

func GetService() *BugtrackerService {
	return bugtrackerService
}

func initTables() {
	app.DB.AutoMigrate(&models.Project{}, &models.Iteration{}, &models.Issue{}, &models.Comment{})
}

// Project methods
func (s *BugtrackerService) CreateProject(project *models.Project) error {
	return app.DB.Create(project).Error
}

func (s *BugtrackerService) UpdateProject(project *models.Project) error {
	return app.DB.Save(project).Error
}

func (s *BugtrackerService) GetProject(id int64) (*models.Project, error) {
	var project models.Project
	if err := app.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *BugtrackerService) ListProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := app.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Iteration methods
func (s *BugtrackerService) CreateIteration(iteration *models.Iteration) error {
	return app.DB.Create(iteration).Error
}

func (s *BugtrackerService) UpdateIteration(iteration *models.Iteration) error {
	return app.DB.Save(iteration).Error
}

func (s *BugtrackerService) GetIteration(id int64) (*models.Iteration, error) {
	var iteration models.Iteration
	if err := app.DB.First(&iteration, id).Error; err != nil {
		return nil, err
	}
	return &iteration, nil
}

func (s *BugtrackerService) ListProjectIterations(projectID int64) ([]models.Iteration, error) {
	var iterations []models.Iteration
	if err := app.DB.Where("project_id = ?", projectID).Find(&iterations).Error; err != nil {
		return nil, err
	}
	return iterations, nil
}

// Issue methods
func (s *BugtrackerService) CreateIssue(issue *models.Issue) error {
	if issue.ProjectID == 0 {
		return errors.New("project ID is required")
	}
	return app.DB.Create(issue).Error
}

func (s *BugtrackerService) UpdateIssue(issue *models.Issue) error {
	return app.DB.Save(issue).Error
}

func (s *BugtrackerService) GetIssue(id int64) (*models.Issue, error) {
	var issue models.Issue
	if err := app.DB.First(&issue, id).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func (s *BugtrackerService) ListProjectIssues(projectID int64) ([]models.Issue, error) {
	var issues []models.Issue
	if err := app.DB.Where("project_id = ?", projectID).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *BugtrackerService) ListIterationIssues(iterationID int64) ([]models.Issue, error) {
	var issues []models.Issue
	if err := app.DB.Where("iteration_id = ?", iterationID).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

// Comment methods
func (s *BugtrackerService) CreateComment(comment *models.Comment) error {
	if comment.IssueID == 0 {
		return errors.New("issue ID is required")
	}
	return app.DB.Create(comment).Error
}

func (s *BugtrackerService) ListIssueComments(issueID int64) ([]models.Comment, error) {
	var comments []models.Comment
	if err := app.DB.Where("issue_id = ?", issueID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
