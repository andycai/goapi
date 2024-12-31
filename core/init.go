package core

import (
	"github.com/andycai/unitool/core/collections"
)

type Module interface {
	Awake(*App) error
	Start() error
	AddPublicRouters() error
	AddAuthRouters() error
}

type BaseModule struct {
	App *App
}

func (m *BaseModule) Awake(a *App) error {
	m.App = a
	return nil
}

func (m *BaseModule) Start() error {
	return nil
}

func (m *BaseModule) AddPublicRouters() error {
	return nil
}

func (m *BaseModule) AddAuthRouters() error {
	return nil
}

var modules = collections.NewMinPriorityQueue[Module]()

// RegisterModule 注册模块
func RegisterModule(module Module, priority int) {
	modules.Enqueue(module, priority)
}

// InitPublicRouters 初始化公共路由
func InitPublicRouters() {
	modules.ForEach(func(module Module, priority int) bool {
		module.AddPublicRouters()
		return true
	})
}

// InitAuthRouters 初始化管理员路由
func InitAuthRouters() {
	modules.ForEach(func(module Module, priority int) bool {
		module.AddAuthRouters()
		return true
	})
}

// AwakeModules 模块初始化
func AwakeModules(app *App) {
	// 模块初始化
	modules.ForEach(func(module Module, priority int) bool {
		module.Awake(app)
		return true
	})

	// 模块启动
	modules.ForEach(func(module Module, priority int) bool {
		module.Start()
		return true
	})
}
