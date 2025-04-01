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

func (m *reposyncModule) AddPublicRouters() error {
	// 公开API
	app.RouterPublicApi.Post("/reposync/rangesync", syncPublicCommitsHandler)
	app.RouterPublicApi.Post("/reposync/autosync", syncPublicAutoHandler)
	return nil
}

func (m *reposyncModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/reposync", app.HasPermission("reposync:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/reposync", fiber.Map{
			"Title": "仓库同步",
			"Scripts": []string{
				"/static/js/admin/reposync.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterApi.Post("/reposync/config", app.HasPermission("reposync:config"), saveConfigHandler)
	app.RouterApi.Get("/reposync/config", app.HasPermission("reposync:config"), getConfigHandler)
	app.RouterApi.Post("/reposync/checkout", app.HasPermission("reposync:checkout"), checkoutHandler)
	app.RouterApi.Get("/reposync/commits", app.HasPermission("reposync:view"), listCommitsHandler)
	app.RouterApi.Post("/reposync/sync", app.HasPermission("reposync:sync"), syncCommitsHandler)
	app.RouterApi.Post("/reposync/refresh", app.HasPermission("reposync:view"), refreshCommitsHandler)
	app.RouterApi.Post("/reposync/clear", app.HasPermission("reposync:config"), clearSyncDataHandler)

	return nil
}
