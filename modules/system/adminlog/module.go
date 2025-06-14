package adminlog

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityAdminLog = 1005 // 系统-管理员活动日志
)

var app *internal.App

type adminlogModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&adminlogModule{}, ModulePriorityAdminLog)
}

func (m *adminlogModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *adminlogModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	subscribeEvents(app.Bus)

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
	app.RouterAdminApi.Get("/adminlog", app.HasPermission("adminlog:view"), listLogsHandler)
	app.RouterAdminApi.Delete("/adminlog", app.HasPermission("adminlog:delete"), deleteLogsHandler)

	return nil
}
