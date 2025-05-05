package role

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityRole = 1001 // 系统-角色管理
)

var app *core.App

type roleModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&roleModule{}, ModulePriorityRole)
}

func (m *roleModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *roleModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *roleModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/roles", app.HasPermission("role:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/roles", fiber.Map{
			"Title": "角色管理",
			"Scripts": []string{
				"/static/js/admin/roles.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/roles", app.HasPermission("role:view"), listRolesHandler)
	app.RouterAdminApi.Post("/roles", app.HasPermission("role:create"), createRoleHandler)
	app.RouterAdminApi.Put("/roles/:id", app.HasPermission("role:update"), updateRoleHandler)
	app.RouterAdminApi.Delete("/roles/:id", app.HasPermission("role:delete"), deleteRoleHandler)

	return nil
}
