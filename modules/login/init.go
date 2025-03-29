package login

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type loginModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&loginModule{}, enum.ModulePriorityLogin)
}

func (m *loginModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *loginModule) AddPublicRouters() error {
	// 登录页面路由（不需要认证）
	app.RouterPublic.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}, "login")
	})

	// 登录 API 路由（不需要认证）
	app.RouterPublic.Post("/login", loginHandler)
	// 退出登录 API 路由（不需要认证）
	app.RouterPublic.Get("/logout", logoutHandler)
	app.RouterPublic.Post("/logout", logoutHandler)

	// 修改密码路由（不需要认证）
	app.RouterPublic.Post("/change-password", changePasswordHandler)

	return nil
}

func (m *loginModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/", func(c *fiber.Ctx) error {
		return c.Render("admin/index", fiber.Map{
			"Title": "管理后台",
		}, "admin/layout")
	})

	// api

	return nil
}
