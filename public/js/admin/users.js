// User management functionality
function userManagement() {
    return {
        users: [],
        roles: [],
        currentUser: {
            id: null,
            username: '',
            password: '',
            nickname: '',
            role_id: '',
            status: 1
        },
        showPanel: false,
        isEditing: false,
        panelTitle: '创建用户',

        init() {
            this.loadUsers();
            this.loadRoles();
        },

        async loadUsers() {
            try {
                const response = await fetch('/api/admin/users');
                if (!response.ok) throw new Error('Failed to load users');
                this.users = await response.json();
            } catch (error) {
                console.error('Error loading users:', error);
                ShowError('加载用户列表失败');
            }
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

        openCreatePanel() {
            this.currentUser = {
                id: null,
                username: '',
                password: '',
                nickname: '',
                role_id: '',
                status: 1
            };
            this.isEditing = false;
            this.panelTitle = '创建用户';
            this.showPanel = true;
        },

        editUser(user) {
            this.currentUser = {
                id: user.id,
                username: user.username,
                password: '',
                nickname: user.nickname,
                role_id: user.role_id,
                status: user.status
            };
            this.isEditing = true;
            this.panelTitle = '编辑用户';
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
        },

        async createUser() {
            try {
                const response = await fetch('/api/admin/users', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.currentUser)
                });

                if (!response.ok) throw new Error('Failed to create user');
                
                await this.loadUsers();
                this.closePanel();
                ShowMessage('用户创建成功');
            } catch (error) {
                console.error('Error creating user:', error);
                ShowError('创建用户失败');
            }
        },

        async updateUser() {
            try {
                const response = await fetch(`/api/admin/users/${this.currentUser.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.currentUser)
                });

                if (!response.ok) throw new Error('Failed to update user');
                
                await this.loadUsers();
                this.closePanel();
                ShowMessage('用户更新成功');
            } catch (error) {
                console.error('Error updating user:', error);
                ShowError('更新用户失败');
            }
        },

        async deleteUser(id) {
            if (!confirm('确定要删除这个用户吗？')) return;

            try {
                const response = await fetch(`/api/admin/users/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) throw new Error('Failed to delete user');
                
                await this.loadUsers();
                ShowMessage('用户删除成功');
            } catch (error) {
                console.error('Error deleting user:', error);
                ShowError('删除用户失败');
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