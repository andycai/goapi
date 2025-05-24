package svn

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePrioritySVN = 9002 // 接口-SVN 接口

var app *core.App

type svnModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&svnModule{}, ModulePrioritySVN)
}

func (m *svnModule) Awake(a *core.App) error {
	app = a
	return autoMigrate()
}

func (m *svnModule) Start() error {
	if err := initData(); err != nil {
		return err
	}

	// Initialize SVN service
	initService()
	return nil
}

func (m *svnModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/svn", app.HasPermission("svn:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/svn", fiber.Map{
			"Title": "SVN管理",
			"Scripts": []string{
				"/static/js/admin/svn.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterAdminApi.Post("/svn/checkout", app.HasPermission("svn:checkout"), checkoutHandler)
	app.RouterAdminApi.Post("/svn/update", app.HasPermission("svn:update"), updateHandler)
	app.RouterAdminApi.Post("/svn/commit", app.HasPermission("svn:commit"), commitHandler)
	app.RouterAdminApi.Get("/svn/status", app.HasPermission("svn:status"), statusHandler)
	app.RouterAdminApi.Get("/svn/info", app.HasPermission("svn:info"), infoHandler)
	app.RouterAdminApi.Get("/svn/log", app.HasPermission("svn:log"), logHandler)
	app.RouterAdminApi.Post("/svn/revert", app.HasPermission("svn:revert"), revertHandler)
	app.RouterAdminApi.Post("/svn/add", app.HasPermission("svn:add"), addHandler)
	app.RouterAdminApi.Delete("/svn/delete", app.HasPermission("svn:delete"), deleteHandler)

	return nil
}
