package git

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityGit = 9003 // 接口-Git 接口

var app *core.App

type gitModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&gitModule{}, ModulePriorityGit)
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
	app.RouterAdminApi.Post("/git/clone", app.HasPermission("git:clone"), cloneHandler)
	app.RouterAdminApi.Post("/git/pull", app.HasPermission("git:pull"), pullHandler)
	app.RouterAdminApi.Post("/git/push", app.HasPermission("git:push"), pushHandler)
	app.RouterAdminApi.Get("/git/status", app.HasPermission("git:status"), statusHandler)
	app.RouterAdminApi.Get("/git/log", app.HasPermission("git:log"), logHandler)
	app.RouterAdminApi.Post("/git/commit", app.HasPermission("git:commit"), commitHandler)
	app.RouterAdminApi.Post("/git/checkout", app.HasPermission("git:checkout"), checkoutHandler)
	app.RouterAdminApi.Post("/git/branch", app.HasPermission("git:branch"), branchHandler)
	app.RouterAdminApi.Post("/git/merge", app.HasPermission("git:merge"), mergeHandler)
	app.RouterAdminApi.Post("/git/reset", app.HasPermission("git:reset"), resetHandler)
	app.RouterAdminApi.Post("/git/stash", app.HasPermission("git:stash"), stashHandler)

	return nil
}
