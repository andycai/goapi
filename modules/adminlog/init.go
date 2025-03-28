package adminlog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type adminlogModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&adminlogModule{}, enum.ModulePriorityAdminLog)
}

func (m *adminlogModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// 初始化数据
	return initData()
}

func (m *adminlogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/adminlog", app.HasPermission("adminlog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/adminlog", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/adminlog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/adminlog", app.HasPermission("adminlog:list"), getAdminLogs)
	app.RouterApi.Delete("/adminlog", app.HasPermission("adminlog:delete"), deleteAdminLogs)

	return nil
}
