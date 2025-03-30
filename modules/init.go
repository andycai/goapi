package modules

// 新增的模块必须在这里进行导入，不然模块 init 方法不会执行
import (
	_ "github.com/andycai/unitool/modules/adminlog"   // 后台操作日志
	_ "github.com/andycai/unitool/modules/browse"     // 浏览文件
	_ "github.com/andycai/unitool/modules/citask"     // 构建任务
	_ "github.com/andycai/unitool/modules/gamelog"    // 游戏日志
	_ "github.com/andycai/unitool/modules/login"      // 登录
	_ "github.com/andycai/unitool/modules/menu"       // 菜单
	_ "github.com/andycai/unitool/modules/note"       // 笔记
	_ "github.com/andycai/unitool/modules/parameter"  // 参数配置
	_ "github.com/andycai/unitool/modules/patch"      // 补丁管理
	_ "github.com/andycai/unitool/modules/permission" // 权限
	_ "github.com/andycai/unitool/modules/reposync"   // 仓库文件同步
	_ "github.com/andycai/unitool/modules/role"       // 角色
	_ "github.com/andycai/unitool/modules/serverconf" // 服务器配置
	_ "github.com/andycai/unitool/modules/shell"      // 命令脚本执行
	_ "github.com/andycai/unitool/modules/stats"      // 游戏计
	_ "github.com/andycai/unitool/modules/unibuild"   // Unity构建
	_ "github.com/andycai/unitool/modules/user"       // 用户
	// _ "github.com/andycai/unitool/modules/bugtracker"
	// _ "github.com/andycai/unitool/modules/filemanager"
	// _ "github.com/andycai/unitool/modules/imagemanager"
)
