// Permission management functionality
function permissionManagement() {
    return {
        permissions: [],
        showPanel: false,
        isEditing: false,
        panelTitle: '创建权限',
        currentPermission: {
            id: 0,
            name: '',
            code: '',
            description: ''
        },
        loading: false,
        // Pagination
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,

        init() {
            this.loadPermissions();
        },

        async loadPermissions() {
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/permissions?page=${this.currentPage}&pageSize=${this.pageSize}`);
                if (!response.ok) throw new Error('获取权限列表失败');
                
                const data = await response.json();
                this.permissions = data.items;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading permissions:', error);
                ShowError('加载权限失败');
            } finally {
                this.loading = false;
            }
        },

        openCreatePanel() {
            this.currentPermission = {
                id: 0,
                name: '',
                code: '',
                description: ''
            };
            this.isEditing = false;
            this.panelTitle = '创建权限';
            this.showPanel = true;
        },

        editPermission(permission) {
            this.currentPermission = { ...permission };
            this.isEditing = true;
            this.panelTitle = '编辑权限';
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
        },

        async createPermission() {
            try {
                this.loading = true;
                const response = await fetch('/api/admin/permissions', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentPermission)
                });

                if (!response.ok) throw new Error('创建权限失败');
                
                await this.loadPermissions();
                this.closePanel();
                ShowMessage('权限创建成功');
            } catch (error) {
                console.error('Error creating permission:', error);
                ShowError('创建权限失败');
            } finally {
                this.loading = false;
            }
        },

        async updatePermission() {
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/permissions/${this.currentPermission.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentPermission)
                });

                if (!response.ok) throw new Error('更新权限失败');
                
                await this.loadPermissions();
                this.closePanel();
                ShowMessage('权限更新成功');
            } catch (error) {
                console.error('Error updating permission:', error);
                ShowError('更新权限失败');
            } finally {
                this.loading = false;
            }
        },

        async deletePermission(id) {
            if (!confirm('确定要删除这个权限吗？')) return;

            try {
                const response = await fetch(`/api/admin/permissions/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) throw new Error('删除权限失败');

                await this.loadPermissions();
                ShowMessage('权限删除成功');
            } catch (error) {
                console.error('Error deleting permission:', error);
                ShowError('删除权限失败');
            }
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            this.loadPermissions();
        },

        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
            });
        }
    }
} 