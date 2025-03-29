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

	// Initialize Bugtracker service
	initService()

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
	app.RouterApi.Get("/bugtracker/projects", app.HasPermission("bugtracker:project:list"), listProjects)
	app.RouterApi.Post("/bugtracker/projects", app.HasPermission("bugtracker:project:create"), createProject)
	app.RouterApi.Put("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:update"), updateProject)
	app.RouterApi.Get("/bugtracker/projects/:id", app.HasPermission("bugtracker:project:list"), getProject)

	// Iteration routes
	app.RouterApi.Get("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:list"), listIterations)
	app.RouterApi.Post("/bugtracker/iterations", app.HasPermission("bugtracker:iteration:create"), createIteration)
	app.RouterApi.Put("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:update"), updateIteration)
	app.RouterApi.Get("/bugtracker/iterations/:id", app.HasPermission("bugtracker:iteration:list"), getIteration)

	// Issue routes
	app.RouterApi.Get("/bugtracker/issues", app.HasPermission("bugtracker:issue:list"), listIssues)
	app.RouterApi.Post("/bugtracker/issues", app.HasPermission("bugtracker:issue:create"), createIssue)
	app.RouterApi.Put("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:update"), updateIssue)
	app.RouterApi.Get("/bugtracker/issues/:id", app.HasPermission("bugtracker:issue:list"), getIssue)

	// Comment routes
	app.RouterApi.Get("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:list"), listComments)
	app.RouterApi.Post("/bugtracker/issues/:id/comments", app.HasPermission("bugtracker:comment:create"), createComment)

	return nil
}
