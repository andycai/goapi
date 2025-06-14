package user

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityUser = 1000 // 系统-用户管理
)

var app *internal.App

type userModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&userModule{}, ModulePriorityUser)
}

func (m *userModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *userModule) Start() error {
	return initData()
}

func (m *userModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/users", app.HasPermission("user:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/users", fiber.Map{
			"Title": "用户管理",
			"Scripts": []string{
				"/static/js/admin/users.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/users", app.HasPermission("user:view"), listUsersHandler)
	app.RouterAdminApi.Post("/users", app.HasPermission("user:create"), createUserHandler)
	app.RouterAdminApi.Put("/users/:id", app.HasPermission("user:update"), updateUserHandler)
	app.RouterAdminApi.Delete("/users/:id", app.HasPermission("user:delete"), deleteUserHandler)

	return nil
}
