-- 系统管理菜单
INSERT INTO menus (id, parent_id, name, path, icon, permission, sort, is_show) VALUES
(1, 0, '系统管理', '', 'system', 'system', 1, 1),
(2, 1, '用户管理', '/admin/users', 'users', 'system:users', 1, 1),
(3, 1, '角色管理', '/admin/roles', 'roles', 'system:roles', 2, 1),
(4, 1, '权限管理', '/admin/permissions', 'permissions', 'system:permissions', 3, 1),
(5, 1, '菜单管理', '/admin/menus', 'menu', 'system:menus', 4, 1),
(6, 1, '操作日志', '/admin/logs', 'logs', 'system:logs', 5, 1);

-- 游戏管理菜单
INSERT INTO menus (id, parent_id, name, path, icon, permission, sort, is_show) VALUES
(7, 0, '游戏管理', '', 'game', 'game', 2, 1),
(8, 7, '游戏日志', '/admin/game/logs', 'game-log', 'game:logs', 1, 1),
(9, 7, '性能统计', '/admin/game/stats', 'stats', 'game:stats', 2, 1);

-- 系统工具菜单
INSERT INTO menus (id, parent_id, name, path, icon, permission, sort, is_show) VALUES
(10, 0, '系统工具', '', 'tools', 'tools', 3, 1),
(11, 10, '文件浏览', '/admin/tools/files', 'files', 'tools:files', 1, 1),
(12, 10, 'FTP上传', '/admin/tools/upload', 'upload', 'tools:upload', 2, 1),
(13, 10, '服务器配置', '/admin/tools/config', 'server', 'tools:server', 3, 1),
(14, 10, '命令执行', '/admin/tools/terminal', 'terminal', 'tools:terminal', 4, 1),
(15, 10, '打包工具', '/admin/tools/package', 'package', 'tools:package', 5, 1); 