package note

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
)

var app *core.App

type noteModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&noteModule{}, enum.ModulePriorityNote)
}

func (m *noteModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *noteModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *noteModule) AddPublicRouters() error {
	// 公开API路由
	app.RouterPublicApi.Get("/notes/public", getPublicNotesHandler)
	app.RouterPublicApi.Get("/notes/public/:id", getPublicNoteDetailHandler)
	app.RouterPublicApi.Get("/notes/categories/public", getPublicCategoriesHandler)

	return nil
}

func (m *noteModule) AddAuthRouters() error {
	// 管理后台路由
	app.RouterAdmin.Get("/notes", app.HasPermission("note:view"), listNotesHandler)

	// API路由 - 调整顺序，将具体路径放在参数路径之前
	app.RouterApi.Get("/notes/tree", app.HasPermission("note:view"), getNoteTreeHandler)

	// 分类操作 - 移动到参数路由之前
	app.RouterApi.Get("/notes/categories", app.HasPermission("note:category:view"), listCategoriesHandler)
	app.RouterApi.Post("/notes/categories", app.HasPermission("note:category:create"), createCategoryHandler)
	app.RouterApi.Put("/notes/categories/:id", app.HasPermission("note:category:update"), updateCategoryHandler)
	app.RouterApi.Delete("/notes/categories/:id", app.HasPermission("note:category:delete"), deleteCategoryHandler)

	// 参数路由放在最后
	app.RouterApi.Get("/notes/:id", app.HasPermission("note:view"), getNoteDetailHandler)
	app.RouterApi.Post("/notes", app.HasPermission("note:create"), createNoteHandler)
	app.RouterApi.Put("/notes/:id", app.HasPermission("note:update"), updateNoteHandler)
	app.RouterApi.Delete("/notes/:id", app.HasPermission("note:delete"), deleteNoteHandler)

	return nil
}
