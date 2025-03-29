package git

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var Srv *GitService

type gitModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&gitModule{}, enum.ModulePriorityGit)
}

func (m *gitModule) Awake(a *core.App) error {
	app = a

	// Initialize Git service
	Srv = initService()

	return autoMigrate()
}

func (m *gitModule) Start() error {
	return initData()
}

func (m *gitModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/git", app.HasPermission("git:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/git", fiber.Map{
			"Title": "Git管理",
			"Scripts": []string{
				"/static/js/admin/git.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterApi.Post("/git/clone", app.HasPermission("git:clone"), clone)
	app.RouterApi.Post("/git/pull", app.HasPermission("git:pull"), pull)
	app.RouterApi.Post("/git/push", app.HasPermission("git:push"), push)
	app.RouterApi.Get("/git/status", app.HasPermission("git:status"), status)
	app.RouterApi.Get("/git/log", app.HasPermission("git:log"), log)
	app.RouterApi.Post("/git/commit", app.HasPermission("git:commit"), commit)
	app.RouterApi.Post("/git/checkout", app.HasPermission("git:checkout"), checkout)
	app.RouterApi.Post("/git/branch", app.HasPermission("git:branch"), branch)
	app.RouterApi.Post("/git/merge", app.HasPermission("git:merge"), merge)
	app.RouterApi.Post("/git/reset", app.HasPermission("git:reset"), reset)
	app.RouterApi.Post("/git/stash", app.HasPermission("git:stash"), stash)

	return nil
}
