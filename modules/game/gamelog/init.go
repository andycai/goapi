package gamelog

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type gamelogModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&gamelogModule{}, enum.ModulePriorityGameLog)
}

func (m *gamelogModule) Awake(a *core.App) error {
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
	app.RouterApi.Get("/gamelog", app.HasPermission("gamelog:view"), listLogsHandler)
	app.RouterApi.Delete("/gamelog/before", app.HasPermission("gamelog:delete"), deleteLogsBeforeHandler)
	app.RouterApi.Delete("/gamelog/:id", app.HasPermission("gamelog:view"), deleteLogHandler)

	return nil
}
