package enum

// 定义模块优先级
const (
	ModulePriorityUser         = 1000 // 系统-用户管理
	ModulePriorityRole         = 1001 // 系统-角色管理
	ModulePriorityPermission   = 1002 // 系统-权限管理
	ModulePriorityParameter    = 1003 // 系统-参数管理
	ModulePriorityLogin        = 1004 // 系统-登录管理
	ModulePriorityMenu         = 1005 // 系统-菜单管理
	ModulePriorityAdminLog     = 1006 // 系统-管理员活动日志
	ModulePriorityStats        = 2000 // 游戏-性能统计
	ModulePriorityGameLog      = 2001 // 游戏-日志
	ModulePriorityServerConf   = 2002 // 游戏-服务器配置
	ModulePriorityShell        = 2003 // 游戏-命令行
	ModulePriorityBrowse       = 2004 // 游戏-文件浏览
	ModulePriorityUnibuild     = 2005 // 游戏-Unity构建
	ModulePriorityGameConf     = 2006 // 游戏-配置
	ModulePriorityLuban        = 2007 // 游戏-Luban
	ModulePrioritySVN          = 3000 // 工具-SVN
	ModulePriorityGit          = 3001 // 工具-Git
	ModulePriorityRepoSync     = 3002 // 工具-仓库同步
	ModulePriorityPatch        = 3003 // 工具-补丁管理
	ModulePriorityFileManager  = 8000 // 功能-文件管理
	ModulePriorityImageManager = 8001 // 功能-图片管理
	ModulePriorityBugTracker   = 8002 // 功能-Bug 跟踪
	ModulePriorityNote         = 8003 // 功能-笔记
	ModulePriorityCiTask       = 8004 // 功能-CI/CD 任务
)
