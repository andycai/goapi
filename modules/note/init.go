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
	app.RouterPublicApi.Get("/notes/public", getPublicNotes)
	app.RouterPublicApi.Get("/notes/public/:id", getPublicNoteDetail)
	app.RouterPublicApi.Get("/notes/categories/public", getPublicCategories)

	return nil
}

func (m *noteModule) AddAuthRouters() error {
	// 管理后台路由
	app.RouterAdmin.Get("/notes", app.HasPermission("note:list"), getNoteList)

	// API路由
	app.RouterApi.Get("/notes/tree", app.HasPermission("note:list"), getNoteTree)
	app.RouterApi.Get("/notes/:id", app.HasPermission("note:list"), getNoteDetail)
	app.RouterApi.Post("/notes", app.HasPermission("note:create"), createNote)
	app.RouterApi.Put("/notes/:id", app.HasPermission("note:update"), updateNote)
	app.RouterApi.Delete("/notes/:id", app.HasPermission("note:delete"), deleteNote)

	// 分类操作
	app.RouterApi.Get("/notes/categories", app.HasPermission("note:category:list"), getCategoryList)
	app.RouterApi.Post("/notes/categories", app.HasPermission("note:category:create"), createCategory)
	app.RouterApi.Put("/notes/categories/:id", app.HasPermission("note:category:update"), updateCategory)
	app.RouterApi.Delete("/notes/categories/:id", app.HasPermission("note:category:delete"), deleteCategory)

	return nil
}
