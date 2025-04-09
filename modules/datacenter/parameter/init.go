package parameter

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type parameterModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&parameterModule{}, enum.ModulePriorityParameter)
}

func (m *parameterModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *parameterModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *parameterModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/parameter", app.HasPermission("parameter:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/parameter", fiber.Map{
			"Title": "参数管理",
			"Scripts": []string{
				"/static/js/admin/parameter.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterApi.Get("/parameters", app.HasPermission("parameter:view"), listParametersHandler)
	app.RouterApi.Get("/parameters/:id", app.HasPermission("parameter:view"), getParameterHandler)
	app.RouterApi.Post("/parameters", app.HasPermission("parameter:create"), createParameterHandler)
	app.RouterApi.Put("/parameters/:id", app.HasPermission("parameter:edit"), updateParameterHandler)
	app.RouterApi.Delete("/parameters/:id", app.HasPermission("parameter:delete"), deleteParameterHandler)

	return nil
}
