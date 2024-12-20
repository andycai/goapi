package git

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
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

	// Initialize Git service
	InitService()

	return nil
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
	app.RouterApi.Post("/git/clone", app.HasPermission("git:clone"), CloneHandler)
	app.RouterApi.Post("/git/pull", app.HasPermission("git:pull"), PullHandler)
	app.RouterApi.Post("/git/push", app.HasPermission("git:push"), PushHandler)
	app.RouterApi.Get("/git/status", app.HasPermission("git:status"), StatusHandler)
	app.RouterApi.Get("/git/log", app.HasPermission("git:log"), LogHandler)
	app.RouterApi.Post("/git/commit", app.HasPermission("git:commit"), CommitHandler)
	app.RouterApi.Post("/git/checkout", app.HasPermission("git:checkout"), CheckoutHandler)
	app.RouterApi.Post("/git/branch", app.HasPermission("git:branch"), BranchHandler)
	app.RouterApi.Post("/git/merge", app.HasPermission("git:merge"), MergeHandler)
	app.RouterApi.Post("/git/reset", app.HasPermission("git:reset"), ResetHandler)
	app.RouterApi.Post("/git/stash", app.HasPermission("git:stash"), StashHandler)

	return nil
}
