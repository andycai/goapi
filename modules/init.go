package modules

// 新增的模块必须在这里进行导入，不然模块 init 方法不会执行
import (
	_ "github.com/andycai/unitool/modules/datacenter/parameter" // 参数配置
	_ "github.com/andycai/unitool/modules/development/citask"   // 构建任务
	_ "github.com/andycai/unitool/modules/game/browse"          // 浏览文件
	_ "github.com/andycai/unitool/modules/game/gamelog"         // 游戏日志
	_ "github.com/andycai/unitool/modules/game/patch"           // 补丁管理
	_ "github.com/andycai/unitool/modules/game/reposync"        // 仓库文件同步
	_ "github.com/andycai/unitool/modules/game/serverconf"      // 服务器配置
	_ "github.com/andycai/unitool/modules/game/stats"           // 游戏计
	_ "github.com/andycai/unitool/modules/game/unibuild"        // Unity构建
	_ "github.com/andycai/unitool/modules/interface/shell"      // 命令脚本执行
	_ "github.com/andycai/unitool/modules/knowledge/note"       // 笔记
	_ "github.com/andycai/unitool/modules/login"                // 登录
	_ "github.com/andycai/unitool/modules/system/adminlog"      // 后台操作日志
	_ "github.com/andycai/unitool/modules/system/menu"          // 菜单
	_ "github.com/andycai/unitool/modules/system/permission"    // 权限
	_ "github.com/andycai/unitool/modules/system/role"          // 角色
	_ "github.com/andycai/unitool/modules/system/user"          // 用户
	_ "github.com/andycai/unitool/modules/webapp/fund"          // 基金
	// _ "github.com/andycai/unitool/modules/datacenter/dict"       // 字典
	// _ "github.com/andycai/unitool/modules/development/bugtracker" // 缺陷管理
	// _ "github.com/andycai/unitool/modules/knowledge/filemanager" // 文件管理
	// _ "github.com/andycai/unitool/modules/knowledge/imagemanager" // 图片管理
)
