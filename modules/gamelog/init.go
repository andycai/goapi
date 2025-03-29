package gamelog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
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
	app.RouterPublicApi.Post("/gamelog", createLog)
	return nil
}

func (m *gamelogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/gamelog", app.HasPermission("gamelog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/gamelog", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/gamelog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/gamelog", app.HasPermission("gamelog:list"), getLogs)
	app.RouterApi.Delete("/gamelog/before", app.HasPermission("gamelog:delete"), deleteLogsBefore)
	app.RouterApi.Delete("/gamelog/:id", app.HasPermission("gamelog:list"), deleteLog)

	return nil
}
