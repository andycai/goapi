package bugtracker

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type bugtrackerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&bugtrackerModule{}, enum.ModulePriorityBugtracker)
}

func (m *bugtrackerModule) Awake(a *core.App) error {
	app = a

	// Initialize Bugtracker service
	InitService()

	return nil
}

func (m *bugtrackerModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/bugtracker", app.HasPermission("bugtracker:project:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/bugtracker", fiber.Map{
			"Title": "Bug Tracker",
			"Scripts": []string{
				"/static/js/admin/bugtracker.js",
			},
		}, "admin/layout")
	})

	// Project routes
	app.RouterApi.Get("/bugtracker/projects", app.HasPermission("bugtracker:project:list"), ListProjectsHandler)
	app.RouterApi.Post("/bugtracker/projects", app.HasPermission("bugtracker:project:create"), CreateProjectHandler)
	app.RouterApi.Put("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:update"), UpdateProjectHandler)
	app.RouterApi.Get("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:list"), GetProjectHandler)

	// Iteration routes
	app.RouterApi.Get("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:list"), ListIterationsHandler)
	app.RouterApi.Post("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:create"), CreateIterationHandler)
	app.RouterApi.Put("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:update"), UpdateIterationHandler)
	app.RouterApi.Get("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:list"), GetIterationHandler)

	// Issue routes
	app.RouterApi.Get("/bugtracker/issues", app.HasPermission("bugtracker:issue:list"), ListIssuesHandler)
	app.RouterApi.Post("/bugtracker/issues", app.HasPermission("bugtracker:issue:create"), CreateIssueHandler)
	app.RouterApi.Put("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:update"), UpdateIssueHandler)
	app.RouterApi.Get("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:list"), GetIssueHandler)

	// Comment routes
	app.RouterApi.Get("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:list"), ListCommentsHandler)
	app.RouterApi.Post("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:create"), CreateCommentHandler)

	return nil
}
