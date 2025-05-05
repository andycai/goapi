package bugtracker

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityBugTracker = 4003 // 功能-Bug 跟踪

var app *core.App

type bugtrackerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&bugtrackerModule{}, ModulePriorityBugTracker)
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
	app.RouterAdminApi.Get("/bugtracker/projects", app.HasPermission("bugtracker:project:view"), listProjectsHandler)
	app.RouterAdminApi.Post("/bugtracker/projects", app.HasPermission("bugtracker:project:create"), createProjectHandler)
	app.RouterAdminApi.Put("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:update"), updateProjectHandler)
	app.RouterAdminApi.Get("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:view"), getProjectHandler)

	// Iteration routes
	app.RouterAdminApi.Get("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:view"), listIterationsHandler)
	app.RouterAdminApi.Post("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:create"), createIterationHandler)
	app.RouterAdminApi.Put("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:update"), updateIterationHandler)
	app.RouterAdminApi.Get("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:view"), getIterationHandler)

	// Issue routes
	app.RouterAdminApi.Get("/bugtracker/issues", app.HasPermission("bugtracker:issue:view"), listIssuesHandler)
	app.RouterAdminApi.Post("/bugtracker/issues", app.HasPermission("bugtracker:issue:create"), createIssueHandler)
	app.RouterAdminApi.Put("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:update"), updateIssueHandler)
	app.RouterAdminApi.Get("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:view"), getIssueHandler)

	// Comment routes
	app.RouterAdminApi.Get("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:view"), listCommentsHandler)
	app.RouterAdminApi.Post("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:create"), createCommentHandler)

	return nil
}
