package filemanager

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type filemanagerModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&filemanagerModule{}, enum.ModulePriorityFileManager)
}

func (m *filemanagerModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *filemanagerModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	// Initialize Filemanager service
	initService()
	return nil
}

func (m *filemanagerModule) AddAuthRouters() error {
	// admin page
	app.RouterAdmin.Get("/filemanager", app.HasPermission("filemanager:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/filemanager", fiber.Map{
			"Title": "文件管理",
			"Scripts": []string{
				"/static/js/admin/filemanager.js",
			},
		}, "admin/layout")
	})

	// api routes
	app.RouterApi.Get("/filemanager/list", app.HasPermission("filemanager:view"), listFilesHandler)
	app.RouterApi.Post("/filemanager/upload", app.HasPermission("filemanager:upload"), uploadFileHandler)
	app.RouterApi.Post("/filemanager/create", app.HasPermission("filemanager:create"), createHandler)
	app.RouterApi.Post("/filemanager/delete", app.HasPermission("filemanager:delete"), deleteHandler)
	app.RouterApi.Post("/filemanager/rename", app.HasPermission("filemanager:rename"), renameHandler)
	app.RouterApi.Post("/filemanager/move", app.HasPermission("filemanager:move"), moveHandler)
	app.RouterApi.Post("/filemanager/copy", app.HasPermission("filemanager:copy"), copyHandler)
	app.RouterApi.Get("/filemanager/download", app.HasPermission("filemanager:download"), downloadHandler)
	app.RouterApi.Get("/filemanager/info", app.HasPermission("filemanager:info"), infoHandler)

	return nil
}
