package imagemanager

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityImageManager = 5003 // 功能-图片管理

var app *core.App

type imagemanagerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&imagemanagerModule{}, ModulePriorityImageManager)
}

func (m *imagemanagerModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// Initialize Imagemanager service
	initService()
	return nil
}

func (m *imagemanagerModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *imagemanagerModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/imagemanager", app.HasPermission("imagemanager:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/imagemanager", fiber.Map{
			"Title": "图片管理",
			"Scripts": []string{
				"/static/js/admin/imagemanager.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterAdminApi.Get("/imagemanager/list", app.HasPermission("imagemanager:view"), listFilesHandler)
	app.RouterAdminApi.Post("/imagemanager/upload", app.HasPermission("imagemanager:upload"), uploadHandler)
	app.RouterAdminApi.Post("/imagemanager/delete", app.HasPermission("imagemanager:delete"), deleteHandler)
	app.RouterAdminApi.Post("/imagemanager/rename", app.HasPermission("imagemanager:rename"), renameHandler)
	app.RouterAdminApi.Post("/imagemanager/move", app.HasPermission("imagemanager:move"), moveHandler)
	app.RouterAdminApi.Post("/imagemanager/copy", app.HasPermission("imagemanager:copy"), copyHandler)
	app.RouterAdminApi.Get("/imagemanager/info", app.HasPermission("imagemanager:info"), infoHandler)
	app.RouterAdminApi.Get("/imagemanager/thumbnail", app.HasPermission("imagemanager:view"), thumbnailHandler)
	app.RouterAdminApi.Get("/imagemanager/view", app.HasPermission("imagemanager:view"), viewHandler)

	return nil
}
