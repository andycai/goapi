package bases

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityBases = 2002 // 数据中心-基础数据
)

var app *internal.App

type basesModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&basesModule{}, ModulePriorityBases)
}

func (m *basesModule) Awake(a *internal.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *basesModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *basesModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/bases", app.HasPermission("bases:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/bases", fiber.Map{
			"Title":   "基础数据",
			"Scripts": []string{"/static/js/admin/bases.js"},
		}, "admin/layout")
	})

	// API路由
	app.RouterAdminApi.Get("/bases/entities", app.HasPermission("bases:view"), listEntityHandler)
	app.RouterAdminApi.Get("/bases/entities/:id", app.HasPermission("bases:view"), getEntityHandler)
	app.RouterAdminApi.Post("/bases/entities", app.HasPermission("bases:create"), addEntityHandler)
	app.RouterAdminApi.Put("/bases/entities/:id", app.HasPermission("bases:edit"), editEntityHandler)
	app.RouterAdminApi.Delete("/bases/entities/:id", app.HasPermission("bases:delete"), deleteEntityHandler)

	app.RouterAdminApi.Get("/bases/fields", app.HasPermission("bases:view"), listFieldHandler)
	app.RouterAdminApi.Get("/bases/fields/:id", app.HasPermission("bases:view"), getFieldHandler)
	app.RouterAdminApi.Post("/bases/fields", app.HasPermission("bases:create"), addFieldHandler)
	app.RouterAdminApi.Put("/bases/fields/:id", app.HasPermission("bases:edit"), editFieldHandler)
	app.RouterAdminApi.Delete("/bases/fields/:id", app.HasPermission("bases:delete"), deleteFieldHandler)

	app.RouterAdminApi.Get("/bases/data", app.HasPermission("bases:view"), listEntityDataHandler)
	app.RouterAdminApi.Get("/bases/data/:id", app.HasPermission("bases:view"), getEntityDataHandler)
	app.RouterAdminApi.Post("/bases/data", app.HasPermission("bases:create"), addEntityDataHandler)
	app.RouterAdminApi.Put("/bases/data/:id", app.HasPermission("bases:edit"), editEntityDataHandler)
	app.RouterAdminApi.Delete("/bases/data/:id", app.HasPermission("bases:delete"), deleteEntityDataHandler)

	return nil
}
