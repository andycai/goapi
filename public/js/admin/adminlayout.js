document.addEventListener('alpine:init', () => {
    Alpine.store('notification', {
        notifications: [],
        notificationCount: 0,
        show(message, type = 'success') {
            const id = ++this.notificationCount;
            const notification = {
                id,
                message,
                type,
                show: true,
                progress: 100
            };

            this.notifications.push(notification);

            setTimeout(() => {
                notification.progress = 0;
                setTimeout(() => {
                    notification.show = false;
                    setTimeout(() => {
                        this.notifications = this.notifications.filter(n => n.id !== id);
                    }, 300);
                }, 2000);
            }, 100);
        }
    });
});

// 显示普通消息
function ShowMessage(message) {
    Alpine.store('notification').show(message, 'success');
}

// 显示错误消息
function ShowError(message) {
    Alpine.store('notification').show(message, 'error');
}

// 格式化日期
function FormatDate(timestamp) {
    if (!timestamp) return '';
    return new Date(timestamp).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function adminLayoutManagement() {
    return {
        user: JSON.parse(localStorage.getItem('user') || '{}'),
        collapsed: localStorage.getItem('menuCollapsed') === 'true',
        mobileMenuOpen: false,
        theme: localStorage.theme || 'light',
        currentPath: window.location.pathname,
        expandedGroup: localStorage.getItem('expandedMenuGroup'),
        recentTabs: [],
        maxTabs: 8,
        menuTree: [],
        loading: false,
        justCollapsed: false,
        menuIcons: {
            default: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>',
            home: `
<svg
  class="lucide lucide-house"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8" />
  <path d="M3 10a2 2 0 0 1 .709-1.528l7-5.999a2 2 0 0 1 2.582 0l7 5.999A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
</svg>
`,
            system: `
<svg
  class="lucide lucide-presentation"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M2 3h20" />
  <path d="M21 3v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V3" />
  <path d="m7 21 5-5 5 5" />
</svg>
`,
            filemanager: `
<svg
  class="lucide lucide-files"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M20 7h-3a2 2 0 0 1-2-2V2" />
  <path d="M9 18a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h7l4 4v10a2 2 0 0 1-2 2Z" />
  <path d="M3 7.6v12.8A1.6 1.6 0 0 0 4.6 22h9.8" />
</svg>
`,
            imagemanager: `
<svg
  class="lucide lucide-images"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M18 22H4a2 2 0 0 1-2-2V6" />
  <path d="m22 13-1.296-1.296a2.41 2.41 0 0 0-3.408 0L11 18" />
  <circle cx="12" cy="8" r="2" />
  <rect width="16" height="16" x="6" y="2" rx="2" />
</svg>
`,
            reposync: `
<svg
  class="lucide lucide-folder-sync"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M9 20H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h3.9a2 2 0 0 1 1.69.9l.81 1.2a2 2 0 0 0 1.67.9H20a2 2 0 0 1 2 2v.5" />
  <path d="M12 10v4h4" />
  <path d="m12 14 1.535-1.605a5 5 0 0 1 8 1.5" />
  <path d="M22 22v-4h-4" />
  <path d="m22 18-1.535 1.605a5 5 0 0 1-8-1.5" />
</svg>
`,
            user: `
<svg
  class="lucide lucide-circle-user-round"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M18 20a6 6 0 0 0-12 0" />
  <circle cx="12" cy="10" r="4" />
  <circle cx="12" cy="12" r="10" />
</svg>
`,
            role: `
<svg
  class="lucide lucide-id-card"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M16 10h2" />
  <path d="M16 14h2" />
  <path d="M6.17 15a3 3 0 0 1 5.66 0" />
  <circle cx="9" cy="11" r="2" />
  <rect x="2" y="5" width="20" height="14" rx="2" />
</svg>
`,
            permission: `
<svg
  class="lucide lucide-lock-keyhole"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <circle cx="12" cy="16" r="1" />
  <rect x="3" y="10" width="18" height="12" rx="2" />
  <path d="M7 10V7a5 5 0 0 1 10 0v3" />
</svg>
`,
            menu: `
<svg
  class="lucide lucide-menu"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M4 12h16" />
  <path d="M4 18h16" />
  <path d="M4 6h16" />
</svg>
`,
            adminlog: `
<svg
  class="lucide lucide-chart-gantt"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M10 6h8" />
  <path d="M12 16h6" />
  <path d="M3 3v16a2 2 0 0 0 2 2h16" />
  <path d="M8 11h7" />
</svg>
`,
            game: `
<svg
  class="lucide lucide-gamepad-2"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <line x1="6" x2="10" y1="11" y2="11" />
  <line x1="8" x2="8" y1="9" y2="13" />
  <line x1="15" x2="15.01" y1="12" y2="12" />
  <line x1="18" x2="18.01" y1="10" y2="10" />
  <path d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.545-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z" />
</svg>
`,
            gamelog: `
<svg
  class="lucide lucide-logs"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M13 12h8" />
  <path d="M13 18h8" />
  <path d="M13 6h8" />
  <path d="M3 12h1" />
  <path d="M3 18h1" />
  <path d="M3 6h1" />
  <path d="M8 12h1" />
  <path d="M8 18h1" />
  <path d="M8 6h1" />
</svg>
`,
            stats: `
<svg
  class="lucide lucide-chart-column-big"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M3 3v16a2 2 0 0 0 2 2h16" />
  <rect x="15" y="5" width="4" height="12" rx="1" />
  <rect x="7" y="8" width="4" height="9" rx="1" />
</svg>
`,
            tools: `
<svg
  class="lucide lucide-wrench"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
</svg>
`,
            browse: `
<svg
  class="lucide lucide-folders"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M20 17a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3.9a2 2 0 0 1-1.69-.9l-.81-1.2a2 2 0 0 0-1.67-.9H8a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2Z" />
  <path d="M2 8v11a2 2 0 0 0 2 2h14" />
</svg>
`,
            serverconf: `
<svg
  class="lucide lucide-badge-info"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z" />
  <line x1="12" x2="12" y1="16" y2="12" />
  <line x1="12" x2="12.01" y1="8" y2="8" />
</svg>
`,
            citask: `
<svg
  class="lucide lucide-timer"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <line x1="10" x2="14" y1="2" y2="2" />
  <line x1="12" x2="15" y1="14" y2="11" />
  <circle cx="12" cy="14" r="8" />
</svg>
`,
            patch: `
<svg
  class="lucide lucide-send-to-back"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <rect x="14" y="14" width="8" height="8" rx="2" />
  <rect x="2" y="2" width="8" height="8" rx="2" />
  <path d="M7 14v1a2 2 0 0 0 2 2h1" />
  <path d="M14 7h1a2 2 0 0 1 2 2v1" />
</svg>
`,
            parameter: `
<svg
  class="lucide lucide-scan-text"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M3 7V5a2 2 0 0 1 2-2h2" />
  <path d="M17 3h2a2 2 0 0 1 2 2v2" />
  <path d="M21 17v2a2 2 0 0 1-2 2h-2" />
  <path d="M7 21H5a2 2 0 0 1-2-2v-2" />
  <path d="M7 8h8" />
  <path d="M7 12h10" />
  <path d="M7 16h6" />
</svg>
`,
            luban: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" /></svg>',
            webapp: `
<svg
  class="lucide lucide-chrome"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <circle cx="12" cy="12" r="10" />
  <circle cx="12" cy="12" r="4" />
  <line x1="21.17" x2="12" y1="8" y2="8" />
  <line x1="3.95" x2="8.54" y1="6.06" y2="14" />
  <line x1="10.88" x2="15.46" y1="21.94" y2="14" />
</svg>
`,
            notes: `
<svg
  class="lucide lucide-notebook-pen"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M13.4 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-7.4" />
  <path d="M2 6h4" />
  <path d="M2 10h4" />
  <path d="M2 14h4" />
  <path d="M2 18h4" />
  <path d="M21.378 5.626a1 1 0 1 0-3.004-3.004l-5.01 5.012a2 2 0 0 0-.506.854l-.837 2.87a.5.5 0 0 0 .62.62l2.87-.837a2 2 0 0 0 .854-.506z" />
</svg>
`,
            physical_server: `
<svg
  class="lucide lucide-server"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <rect width="20" height="8" x="2" y="2" rx="2" ry="2" />
  <rect width="20" height="8" x="2" y="14" rx="2" ry="2" />
  <line x1="6" x2="6.01" y1="6" y2="6" />
  <line x1="6" x2="6.01" y1="18" y2="18" />
</svg>
`,
            channel: `
<svg
  class="lucide lucide-podcast"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M16.85 18.58a9 9 0 1 0-9.7 0" />
  <path d="M8 14a5 5 0 1 1 8 0" />
  <circle cx="12" cy="11" r="1" />
  <path d="M13 17a1 1 0 1 0-2 0l.5 4.5a.5.5 0 1 0 1 0Z" />
</svg>
`,
            server_group: `
<svg
  class="lucide lucide-group"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M3 7V5c0-1.1.9-2 2-2h2" />
  <path d="M17 3h2c1.1 0 2 .9 2 2v2" />
  <path d="M21 17v2c0 1.1-.9 2-2 2h-2" />
  <path d="M7 21H5c-1.1 0-2-.9-2-2v-2" />
  <rect width="7" height="5" x="7" y="7" rx="1" />
  <rect width="7" height="5" x="10" y="12" rx="1" />
</svg>
`,
            announcement: `
<svg
  class="lucide lucide-message-square-more"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
  <path d="M8 10h.01" />
  <path d="M12 10h.01" />
  <path d="M16 10h.01" />
</svg>
`,
            dict: `
<svg
  class="lucide lucide-whole-word"
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <circle cx="7" cy="12" r="3" />
  <path d="M10 9v6" />
  <circle cx="17" cy="12" r="3" />
  <path d="M14 7v8" />
  <path d="M22 17v1c0 .5-.5 1-1 1H3c-.5 0-1-.5-1-1v-1" />
</svg>
`,
            
        },
        get recentTabsKey() {
            return 'recentTabs_' + this.user.id;
        },
        async loadMenus() {
            try {
                const response = await fetch('/api/menus/public/tree');
                if (!response.ok) throw new Error('加载菜单失败');
                this.menuTree = await response.json();

                // 加载菜单后，只初始化标签，不自动展开菜单
                this.initializeTabs();

                // 如果菜单是展开状态，并且localStorage中有保存的展开组，则恢复
                if (!this.collapsed) {
                    const savedGroup = localStorage.getItem('expandedMenuGroup');
                    if (savedGroup) {
                        this.expandedGroup = parseInt(savedGroup);
                    }
                }
            } catch (error) {
                console.error('Failed to load menus:', error);
                ShowError('加载菜单失败');
            }
        },
        initializeTabs() {
            // 从 localStorage 获取保存的标签
            const savedTabs = JSON.parse(localStorage.getItem(this.recentTabsKey) || '[]');
            this.recentTabs = savedTabs;
        },
        hasPermission(permission) {
            // 如果没有配置权限，则不显示
            if (!permission) return false;

            // 检查用户是否存在且有角色
            if (!this.user || !this.user.role) return false;

            // 获取用户权限列表
            const userPermissions = this.user.role.permissions?.map(p => p.code) || [];

            // 如果用户有 admin 权限，允许访问所有内容
            if (userPermissions.includes('admin')) return true;

            // 检查具体权限
            return userPermissions.includes(permission);
        },
        hasMenuPermission(menuItem) {
            // 如果是父菜单，检查是否有可见的子菜单
            if (menuItem.children && menuItem.children.length > 0) {
                return this.hasVisibleChildren(menuItem);
            }

            // 如果是子菜单或没有子菜单的菜单项，检查自身权限
            return this.hasPermission(menuItem.menu.permission);
        },
        hasVisibleChildren(menuItem) {
            // 检查是否有任何子菜单有权限显示
            return menuItem.children && menuItem.children.some(child => this.hasPermission(child.menu.permission));
        },
        getMenuName(path) {
            // 首页特殊处理
            if (path === '/admin') return '首页';

            // 在菜单树中查找
            for (const menu of this.menuTree) {
                if (menu.menu.path === path) {
                    return menu.menu.name;
                }
                for (const child of menu.children) {
                    if (child.menu.path === path) {
                        return child.menu.name;
                    }
                }
            }
            return path;
        },
        addTab(path) {
            // 检查路径是否在菜单中定义
            const isMenuPath = this.isPathInMenu(path);
            if (!isMenuPath) {
                return; // 如果路径不在菜单中，不添加标签
            }

            // 查找菜单项以获取标题
            const menuItem = this.findMenuItemByPath(path);
            if (!menuItem) {
                return; // 如果找不到对应的菜单项，不添加标签
            }

            // 检查标签是否已存在
            const existingTabIndex = this.recentTabs.findIndex(tab => tab.path === path);
            if (existingTabIndex !== -1) {
                // 如果标签已存在，将其移动到最后
                const tab = this.recentTabs.splice(existingTabIndex, 1)[0];
                this.recentTabs.push(tab);
            } else {
                // 如果标签不存在，添加新标签
                this.recentTabs.push({
                    path: path,
                    title: menuItem.name
                });

                // 如果标签数量超过最大值，移除最早的标签
                if (this.recentTabs.length > this.maxTabs) {
                    this.recentTabs.shift();
                }
            }

            // 保存到 localStorage
            localStorage.setItem(this.recentTabsKey, JSON.stringify(this.recentTabs));
        },
        // 检查路径是否在菜单中定义
        isPathInMenu(path) {
            return this.findMenuItemByPath(path) !== null;
        },
        // 在菜单树中查找指定路径的菜单项
        findMenuItemByPath(path) {
            const searchInMenu = (items) => {
                for (const item of items) {
                    // 检查当前菜单
                    if (item.menu && item.menu.path === path) {
                        return item.menu;
                    }
                    // 检查子菜单项
                    if (item.children && item.children.length > 0) {
                        const found = searchInMenu(item.children);
                        if (found) return found;
                    }
                }
                return null;
            };
            return searchInMenu(this.menuTree);
        },
        closeTab(path) {
            const index = this.recentTabs.findIndex(tab => tab.path === path);
            if (index === -1) return;

            this.recentTabs.splice(index, 1);
            localStorage.setItem(this.recentTabsKey, JSON.stringify(this.recentTabs));

            // 如果关闭的是当前标签，导航到前一个标签
            if (path === this.currentPath) {
                const prevTab = this.recentTabs[Math.max(0, index - 1)];
                if (prevTab) {
                    this.navigate(prevTab.path);
                } else {
                    this.navigate('/admin');
                }
            }
        },
        navigate(path) {
            if (this.loading || path === this.currentPath) return;
            this.loading = true;

            try {
                window.location.href = path;
            } catch (error) {
                console.error('Navigation error:', error);
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },
        toggleMenuGroup(group) {
            if (this.expandedGroup === group) {
                this.expandedGroup = null;
                localStorage.removeItem('expandedMenuGroup');
            } else {
                this.expandedGroup = group;
                localStorage.setItem('expandedMenuGroup', group);
            }
        },
        toggleTheme() {
            this.theme = this.theme === 'light' ? 'dark' : 'light';
            localStorage.setItem('theme', this.theme);

            if (this.theme === 'dark') {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }
        },
        toggleCollapse() {
            this.collapsed = !this.collapsed;
            localStorage.setItem('menuCollapsed', this.collapsed);
            if (this.collapsed) {
                this.justCollapsed = true;
                this.expandedGroup = null;
                // 300ms 后重置 justCollapsed，这个时间要比菜单收起的动画时间长
                setTimeout(() => {
                    this.justCollapsed = false;
                }, 300);
            }
        },
        logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            // post 请求 /logout
            window.location.href = '/logout';
        },
        init() {
            // 检查登录状态
            const token = localStorage.getItem('token');
            if (!token) {
                window.location.href = '/login';
                return;
            }

            // 设置初始主题
            if (this.theme === 'dark') {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }

            // 加载菜单数据
            this.loadMenus();

            // 初始化时不显示浮动子菜单
            this.justCollapsed = true;
            setTimeout(() => {
                this.justCollapsed = false;
            }, 300);
        },
        // 查找当前路径对应的菜单项及其父菜单
        findCurrentMenu(items = this.menuTree, parent = null) {
            for (const item of items) {
                // 检查当前菜单项
                if (item.menu.path === this.currentPath) {
                    return { current: item.menu, parent };
                }
                // 检查子菜单项
                if (item.children && item.children.length > 0) {
                    const found = this.findCurrentMenu(item.children, item.menu);
                    if (found) return found;
                }
            }
            return null;
        },
        // 展开当前菜单的父菜单
        expandCurrentMenuParent() {
            const found = this.findCurrentMenu();
            if (found && found.parent) {
                this.expandedGroup = found.parent.id;
                localStorage.setItem('expandedMenuGroup', found.parent.id);
            }
        },
        async navigate(path) {
            if (this.loading || path === this.currentPath) return;
            this.loading = true;

            // 更新当前路径
            this.currentPath = path;

            // 更新标签
            const menuInfo = this.findCurrentMenu();
            if (menuInfo) {
                const { current } = menuInfo;

                // 添加到最近标签
                const existingTabIndex = this.recentTabs.findIndex(tab => tab.path === path);
                if (existingTabIndex === -1) {
                    // 添加新标签
                    this.recentTabs.push({
                        path: current.path,
                        name: current.name
                    });
                    // 如果超过最大数量，删除最早的标签
                    if (this.recentTabs.length > this.maxTabs) {
                        this.recentTabs.shift();
                    }
                    // 保存到 localStorage
                    localStorage.setItem(this.recentTabsKey, JSON.stringify(this.recentTabs));
                }
            }

            // 如果菜单是收起状态，确保不会展开
            if (this.collapsed) {
                this.expandedGroup = null;
                this.justCollapsed = true;
                setTimeout(() => {
                    this.justCollapsed = false;
                }, 300);
            }

            try {
                window.location.href = path;
            } catch (error) {
                console.error('Navigation error:', error);
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        }
    }
}

