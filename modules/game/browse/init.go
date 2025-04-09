package browse

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
)

var app *core.App

type browseModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&browseModule{}, enum.ModulePriorityBrowse)
}

func (m *browseModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *browseModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *browseModule) AddAuthRouters() error {
	// admin
	// 浏览目录和文件的路由
	app.RouterAdmin.Get("/browse/*", app.HasPermission("browse:view"), listFilesHandler)

	// 文件删除路由
	app.RouterAdmin.Delete("/browse/*", app.HasPermission("browse:delete"), deleteFileHandler)

	// FTP 上传路由
	app.RouterAdmin.Post("/ftp/upload", app.HasPermission("browse:ftp"), uploadFileHandler)

	// api

	return nil
}
