package post

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const (
	ModulePriorityPost = 2003 // 数据中心-文章管理
)

var app *core.App

type postModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&postModule{}, ModulePriorityPost)
}

func (m *postModule) Awake(a *core.App) error {
	app = a

	// 数据迁移
	return autoMigrate()
}

func (m *postModule) Start() error {
	// 初始化数据
	if err := initData(); err != nil {
		return err
	}

	return nil
}

func (m *postModule) AddPublicRouters() error {
	// 公开API
	app.RouterPublic.Get("/posts", listPostHandler)
	app.RouterPublic.Get("/posts/:slug", getPostBySlugHandler)
	app.RouterPublic.Get("/posts/search", searchPostHandler)

	return nil
}

func (m *postModule) AddAuthRouters() error {
	// 管理页面
	app.RouterAdmin.Get("/post", app.HasPermission("post:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/post", fiber.Map{
			"Title": "文章管理",
			"Scripts": []string{
				"/static/js/admin/post.js",
			},
		}, "admin/layout")
	})

	// API路由
	postGroup := app.RouterAdminApi.Group("/post")
	{
		postGroup.Get("/list", app.HasPermission("post:view"), listPostHandler)
		postGroup.Post("/add", app.HasPermission("post:add"), addPostHandler)
		postGroup.Post("/edit", app.HasPermission("post:edit"), editPostHandler)
		postGroup.Post("/delete", app.HasPermission("post:delete"), deletePostHandler)
	}

	return nil
}
