package git

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type gitModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&gitModule{}, enum.ModulePriorityGit)
}

func (m *gitModule) Awake(a *core.App) error {
	app = a
	return autoMigrate()
}

func (m *gitModule) Start() error {
	if err := initData(); err != nil {
		return err
	}

	// Initialize Git service
	initService()
	return nil
}

func (m *gitModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/git", app.HasPermission("git:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/git", fiber.Map{
			"Title": "Git管理",
			"Scripts": []string{
				"/static/js/admin/git.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterApi.Post("/git/clone", app.HasPermission("git:clone"), cloneHandler)
	app.RouterApi.Post("/git/pull", app.HasPermission("git:pull"), pullHandler)
	app.RouterApi.Post("/git/push", app.HasPermission("git:push"), pushHandler)
	app.RouterApi.Get("/git/status", app.HasPermission("git:status"), statusHandler)
	app.RouterApi.Get("/git/log", app.HasPermission("git:log"), logHandler)
	app.RouterApi.Post("/git/commit", app.HasPermission("git:commit"), commitHandler)
	app.RouterApi.Post("/git/checkout", app.HasPermission("git:checkout"), checkoutHandler)
	app.RouterApi.Post("/git/branch", app.HasPermission("git:branch"), branchHandler)
	app.RouterApi.Post("/git/merge", app.HasPermission("git:merge"), mergeHandler)
	app.RouterApi.Post("/git/reset", app.HasPermission("git:reset"), resetHandler)
	app.RouterApi.Post("/git/stash", app.HasPermission("git:stash"), stashHandler)

	return nil
}
