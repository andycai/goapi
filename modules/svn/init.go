package svn

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var Srv *SVNService

type svnModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&svnModule{}, enum.ModulePrioritySVN)
}

func (m *svnModule) Awake(a *core.App) error {
	app = a

	// Initialize SVN service
	Srv = initService()

	return autoMigrate()
}

func (m *svnModule) Start() error {
	return initData()
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
	app.RouterApi.Post("/svn/checkout", app.HasPermission("svn:checkout"), checkout)
	app.RouterApi.Post("/svn/update", app.HasPermission("svn:update"), update)
	app.RouterApi.Post("/svn/commit", app.HasPermission("svn:commit"), commit)
	app.RouterApi.Get("/svn/status", app.HasPermission("svn:status"), status)
	app.RouterApi.Get("/svn/info", app.HasPermission("svn:info"), info)
	app.RouterApi.Get("/svn/log", app.HasPermission("svn:log"), log)
	app.RouterApi.Post("/svn/revert", app.HasPermission("svn:revert"), revert)
	app.RouterApi.Post("/svn/add", app.HasPermission("svn:add"), add)
	app.RouterApi.Delete("/svn/delete", app.HasPermission("svn:delete"), delete)

	return nil
}
