package internal

import (
	"github.com/andycai/goapi/pkg/collections"
)

type Module interface {
	Awake(*App) error
	Start() error
	Dispose() error
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

func (m *BaseModule) Dispose() error {
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

// InitModules 模块初始化
func InitModules(app *App) {
	// Awake
	modules.ForEach(func(module Module, priority int) bool {
		module.Awake(app)
		return true
	})

	// Start
	modules.ForEach(func(module Module, priority int) bool {
		module.Start()
		return true
	})
}
