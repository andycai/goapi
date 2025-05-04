package channel

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type channelModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&channelModule{}, enum.ModulePriorityChannel)
}

func (m *channelModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *channelModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// 初始化服务
	initService()

	return nil
}

func (m *channelModule) AddPublicRouters() error {
	// 公开API
	return nil
}

func (m *channelModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/channel", app.HasPermission("channel:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/channel", fiber.Map{
			"Title": "渠道管理",
			"Scripts": []string{
				"/static/js/admin/channel.js",
			},
		}, "admin/layout")
	})

	app.RouterAdmin.Get("/server", app.HasPermission("server:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/server", fiber.Map{
			"Title": "服务器管理",
			"Scripts": []string{
				"/static/js/admin/server.js",
			},
		}, "admin/layout")
	})

	app.RouterAdmin.Get("/announcement", app.HasPermission("announcement:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/announcement", fiber.Map{
			"Title": "公告管理",
			"Scripts": []string{
				"/static/js/admin/announcement.js",
			},
		}, "admin/layout")
	})

	// API路由
	// 渠道相关
	app.RouterApi.Get("/channel/list", app.HasPermission("channel:view"), getChannelsHandler)
	app.RouterApi.Post("/channel", app.HasPermission("channel:manage"), createChannelHandler)
	app.RouterApi.Put("/channel/:id", app.HasPermission("channel:manage"), updateChannelHandler)
	app.RouterApi.Delete("/channel/:id", app.HasPermission("channel:manage"), deleteChannelHandler)

	// 物理服务器相关
	app.RouterApi.Get("/server/list", app.HasPermission("server:view"), getPhysicalServersHandler)
	app.RouterApi.Post("/server", app.HasPermission("server:manage"), createPhysicalServerHandler)
	app.RouterApi.Put("/server/:id", app.HasPermission("server:manage"), updatePhysicalServerHandler)
	app.RouterApi.Delete("/server/:id", app.HasPermission("server:manage"), deletePhysicalServerHandler)

	// 服务器分组相关
	app.RouterApi.Get("/server/group/list", app.HasPermission("server:view"), getServerGroupsHandler)
	app.RouterApi.Post("/server/group", app.HasPermission("server:manage"), createServerGroupHandler)
	app.RouterApi.Put("/server/group/:id", app.HasPermission("server:manage"), updateServerGroupHandler)
	app.RouterApi.Delete("/server/group/:id", app.HasPermission("server:manage"), deleteServerGroupHandler)

	// 公告相关
	app.RouterApi.Get("/announcement/list", app.HasPermission("announcement:view"), getAnnouncementsHandler)
	app.RouterApi.Post("/announcement", app.HasPermission("announcement:manage"), createAnnouncementHandler)
	app.RouterApi.Put("/announcement/:id", app.HasPermission("announcement:manage"), updateAnnouncementHandler)
	app.RouterApi.Delete("/announcement/:id", app.HasPermission("announcement:manage"), deleteAnnouncementHandler)

	return nil
}
