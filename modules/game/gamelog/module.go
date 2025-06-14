package gamelog

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityGameLog = 9902 // 游戏-游戏日志

var app *internal.App

type gamelogModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&gamelogModule{}, ModulePriorityGameLog)
}

func (m *gamelogModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *gamelogModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *gamelogModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/gamelog", createLogHandler)
	return nil
}

func (m *gamelogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/gamelog", app.HasPermission("gamelog:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/gamelog", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/gamelog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/gamelog", app.HasPermission("gamelog:view"), listLogsHandler)
	app.RouterAdminApi.Delete("/gamelog/before", app.HasPermission("gamelog:delete"), deleteLogsBeforeHandler)
	app.RouterAdminApi.Delete("/gamelog/:id", app.HasPermission("gamelog:view"), deleteLogHandler)

	return nil
}
