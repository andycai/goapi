package menu

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityMenu = 1004 // 系统-菜单管理
)

var app *internal.App
var dao *MenuDao

type menuModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&menuModule{}, ModulePriorityMenu)
}

func (m *menuModule) Awake(a *internal.App) error {
	app = a
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *menuModule) Start() error {
	dao = NewMenuDao()

	// 初始化数据
	return initData()
}

func (m *menuModule) AddPublicRouters() error {
	app.RouterPublicApi.Get("/menus/public/tree", getMenuTreeHandler)
	return nil
}

func (m *menuModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/menus", app.HasPermission("menu:view"), func(c *fiber.Ctx) error {
		user := app.CurrentUser(c)

		return c.Render("admin/menus", fiber.Map{
			"Title": "菜单管理",
			"Scripts": []string{
				"/static/js/admin/menus.js",
			},
			"user": user,
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/menus", app.HasPermission("menu:view"), listMenusHandler)
	app.RouterAdminApi.Get("/menus/tree", app.HasPermission("menu:view"), getMenuTreeHandler)
	app.RouterAdminApi.Post("/menus", app.HasPermission("menu:create"), createMenuHandler)
	app.RouterAdminApi.Put("/menus/:id", app.HasPermission("menu:update"), updateMenuHandler)
	app.RouterAdminApi.Delete("/menus/:id", app.HasPermission("menu:delete"), deleteMenuHandler)

	return nil
}
