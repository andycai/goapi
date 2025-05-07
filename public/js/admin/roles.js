// Role management functionality
function roleManagement() {
    return {
        roles: [],
        permissions: [],
        currentRole: {
            id: null,
            name: '',
            description: '',
            permissions: []
        },
        showPanel: false,
        isEditing: false,
        panelTitle: '创建角色',

        init() {
            this.loadRoles();
            this.loadPermissions();
        },

        async loadRoles() {
            try {
                const response = await fetch('/api/admin/roles');
                if (!response.ok) throw new Error('Failed to load roles');
                this.roles = await response.json();
            } catch (error) {
                console.error('Error loading roles:', error);
                ShowError('加载角色列表失败');
            }
        },

        async loadPermissions() {
            try {
                const response = await fetch('/api/admin/permissions/all');
                if (!response.ok) throw new Error('Failed to load permissions');
                data = await response.json();
                this.permissions = data.items;
            } catch (error) {
                console.error('Error loading permissions:', error);
                ShowError('加载权限列表失败');
            }
        },

        openCreatePanel() {
            this.currentRole = {
                id: null,
                name: '',
                description: '',
                permissions: []
            };
            this.isEditing = false;
            this.panelTitle = '创建角色';
            this.showPanel = true;
        },

        editRole(role) {
            this.currentRole = {
                id: role.id,
                name: role.name,
                description: role.description,
                permissions: role.permissions.map(p => parseInt(p.id))
            };
            this.isEditing = true;
            this.panelTitle = '编辑角色';
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
        },

        async createRole() {
            try {
                const formData = {
                    ...this.currentRole,
                    permissions: this.currentRole.permissions.map(id => parseInt(id))
                };
                const response = await fetch('/api/admin/roles', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) throw new Error('Failed to create role');
                
                await this.loadRoles();
                this.closePanel();
                ShowMessage('角色创建成功');
            } catch (error) {
                console.error('Error creating role:', error);
                ShowError('创建角色失败');
            }
        },

        async updateRole() {
            try {
                const formData = {
                    ...this.currentRole,
                    permissions: this.currentRole.permissions.map(id => parseInt(id))
                };
                const response = await fetch(`/api/admin/roles/${this.currentRole.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) throw new Error('Failed to update role');
                
                await this.loadRoles();
                this.closePanel();
                ShowMessage('角色更新成功');
            } catch (error) {
                console.error('Error updating role:', error);
                ShowError('更新角色失败');
            }
        },

        async deleteRole(id) {
            if (!confirm('确定要删除这个角色吗？')) return;

            try {
                const response = await fetch(`/api/admin/roles/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) throw new Error('Failed to delete role');
                
                await this.loadRoles();
                ShowMessage('角色删除成功');
            } catch (error) {
                console.error('Error deleting role:', error);
                ShowError('删除角色失败');
            }
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