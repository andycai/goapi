package adminlog

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityAdminLog = 1005 // 系统-管理员活动日志
)

var app *core.App

type adminlogModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&adminlogModule{}, ModulePriorityAdminLog)
}

func (m *adminlogModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *adminlogModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	return nil
}

func (m *adminlogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/adminlog", app.HasPermission("adminlog:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/adminlog", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/adminlog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/adminlog", app.HasPermission("adminlog:view"), listLogsHandler)
	app.RouterApi.Delete("/adminlog", app.HasPermission("adminlog:delete"), deleteLogsHandler)

	return nil
}
