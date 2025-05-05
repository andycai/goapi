package serverconf

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityServerConf = 9904 // 游戏-服务器配置

var app *core.App

type serverconfModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&serverconfModule{}, ModulePriorityServerConf)
}

func (m *serverconfModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *serverconfModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *serverconfModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Get("/game/serverlist", getServerList)
	app.RouterPublicApi.Get("/game/lastserver", getLastServer)
	app.RouterPublicApi.Get("/game/serverinfo", getServerInfo)
	app.RouterPublicApi.Get("/game/noticelist", getNoticeList)
	app.RouterPublicApi.Get("/game/noticenum", getNoticeNum)
	app.RouterPublicApi.Get("/serverlist", getServerList)
	app.RouterPublicApi.Get("/lastserver", getLastServer)
	app.RouterPublicApi.Get("/serverinfo", getServerInfo)
	app.RouterPublicApi.Get("/noticelist", getNoticeList)
	app.RouterPublicApi.Get("/noticenum", getNoticeNum)

	return nil
}

func (m *serverconfModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/serverconf", app.HasPermission("serverconf:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/serverconf", fiber.Map{
			"Title": "服务器配置",
			"Scripts": []string{
				"/static/js/admin/serverconf.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Post("/game/serverlist", app.HasPermission("serverconf:update"), updateServerListHandler)
	app.RouterAdminApi.Post("/game/lastserver", app.HasPermission("serverconf:update"), updateLastServerHandler)
	app.RouterAdminApi.Post("/game/serverinfo", app.HasPermission("serverconf:update"), updateServerInfoHandler)
	app.RouterAdminApi.Post("/game/noticelist", app.HasPermission("serverconf:update"), updateNoticeListHandler)
	app.RouterAdminApi.Post("/game/noticenum", app.HasPermission("serverconf:update"), updateNoticeNumHandler)

	return nil
}
