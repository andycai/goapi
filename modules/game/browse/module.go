package browse

import (
	"github.com/andycai/goapi/internal"
)

const ModulePriorityBrowse = 9901 // 游戏-文件浏览

var app *internal.App

type browseModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&browseModule{}, ModulePriorityBrowse)
}

func (m *browseModule) Awake(a *internal.App) error {
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
	app.RouterAdminApi.Delete("/browse/*", app.HasPermission("browse:delete"), deleteFileHandler)

	// FTP 上传路由
	app.RouterAdminApi.Post("/ftp/upload", app.HasPermission("browse:ftp"), uploadFileHandler)

	// api

	return nil
}
