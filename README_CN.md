# UnityTool

UnityTool 是一个为 Unity 游戏开发团队设计的综合性 Web
管理平台。它提供了一套工具和服务，用于简化游戏开发工作流程，包括版本控制管理、配置管理和开发工具集成。

## 功能特性

- **版本控制集成**
  - 支持 SVN 和 Git
  - 代码仓库浏览
  - 分支和提交管理
  - 代码审查功能

- **游戏配置管理**
  - 游戏配置编辑器
  - 配置版本控制
  - 多环境支持
  - 配置验证

- **开发工具**
  - Luban 配置工具集成
  - 图片资源管理
  - 文件管理系统
  - Bug 跟踪系统

- **团队协作**
  - 用户管理和认证
  - 基于角色的访问控制
  - 团队笔记和文档
  - 活动日志记录

- **CI/CD 集成**
  - Unity 构建自动化
  - 构建任务管理
  - 服务器配置管理
  - 部署自动化

## 项目结构

```
.
├── core/           # 核心应用框架
├── modules/        # 功能模块
│   ├── adminlog/   # 管理员活动日志
│   ├── browse/     # 仓库浏览器
│   ├── bugtracker/ # Bug 跟踪系统
│   ├── citask/     # CI/CD 任务管理
│   ├── filemanager/# 文件管理
│   ├── gameconf/   # 游戏配置
│   ├── git/        # Git 集成
│   ├── luban/      # Luban 工具集成
│   ├── svn/        # SVN 集成
│   └── ...         # 其他模块
├── models/         # 数据模型
├── utils/          # 实用工具函数
├── templates/      # HTML 模板
├── public/         # 静态资源
└── sql/           # 数据库架构
```

## 技术栈

- **后端**
  - Go 1.22+
  - Fiber Web 框架
  - GORM 数据库 ORM
  - HTML 模板引擎

- **前端**
  - HTML/CSS/JavaScript
  - 现代 Web 组件
  - 响应式设计

## 快速开始

### 环境要求

- Go 1.22 或更高版本
- MySQL/MariaDB
- Git
- SVN（可选）

### 安装步骤

1. 克隆仓库：
   ```bash
   git clone https://github.com/andycai/unitool.git
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 配置应用：
   - 复制 `conf.toml.example` 到 `conf.toml`
   - 更新配置设置

4. 构建应用：
   ```bash
   ./build.sh    # Unix 类系统
   build.bat     # Windows 系统
   ```

### 运行

1. 启动服务器：
   ```bash
   ./start.sh    # Unix 类系统
   start.bat     # Windows 系统
   ```

2. 访问 Web 界面：`http://localhost:8080`

### 配置说明

主配置文件 `conf.toml` 包含以下设置：

- 服务器配置
- 数据库连接
- 版本控制系统
- 认证设置
- 构建路径和选项

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 贡献指南

欢迎提交贡献！请随时提交 Pull Request。

1. Fork 本仓库
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建新的 Pull Request
