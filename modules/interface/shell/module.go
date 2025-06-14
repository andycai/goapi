package shell

import (
	"github.com/andycai/goapi/internal"
)

const ModulePriorityShell = 9001 // 接口-shell脚本

var app *internal.App

type shellModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&shellModule{}, ModulePriorityShell)
}

func (m *shellModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// 初始化数据
	return initData()
}

func (m *shellModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/shell", execScriptHandler)

	return nil
}
