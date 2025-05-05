package channel

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityChannel = 9907 // 游戏-渠道管理

var app *core.App

type channelModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&channelModule{}, ModulePriorityChannel)
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
		return c.Render("admin/physical_server", fiber.Map{
			"Title": "服务器管理",
			"Scripts": []string{
				"/static/js/admin/physical_server.js",
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
	app.RouterAdminApi.Get("/channel/list", app.HasPermission("channel:view"), getChannelsHandler)
	app.RouterAdminApi.Post("/channel", app.HasPermission("channel:manage"), createChannelHandler)
	app.RouterAdminApi.Put("/channel/:id", app.HasPermission("channel:manage"), updateChannelHandler)
	app.RouterAdminApi.Delete("/channel/:id", app.HasPermission("channel:manage"), deleteChannelHandler)

	// 物理服务器相关
	app.RouterAdminApi.Get("/physical_servers", app.HasPermission("server:view"), getPhysicalServersHandler)
	app.RouterAdminApi.Post("/physical_servers", app.HasPermission("server:manage"), createPhysicalServerHandler)
	app.RouterAdminApi.Put("/physical_servers/:id", app.HasPermission("server:manage"), updatePhysicalServerHandler)
	app.RouterAdminApi.Delete("/physical_servers/:id", app.HasPermission("server:manage"), deletePhysicalServerHandler)

	// 服务器分组相关
	app.RouterAdminApi.Get("/server/group/list", app.HasPermission("server:view"), getServerGroupsHandler)
	app.RouterAdminApi.Post("/server/group", app.HasPermission("server:manage"), createServerGroupHandler)
	app.RouterAdminApi.Put("/server/group/:id", app.HasPermission("server:manage"), updateServerGroupHandler)
	app.RouterAdminApi.Delete("/server/group/:id", app.HasPermission("server:manage"), deleteServerGroupHandler)

	// 公告相关
	app.RouterAdminApi.Get("/announcement/list", app.HasPermission("announcement:view"), getAnnouncementsHandler)
	app.RouterAdminApi.Post("/announcement", app.HasPermission("announcement:manage"), createAnnouncementHandler)
	app.RouterAdminApi.Put("/announcement/:id", app.HasPermission("announcement:manage"), updateAnnouncementHandler)
	app.RouterAdminApi.Delete("/announcement/:id", app.HasPermission("announcement:manage"), deleteAnnouncementHandler)

	return nil
}
