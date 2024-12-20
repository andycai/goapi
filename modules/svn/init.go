package svn

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type svnModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&svnModule{}, enum.ModulePrioritySVN)
}

func (m *svnModule) Awake(a *core.App) error {
	app = a

	// Initialize SVN service
	InitService()

	return nil
}

func (m *svnModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/svn", app.HasPermission("svn:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/svn", fiber.Map{
			"Title": "SVN管理",
			"Scripts": []string{
				"/static/js/admin/svn.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterApi.Post("/svn/checkout", app.HasPermission("svn:checkout"), CheckoutHandler)
	app.RouterApi.Post("/svn/update", app.HasPermission("svn:update"), UpdateHandler)
	app.RouterApi.Post("/svn/commit", app.HasPermission("svn:commit"), CommitHandler)
	app.RouterApi.Get("/svn/status", app.HasPermission("svn:status"), StatusHandler)
	app.RouterApi.Get("/svn/info", app.HasPermission("svn:info"), InfoHandler)
	app.RouterApi.Get("/svn/log", app.HasPermission("svn:log"), LogHandler)
	app.RouterApi.Post("/svn/revert", app.HasPermission("svn:revert"), RevertHandler)
	app.RouterApi.Post("/svn/add", app.HasPermission("svn:add"), AddHandler)
	app.RouterApi.Delete("/svn/delete", app.HasPermission("svn:delete"), DeleteHandler)

	return nil
}
