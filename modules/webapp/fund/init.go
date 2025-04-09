package fund

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type fundModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&fundModule{}, enum.ModulePriorityFund)
}

func (m *fundModule) Awake(a *core.App) error {
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
	app.RouterApi.Post("/fund/config", app.HasPermission("fund:config"), saveConfigHandler)
	app.RouterApi.Get("/fund/config", app.HasPermission("fund:config"), getConfigHandler)
	app.RouterApi.Post("/fund/sync", app.HasPermission("fund:sync"), syncDataHandler)
	app.RouterApi.Get("/fund/list", app.HasPermission("fund:view"), listFundsHandler)
	app.RouterApi.Get("/fund/market", app.HasPermission("fund:view"), getMarketDataHandler) // 保持原有路径

	return nil
}
