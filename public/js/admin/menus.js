// Menu management functionality
function menuManagement() {
    return {
        menus: [],
        parentMenus: [],
        currentMenu: {
            id: 0,
            menu_id: 0,
            parent_id: 0,
            name: '',
            path: '',
            icon: '',
            permission: '',
            sort: 0,
            is_show: true
        },
        showPanel: false,
        isEditing: false,
        panelTitle: '新建菜单',
        loading: false,

        init() {
            this.loadMenus();
        },

        async loadMenus() {
            try {
                const treeResponse = await fetch('/api/admin/menus/tree');
                if (!treeResponse.ok) throw new Error('获取菜单树失败');
                this.menus = await treeResponse.json();

                const response = await fetch('/api/admin/menus');
                if (!response.ok) throw new Error('获取菜单列表失败');
                const menus = await response.json();
                this.parentMenus = menus.filter(menu => menu.parent_id === 0);
            } catch (error) {
                console.error('Error loading menus:', error);
                ShowError('加载菜单失败');
            }
        },

        get flattenedMenus() {
            if (!this.menus || !Array.isArray(this.menus) || this.menus.length === 0) {
                return [];
            }
            
            const flattened = [];
            const processMenu = (menuNode, level = 0) => {
                if (!menuNode || !menuNode.menu) return;
                
                flattened.push({ ...menuNode.menu, level });
                if (menuNode.children && Array.isArray(menuNode.children) && menuNode.children.length > 0) {
                    menuNode.children.forEach(child => {
                        processMenu(child, level + 1);
                    });
                }
            };
            
            this.menus.forEach(menuNode => processMenu(menuNode));
            return flattened;
        },

        openCreatePanel() {
            this.currentMenu = {
                menu_id: 0,
                parent_id: 0,
                name: '',
                path: '',
                icon: '',
                permission: '',
                sort: 0,
                is_show: true
            };
            this.isEditing = false;
            this.panelTitle = '新建菜单';
            this.showPanel = true;
        },

        editMenu(menu) {
            this.currentMenu = { ...menu };
            this.isEditing = true;
            this.panelTitle = '编辑菜单';
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
        },

        async createMenu() {
            try {
                this.loading = true;
                const response = await fetch('/api/admin/menus', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentMenu)
                });

                if (!response.ok) throw new Error('Failed to create menu');
                
                await this.loadMenus();
                this.closePanel();
                ShowMessage('菜单创建成功');
            } catch (error) {
                console.error('Error creating menu:', error);
                ShowError('创建菜单失败');
            } finally {
                this.loading = false;
            }
        },

        async updateMenu() {
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/menus/${this.currentMenu.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentMenu)
                });

                if (!response.ok) throw new Error('Failed to update menu');
                
                await this.loadMenus();
                this.closePanel();
                ShowMessage('菜单更新成功');
            } catch (error) {
                console.error('Error updating menu:', error);
                ShowError('更新菜单失败');
            } finally {
                this.loading = false;
            }
        },

        async deleteMenu(menuId) {
            if (!confirm('确定要删除这个菜单吗？')) return;

            try {
                const response = await fetch(`/api/admin/menus/${menuId}`, {
                    method: 'DELETE'
                });

                if (!response.ok) throw new Error('Failed to delete menu');
                
                await this.loadMenus();
                ShowMessage('菜单删除成功');
            } catch (error) {
                console.error('Error deleting menu:', error);
                ShowError('删除菜单失败');
            }
        }
    };
} 