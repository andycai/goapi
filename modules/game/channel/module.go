package channel

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityChannel = 9907 // 游戏-渠道管理

var app *internal.App

type channelModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&channelModule{}, ModulePriorityChannel)
}

func (m *channelModule) Awake(a *internal.App) error {
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

	app.RouterAdmin.Get("/physical_servers", app.HasPermission("server:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/physical_server", fiber.Map{
			"Title": "物理服务器管理",
			"Scripts": []string{
				"/static/js/admin/physical_server.js",
			},
		}, "admin/layout")
	})

	app.RouterAdmin.Get("/server_groups", app.HasPermission("server:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/server_group", fiber.Map{
			"Title": "服务器分组管理",
			"Scripts": []string{
				"/static/js/admin/server_group.js",
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
	app.RouterAdminApi.Get("/server_groups", app.HasPermission("server:view"), getServerGroupsHandler)
	app.RouterAdminApi.Post("/server_groups", app.HasPermission("server:manage"), createServerGroupHandler)
	app.RouterAdminApi.Put("/server_groups/:id", app.HasPermission("server:manage"), updateServerGroupHandler)
	app.RouterAdminApi.Delete("/server_groups/:id", app.HasPermission("server:manage"), deleteServerGroupHandler)
	app.RouterAdminApi.Get("/server_groups/:id/servers", app.HasPermission("server:view"), getServerGroupServersHandler)
	app.RouterAdminApi.Post("/server_groups/:id/servers/:serverId", app.HasPermission("server:manage"), addServerToGroupHandler)
	app.RouterAdminApi.Delete("/server_groups/:id/servers/:serverId", app.HasPermission("server:manage"), removeServerFromGroupHandler)
	app.RouterAdminApi.Put("/server_groups/:id/servers/:groupServerId", app.HasPermission("server:manage"), updateServerGroupServerHandler)

	// 公告相关
	app.RouterAdminApi.Get("/announcements", app.HasPermission("announcement:view"), getAnnouncementsHandler)
	app.RouterAdminApi.Post("/announcements", app.HasPermission("announcement:manage"), createAnnouncementHandler)
	app.RouterAdminApi.Put("/announcements/:id", app.HasPermission("announcement:manage"), updateAnnouncementHandler)
	app.RouterAdminApi.Delete("/announcements/:id", app.HasPermission("announcement:manage"), deleteAnnouncementHandler)

	return nil
}
