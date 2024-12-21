package enum

// 定义模块优先级
const (
	ModulePriorityUser         = 1000
	ModulePriorityRole         = 990
	ModulePriorityPermission   = 980
	ModulePriorityLogin        = 970
	ModulePriorityMenu         = 960
	ModulePriorityAdminLog     = 950
	ModulePriorityStats        = 800
	ModulePriorityGameLog      = 790
	ModulePriorityServerConf   = 780
	ModulePriorityShell        = 770
	ModulePriorityBrowse       = 760
	ModulePriorityUnibuild     = 750
	ModulePrioritySVN          = 740
	ModulePriorityGit          = 730
	ModulePriorityCiTask       = 600
	ModulePriorityBugTracker   = 500
	ModulePriorityNote         = 400
	ModulePriorityFileManager  = 300
	ModulePriorityImageManager = 200
	ModulePriorityLuban        = 100
	ModulePriorityGameConf     = 90 // 游戏配置模块
)
