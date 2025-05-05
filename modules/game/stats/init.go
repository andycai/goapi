package stats

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityStats = 9903 // 游戏-游戏统计

var app *core.App

type statsModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&statsModule{}, ModulePriorityStats)
}

func (m *statsModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *statsModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *statsModule) AddPublicRouters() error {
	app.RouterPublicApi.Post("/stats", CreateStats)

	return nil
}

func (m *statsModule) AddAuthRouters() error {

	app.RouterAdminApi.Get("/stats", app.HasPermission("stats:view"), listStatsHandler)
	app.RouterAdminApi.Delete("/stats/before", app.HasPermission("stats:delete"), deleteStatsBeforeHandler)
	app.RouterAdminApi.Get("/stats/details", app.HasPermission("stats:view"), getStatDetailsHandler)
	app.RouterAdminApi.Delete("/stats/:id", app.HasPermission("stats:delete"), deleteStatHandler)

	// admin
	app.RouterAdmin.Get("/stats", app.HasPermission("stats:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/stats", fiber.Map{
			"Title": "游戏统计",
			"Scripts": []string{
				"/static/js/chart-4.4.4.js",
				"/static/js/hammer-2.0.8.js",
				"/static/js/chartjs-plugin-zoom.min.js",
				"/static/js/chartjs-adapter-date-fns.bundle.min.js",
				"/static/js/admin/stats.js",
			},
		}, "admin/layout")
	})

	return nil
}
