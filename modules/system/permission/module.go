package permission

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityPermission = 1002 // 系统-权限管理
)

var app *internal.App

type permissionModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&permissionModule{}, ModulePriorityPermission)
}

func (m *permissionModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *permissionModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *permissionModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/permissions", app.HasPermission("permission:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/permissions", fiber.Map{
			"Title": "权限管理",
			"Scripts": []string{
				"/static/js/admin/permissions.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/permissions", app.HasPermission("permission:view"), listPermissionsHandler)
	app.RouterAdminApi.Get("/permissions/all", app.HasPermission("permission:view"), listAllPermissionsHandler)
	app.RouterAdminApi.Post("/permissions", app.HasPermission("permission:create"), createPermissionHandler)
	app.RouterAdminApi.Put("/permissions/:id", app.HasPermission("permission:update"), updatePermissionHandler)
	app.RouterAdminApi.Delete("/permissions/:id", app.HasPermission("permission:delete"), deletePermissionHandler)

	return nil
}
