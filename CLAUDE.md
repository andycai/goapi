# Goapi 项目 Claude Code 描述

## 项目概述

Goapi 是一个专为Unity游戏开发团队设计的综合性Web管理平台。该项目使用Go语言开发，采用模块化架构，提供了版本控制管理、游戏配置管理、开发工具集成和团队协作等功能。

项目主要技术栈：
- 后端：Go 1.22+、Fiber Web框架、GORM数据库ORM
- 前端：HTML/CSS/JavaScript、Tailwind CSS
- 数据库：支持SQLite和MySQL
- 模板引擎：html/template

## 项目架构

### 核心架构

项目采用模块化设计，通过internal/app.go中的App结构体管理整个应用。主要组件包括：

1. Fiber应用实例
2. 数据库连接（支持多数据库）
3. 配置管理
4. 路由系统（公共路由、API路由、管理路由）
5. 事件总线
6. 会话管理
7. 权限认证系统

### 目录结构

```
.
├── enum/              # 枚举类型定义
├── events/            # 事件定义
├── internal/          # 核心内部组件
│   ├── log/           # 日志系统（使用zap）
│   └── (核心应用文件)
├── models/            # 数据模型
├── modules/           # 功能模块（核心）
├── pkg/               # 公共包
├── public/            # 静态资源文件
├── scripts/           # 脚本文件
├── sql/               # 数据库迁移脚本
├── templates/         # HTML模板
├── utils/             # 工具函数
└── (配置和主程序文件)
```

### 模块系统

项目功能通过modules/目录下的模块组织，每个模块都是一个独立的功能单元。模块遵循统一的接口：

```go
type Module interface {
    // 初始化模块，传入应用实例
    Awake(a *App) error
    
    // 模块启动
    Start() error
    
    // 添加公共路由（无需认证）
    AddPublicRouters() error
    
    // 添加需认证路由
    AddAuthRouters() error
}
```

主要模块分类：
1. datacenter - 数据中心（基础数据、字典、页面、参数、博客等）
2. development - 开发工具（CI任务、缺陷跟踪等）
3. game - 游戏相关（配置、日志、补丁、仓库同步等）
4. interface - 接口集成（Git、SVN、Shell等）
5. knowledge - 知识管理（文件、图片、笔记等）
6. login - 登录认证
7. system - 系统管理（用户、角色、权限、菜单、操作日志等）
8. webapp - Web应用（基金等）

## 核心功能模块

### 认证与权限系统

项目实现了基于角色的访问控制（RBAC）：
- User（用户）：系统用户，关联角色
- Role（角色）：用户角色，关联权限
- Permission（权限）：系统权限，通过代码标识

权限检查通过中间件实现，在路由中可直接使用：
```go
app.RouterAdmin.Get("/path", app.HasPermission("permission_code"), handler)
```

### 配置管理

通过conf.toml文件进行配置，支持命令行参数覆盖。配置项包括：
- 应用配置（开发模式、安全模式）
- 服务器配置（主机、端口、静态路径等）
- 数据库配置（驱动、连接串、连接池等）
- FTP配置
- 认证配置（JWT密钥、过期时间）
- CORS跨域配置

### 数据库设计

使用GORM作为ORM库，支持SQLite和MySQL。主要数据表包括：
- 用户系统：users, roles, permissions
- 系统日志：admin_logs
- 业务数据：根据模块不同而不同

### 路由系统

采用分组路由管理：
1. 公共路由（/）- 无需认证
2. 公共API路由（/api）- 无需认证的API接口
3. 管理API路由（/api/admin）- 需要认证的API接口
4. 管理路由（/admin）- 需要认证的管理页面

### 模板系统

使用Go标准库html/template作为模板引擎，结合Fiber框架。支持自定义模板函数，如：
- hasSuffix：检查文件扩展名
- splitPath：分割路径
- sub：数字减法运算

前端采用Tailwind CSS作为样式框架，实现响应式设计。

## 开发特点

1. 模块化设计：功能通过独立模块实现，易于扩展和维护
2. 权限控制：细粒度的基于角色的权限管理系统
3. 配置灵活：支持文件配置和命令行参数
4. 多数据库支持：支持SQLite（开发）和MySQL（生产）
5. 静态资源管理：支持多静态路径配置
6. 日志系统：使用zap实现高性能日志记录
7. 错误处理：统一的错误处理机制，支持自定义错误页面

## 部署与运行

项目可通过build.sh脚本编译，支持通过命令行参数指定配置：
- -config：配置文件路径
- -host：主机地址
- -port：端口号
- -dev：开发模式

支持通过systemd等进程管理工具部署，静态资源通过内置静态文件服务提供。