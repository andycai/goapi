package imagemanager

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type imagemanagerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&imagemanagerModule{}, enum.ModulePriorityImageManager)
}

func (m *imagemanagerModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// Initialize Imagemanager service
	initService()

	// 初始化数据
	return initData()
}

func (m *imagemanagerModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/imagemanager", app.HasPermission("imagemanager:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/imagemanager", fiber.Map{
			"Title": "图片管理",
			"Scripts": []string{
				"/static/js/admin/imagemanager.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterApi.Get("/imagemanager/list", app.HasPermission("imagemanager:list"), ListHandler)
	app.RouterApi.Post("/imagemanager/upload", app.HasPermission("imagemanager:upload"), UploadHandler)
	app.RouterApi.Post("/imagemanager/delete", app.HasPermission("imagemanager:delete"), DeleteHandler)
	app.RouterApi.Post("/imagemanager/rename", app.HasPermission("imagemanager:rename"), RenameHandler)
	app.RouterApi.Post("/imagemanager/move", app.HasPermission("imagemanager:move"), MoveHandler)
	app.RouterApi.Post("/imagemanager/copy", app.HasPermission("imagemanager:copy"), CopyHandler)
	app.RouterApi.Get("/imagemanager/info", app.HasPermission("imagemanager:info"), InfoHandler)
	app.RouterApi.Get("/imagemanager/thumbnail", app.HasPermission("imagemanager:list"), ThumbnailHandler)
	app.RouterApi.Get("/imagemanager/view", app.HasPermission("imagemanager:list"), ViewHandler)

	return nil
}
