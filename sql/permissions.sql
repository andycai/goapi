-- 系统管理权限
INSERT INTO permissions (id, name, description) VALUES
(1, 'system', '系统管理'),
(2, 'user:list', '用户列表'),
(3, 'user:create', '创建用户'),
(4, 'user:update', '更新用户'),
(5, 'user:delete', '删除用户'),
(6, 'role:list', '角色列表'),
(7, 'role:create', '创建角色'),
(8, 'role:update', '更新角色'),
(9, 'role:delete', '删除角色'),
(10, 'permission:list', '权限列表'),
(11, 'permission:create', '创建权限'),
(12, 'permission:update', '更新权限'),
(13, 'permission:delete', '删除权限'),
(14, 'menu:list', '菜单列表'),
(15, 'menu:create', '创建菜单'),
(16, 'menu:update', '更新菜单'),
(17, 'menu:delete', '删除菜单'),
(18, 'adminlog:list', '操作日志列表'),
(19, 'adminlog:delete', '删除操作日志');

-- 游戏管理权限
INSERT INTO permissions (id, name, description) VALUES
(19, 'game', '游戏管理'),
(20, 'gamelog:list', '游戏日志列表'),
(21, 'gamelog:create', '创建游戏日志'),
(22, 'gamelog:delete', '删除游戏日志'),
(23, 'stats:list', '统计列表'),
(24, 'stats:create', '创建统计'),
(25, 'stats:delete', '删除统计');

-- 系统工具权限
INSERT INTO permissions (id, name, description) VALUES
(26, 'tools', '系统工具'),
(27, 'file:list', '文件列表'),
(28, 'file:ftp', 'FTP上传'),
(29, 'serverconf:list', '服务器配置'),
(30, 'tools:terminal', '命令执行'),
(31, 'package:list', '打包工具'),
(32, 'citask:list', '任务列表'),
(33, 'citask:create', '创建任务'),
(34, 'citask:update', '更新任务'),
(35, 'citask:delete', '删除任务'),
(36, 'citask:run', '执行任务'),
(37, 'svn:list', 'SVN列表'),
(38, 'svn:checkout', 'SVN检出'),
(39, 'svn:update', 'SVN更新'),
(40, 'svn:commit', 'SVN提交'),
(41, 'svn:status', 'SVN状态'),
(42, 'svn:info', 'SVN信息'),
(43, 'svn:log', 'SVN日志'),
(44, 'svn:revert', 'SVN还原'),
(45, 'svn:add', 'SVN添加'),
(46, 'svn:delete', 'SVN删除'),
(47, 'git:list', 'Git列表'),
(48, 'git:clone', 'Git克隆'),
(49, 'git:pull', 'Git拉取'),
(50, 'git:push', 'Git推送'),
(51, 'git:status', 'Git状态'),
(52, 'git:log', 'Git日志'),
(53, 'git:commit', 'Git提交'),
(54, 'git:checkout', 'Git检出'),
(55, 'git:branch', 'Git分支'),
(56, 'git:merge', 'Git合并'),
(57, 'git:reset', 'Git重置'),
(58, 'git:stash', 'Git暂存');

-- Filemanager permissions
INSERT INTO permission (name, description) VALUES ('filemanager:list', '文件管理-列表');
INSERT INTO permission (name, description) VALUES ('filemanager:upload', '文件管理-上传');
INSERT INTO permission (name, description) VALUES ('filemanager:create', '文件管理-创建');
INSERT INTO permission (name, description) VALUES ('filemanager:delete', '文件管理-删除');
INSERT INTO permission (name, description) VALUES ('filemanager:rename', '文件管理-重命名');
INSERT INTO permission (name, description) VALUES ('filemanager:move', '文件管理-移动');
INSERT INTO permission (name, description) VALUES ('filemanager:copy', '文件管理-复制');
INSERT INTO permission (name, description) VALUES ('filemanager:download', '文件管理-下载');
INSERT INTO permission (name, description) VALUES ('filemanager:info', '文件管理-信息');

