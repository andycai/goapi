package menu

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var dao *MenuDao

type menuModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&menuModule{}, enum.ModulePriorityMenu)
}

func (m *menuModule) Awake(a *core.App) error {
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
	app.RouterApi.Get("/menus", app.HasPermission("menu:view"), listMenusHandler)
	app.RouterApi.Get("/menus/tree", app.HasPermission("menu:view"), getMenuTreeHandler)
	app.RouterApi.Post("/menus", app.HasPermission("menu:create"), createMenuHandler)
	app.RouterApi.Put("/menus/:id", app.HasPermission("menu:update"), updateMenuHandler)
	app.RouterApi.Delete("/menus/:id", app.HasPermission("menu:delete"), deleteMenuHandler)

	return nil
}
