package patch

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityPatch = 9905 // 游戏-热更新

var app *core.App

type patchModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&patchModule{}, ModulePriorityPatch)
}

func (m *patchModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *patchModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *patchModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/patch", app.HasPermission("patch:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/patch", fiber.Map{
			"Title": "补丁管理",
			"Scripts": []string{
				"/static/js/admin/patch.js",
			},
		}, "admin/layout")
	})

	// API路由
	app.RouterAdminApi.Post("/patch/config", app.HasPermission("patch:config"), saveConfigHandler)
	app.RouterAdminApi.Get("/patch/config", app.HasPermission("patch:config"), getConfigHandler)
	app.RouterAdminApi.Post("/patch/generate", app.HasPermission("patch:generate"), generatePatchHandler)
	app.RouterAdminApi.Get("/patch/records", app.HasPermission("patch:view"), listPatchRecordsHandler)
	app.RouterAdminApi.Post("/patch/apply", app.HasPermission("patch:apply"), applyPatchHandler)

	return nil
}
