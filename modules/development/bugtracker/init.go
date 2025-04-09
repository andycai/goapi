package bugtracker

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type bugtrackerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&bugtrackerModule{}, enum.ModulePriorityBugTracker)
}

func (m *bugtrackerModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *bugtrackerModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	return nil
}

func (m *bugtrackerModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/bugtracker", app.HasPermission("bugtracker:project:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/bugtracker", fiber.Map{
			"Title": "Bug Tracker",
			"Scripts": []string{
				"/static/js/admin/bugtracker.js",
			},
		}, "admin/layout")
	})

	// Project routes
	app.RouterApi.Get("/bugtracker/projects", app.HasPermission("bugtracker:project:view"), listProjectsHandler)
	app.RouterApi.Post("/bugtracker/projects", app.HasPermission("bugtracker:project:create"), createProjectHandler)
	app.RouterApi.Put("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:update"), updateProjectHandler)
	app.RouterApi.Get("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:view"), getProjectHandler)

	// Iteration routes
	app.RouterApi.Get("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:view"), listIterationsHandler)
	app.RouterApi.Post("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:create"), createIterationHandler)
	app.RouterApi.Put("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:update"), updateIterationHandler)
	app.RouterApi.Get("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:view"), getIterationHandler)

	// Issue routes
	app.RouterApi.Get("/bugtracker/issues", app.HasPermission("bugtracker:issue:view"), listIssuesHandler)
	app.RouterApi.Post("/bugtracker/issues", app.HasPermission("bugtracker:issue:create"), createIssueHandler)
	app.RouterApi.Put("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:update"), updateIssueHandler)
	app.RouterApi.Get("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:view"), getIssueHandler)

	// Comment routes
	app.RouterApi.Get("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:view"), listCommentsHandler)
	app.RouterApi.Post("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:create"), createCommentHandler)

	return nil
}
