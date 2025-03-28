package unibuild

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
)

var app *core.App

type uniBuildModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&uniBuildModule{}, enum.ModulePriorityUnibuild)
}

func (m *uniBuildModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// 初始化数据
	return initData()
}

func (m *uniBuildModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/unibuild/res", buildResources)
	app.RouterPublicApi.Post("/unibuild/app", buildApp)

	return nil
}

func (m *uniBuildModule) AddAuthRouters() error {
	return nil
}
