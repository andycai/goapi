# UnityTool (Goapi) 项目 Gemini CLI 描述

## 项目简介

UnityTool（又名Goapi）是一个专为Unity游戏开发团队设计的综合性Web管理平台。该项目使用Go语言开发，采用模块化架构，提供了版本控制管理、游戏配置管理、开发工具集成和团队协作等功能。

### 技术栈

- 后端：Go 1.22+、Fiber Web框架、GORM数据库ORM
- 前端：HTML/CSS/JavaScript、Tailwind CSS
- 数据库：支持SQLite（开发环境）和MySQL（生产环境）
- 模板引擎：html/template

## 项目架构概览

```
.
├── enum/              # 枚举类型定义
├── events/            # 事件定义
├── internal/          # 核心内部组件
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

## 核心组件

### 1. 模块化系统

项目采用模块化设计，通过`internal/app.go`中的App结构体管理整个应用。主要特性：

- 每个功能都是一个独立模块
- 模块遵循统一接口规范
- 支持模块优先级和依赖关系

模块接口定义：
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

### 2. 认证与权限系统

实现了基于角色的访问控制（RBAC）：
- User（用户）：系统用户，关联角色
- Role（角色）：用户角色，关联权限
- Permission（权限）：系统权限，通过代码标识

权限检查中间件使用示例：
```go
app.RouterAdmin.Get("/path", app.HasPermission("permission_code"), handler)
```

### 3. 配置管理系统

通过`conf.toml`文件进行配置，支持命令行参数覆盖。配置项包括：

- 应用配置（开发模式、安全模式）
- 服务器配置（主机、端口、静态路径等）
- 数据库配置（驱动、连接串、连接池等）
- FTP配置
- 认证配置（JWT密钥、过期时间）
- CORS跨域配置

### 4. 数据库设计

使用GORM作为ORM库，主要数据表：

- 用户系统：users, roles, permissions
- 系统日志：admin_logs
- 业务数据：根据模块不同而异

### 5. 路由系统

采用分组路由管理：

1. 公共路由（`/`）- 无需认证
2. 公共API路由（`/api`）- 无需认证的API接口
3. 管理API路由（`/api/admin`）- 需要认证的API接口
4. 管理路由（`/admin`）- 需要认证的管理页面

### 6. 模板系统

使用Go标准库html/template作为模板引擎，结合Fiber框架。支持自定义模板函数：

- `hasSuffix`：检查文件扩展名
- `splitPath`：分割路径
- `sub`：数字减法运算

## 模块分类详解

### 数据中心模块 (datacenter)
- bases: 基础数据管理
- dict: 字典管理
- page: 静态页面管理
- parameter: 参数配置
- post: 博客文章管理

### 开发工具模块 (development)
- citask: CI/CD任务管理
- bugtracker: 缺陷跟踪系统
- command: 命令执行管理

### 游戏相关模块 (game)
- browse: 文件浏览
- channel: 渠道管理
- gameconf: 游戏配置
- gamelog: 游戏日志
- patch: 补丁管理
- reposync: 仓库同步
- serverconf: 服务器配置
- stats: 游戏统计
- unibuild: Unity构建
- unitool: Unity工具集成

### 接口集成模块 (interface)
- git: Git版本控制集成
- shell: 命令行脚本执行
- svn: SVN版本控制集成

### 知识管理模块 (knowledge)
- filemanager: 文件管理
- imagemanager: 图片资源管理
- note: 笔记管理

### 系统管理模块 (system)
- adminlog: 后台操作日志
- menu: 菜单管理
- permission: 权限管理
- role: 角色管理
- user: 用户管理

## 开发与部署

### 环境要求
- Go 1.22 或更高版本
- MySQL/MariaDB（生产环境）
- SQLite（开发环境，无需额外安装）

### 启动项目

1. 克隆项目：
```bash
git clone https://github.com/andycai/unitool.git
```

2. 安装依赖：
```bash
go mod download
```

3. 配置应用：
```bash
cp conf.toml.example conf.toml
# 根据需要修改配置
```

4. 运行应用：
```bash
go run main.go
```

### 命令行参数

支持通过命令行参数覆盖配置文件：
- `-config`: 配置文件路径
- `-host`: 主机地址
- `-port`: 端口号
- `-dev`: 开发模式

示例：
```bash
go run main.go -host=0.0.0.0 -port=8080 -dev
```

### 构建与部署

使用项目提供的构建脚本：
```bash
./build.sh
```

支持通过systemd等进程管理工具部署。

## 项目特点

1. **模块化设计**：功能通过独立模块实现，易于扩展和维护
2. **权限控制**：细粒度的基于角色的权限管理系统
3. **配置灵活**：支持文件配置和命令行参数
4. **多数据库支持**：支持SQLite（开发）和MySQL（生产）
5. **静态资源管理**：支持多静态路径配置
6. **日志系统**：使用zap实现高性能日志记录
7. **错误处理**：统一的错误处理机制，支持自定义错误页面
8. **响应式前端**：采用Tailwind CSS实现现代化界面

## Gemini CLI 使用建议

当使用Gemini CLI与该项目交互时，请注意：

1. 项目结构清晰，模块划分明确，便于理解和修改
2. 数据模型定义在`models/`目录下，便于查阅数据结构
3. 每个模块的功能相对独立，可根据需要重点关注特定模块
4. 权限系统复杂但完善，建议在开发新功能时遵循现有权限设计模式
5. 配置系统灵活，支持多种部署环境
6. 模板系统使用标准Go模板语法，便于前端开发