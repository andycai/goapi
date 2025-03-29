package reposync

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type reposyncModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&reposyncModule{}, enum.ModulePriorityRepoSync)
}

func (m *reposyncModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *reposyncModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *reposyncModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/reposync", app.HasPermission("reposync:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/reposync", fiber.Map{
			"Title": "仓库同步",
			"Scripts": []string{
				"/static/js/admin/reposync.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterApi.Post("/reposync/config", app.HasPermission("reposync:config"), saveConfig)
	app.RouterApi.Get("/reposync/config", app.HasPermission("reposync:config"), getConfig)
	app.RouterApi.Post("/reposync/checkout", app.HasPermission("reposync:checkout"), checkout)
	app.RouterApi.Get("/reposync/commits", app.HasPermission("reposync:list"), getCommits)
	app.RouterApi.Post("/reposync/sync", app.HasPermission("reposync:sync"), syncCommits)

	return nil
}
