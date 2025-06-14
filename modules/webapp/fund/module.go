package fund

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityFund = 5100 // 功能-基金管理

var app *internal.App

type fundModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&fundModule{}, ModulePriorityFund)
}

func (m *fundModule) Awake(a *internal.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *fundModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *fundModule) AddPublicRouters() error {
	// 公开API
	app.RouterPublicApi.Get("/fund/market", getMarketDataHandler) // 获取市场指数
	app.RouterPublicApi.Get("/fund/hot", getHotFundsHandler)      // 获取热门基金
	return nil
}

func (m *fundModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/fund", app.HasPermission("fund:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/fund", fiber.Map{
			"Title": "基金管理",
			"Scripts": []string{
				"/static/js/admin/fund.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterAdminApi.Post("/fund/config", app.HasPermission("fund:config"), saveConfigHandler)
	app.RouterAdminApi.Get("/fund/config", app.HasPermission("fund:config"), getConfigHandler)
	app.RouterAdminApi.Post("/fund/sync", app.HasPermission("fund:sync"), syncDataHandler)
	app.RouterAdminApi.Get("/fund/list", app.HasPermission("fund:view"), listFundsHandler)
	app.RouterAdminApi.Get("/fund/market", app.HasPermission("fund:view"), getMarketDataHandler) // 保持原有路径

	return nil
}
