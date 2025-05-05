package unitool

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityUnitool = 9910 // 游戏-其他工具集合

var app *core.App

type unitoolModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&unitoolModule{}, ModulePriorityUnitool)
}

func (m *unitoolModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *unitoolModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *unitoolModule) AddPublicRouters() error {
	// 公开API
	app.RouterPublicApi.Post("/unitool/findguid", findGuidHandler)
	return nil
}

func (m *unitoolModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/unitool", app.HasPermission("unitool:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/unitool", fiber.Map{
			"Title": "Unity工具",
			"Scripts": []string{
				"/static/js/admin/unitool.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterAdminApi.Get("/unitool/logs", app.HasPermission("unitool:view"), getFindGuidLogsHandler)
	app.RouterAdminApi.Get("/unitool/duplicates/:id", app.HasPermission("unitool:view"), getDuplicateGuidsHandler)

	return nil
}
