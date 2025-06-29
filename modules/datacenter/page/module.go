package page

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityPage = 2002 // 数据中心-页面管理
)

var app *internal.App

type pageModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&pageModule{}, ModulePriorityPage)
}

func (m *pageModule) Awake(a *internal.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *pageModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	return nil
}

func (m *pageModule) AddPublicRouters() error {
	// 公开API
	app.RouterPublic.Get("/pages", listPageHandler)
	app.RouterPublic.Get("/pages/:slug", getPageBySlugHandler)
	app.RouterPublic.Get("/pages/search", searchPageHandler)

	return nil
}

func (m *pageModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/page", app.HasPermission("page:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/page", fiber.Map{
			"Title": "页面管理",
			"Scripts": []string{
				"/static/js/admin/page.js",
			},
		}, "admin/layout")
	})

	// API路由
	pageGroup := app.RouterAdminApi.Group("/page")
	{
		pageGroup.Get("/list", app.HasPermission("page:view"), listPageHandler)
		pageGroup.Post("/add", app.HasPermission("page:add"), addPageHandler)
		pageGroup.Post("/edit", app.HasPermission("page:edit"), editPageHandler)
		pageGroup.Post("/delete", app.HasPermission("page:delete"), deletePageHandler)
	}

	return nil
}
