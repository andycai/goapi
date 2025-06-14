package parameter

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityParameter = 2001 // 数据中心-参数管理
)

var app *internal.App

type parameterModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&parameterModule{}, ModulePriorityParameter)
}

func (m *parameterModule) Awake(a *internal.App) error {
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
	app.RouterAdminApi.Get("/parameters", app.HasPermission("parameter:view"), listParameterHandler)
	app.RouterAdminApi.Get("/parameters/:id", app.HasPermission("parameter:view"), getParameterHandler)
	app.RouterAdminApi.Post("/parameters", app.HasPermission("parameter:create"), addParameterHandler)
	app.RouterAdminApi.Put("/parameters/:id", app.HasPermission("parameter:edit"), editParameterHandler)
	app.RouterAdminApi.Delete("/parameters/:id", app.HasPermission("parameter:delete"), deleteParameterHandler)

	return nil
}
