package reposync

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityRepoSync = 9906 // 游戏-仓库同步

var app *internal.App

type reposyncModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&reposyncModule{}, ModulePriorityRepoSync)
}

func (m *reposyncModule) Awake(a *internal.App) error {
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
	app.RouterAdminApi.Post("/reposync/config", app.HasPermission("reposync:config"), saveConfigHandler)
	app.RouterAdminApi.Get("/reposync/config", app.HasPermission("reposync:config"), getConfigHandler)
	app.RouterAdminApi.Post("/reposync/checkout", app.HasPermission("reposync:checkout"), checkoutHandler)
	app.RouterAdminApi.Get("/reposync/commits", app.HasPermission("reposync:view"), listCommitsHandler)
	app.RouterAdminApi.Post("/reposync/sync", app.HasPermission("reposync:sync"), syncCommitsHandler)
	app.RouterAdminApi.Post("/reposync/refresh", app.HasPermission("reposync:view"), refreshCommitsHandler)
	app.RouterAdminApi.Post("/reposync/clear", app.HasPermission("reposync:config"), clearSyncDataHandler)

	return nil
}
