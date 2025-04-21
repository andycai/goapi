package command

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type commandModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&commandModule{}, enum.ModulePriorityCommand)
}

func (m *commandModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *commandModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *commandModule) AddPublicRouters() error {
	// 公开API
	return nil
}

func (m *commandModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/command", app.HasPermission("command:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/command", fiber.Map{
			"Title": "命令管理",
			"Scripts": []string{
				"/static/js/admin/command.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterApi.Get("/command/list", app.HasPermission("command:view"), getCommandsHandler)
	app.RouterApi.Post("/command", app.HasPermission("command:manage"), createCommandHandler)
	app.RouterApi.Post("/command/:id/execute", app.HasPermission("command:manage"), executeCommandHandler)
	app.RouterApi.Get("/command/:id/executions", app.HasPermission("command:view"), getCommandExecutionsHandler)

	return nil
}
