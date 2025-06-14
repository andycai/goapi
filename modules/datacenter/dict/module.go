package dict

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityDict = 2000 // 数据中心-字典管理
)

var app *internal.App

type dictModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&dictModule{}, ModulePriorityDict)
}

func (m *dictModule) Awake(a *internal.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *dictModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	return nil
}

func (m *dictModule) AddPublicRouters() error {
	// 公开API
	return nil
}

func (m *dictModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/dict", app.HasPermission("dict:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/dict", fiber.Map{
			"Title": "字典管理",
			"Scripts": []string{
				"/static/js/admin/dict.js",
			},
		}, "admin/layout")
	})

	// API路由
	dictTypeGroup := app.RouterAdminApi.Group("/dict/type")
	{
		dictTypeGroup.Get("/list", app.HasPermission("dict:view"), listDictTypeHandler)
		dictTypeGroup.Post("/add", app.HasPermission("dict:add"), addDictTypeHandler)
		dictTypeGroup.Post("/edit", app.HasPermission("dict:edit"), editDictTypeHandler)
		dictTypeGroup.Post("/delete", app.HasPermission("dict:delete"), deleteDictTypeHandler)
	}

	dictDataGroup := app.RouterAdminApi.Group("/dict/data")
	{
		dictDataGroup.Get("/list", app.HasPermission("dict:view"), listDictDataHandler)
		dictDataGroup.Post("/add", app.HasPermission("dict:add"), addDictDataHandler)
		dictDataGroup.Post("/edit", app.HasPermission("dict:edit"), editDictDataHandler)
		dictDataGroup.Post("/delete", app.HasPermission("dict:delete"), deleteDictDataHandler)
	}

	return nil
}
