package unibuild

import (
	"github.com/andycai/goapi/core"
)

const ModulePriorityUnibuild = 9911 // 游戏-游戏构建

var app *core.App

type uniBuildModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&uniBuildModule{}, ModulePriorityUnibuild)
}

func (m *uniBuildModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *uniBuildModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *uniBuildModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/unibuild/res", buildResourcesHandler)
	app.RouterPublicApi.Post("/unibuild/app", buildAppHandler)

	return nil
}

func (m *uniBuildModule) AddAuthRouters() error {
	return nil
}
