function serverGroupManagement() {
    return {
        serverGroups: [],
        serverGroup: {
            name: ''
        },
        editingServerGroup: {
            id: 0,
            name: ''
        },
        groupServers: [],
        availableServers: [],
        showEditModal: false,
        showManageServersModal: false,
        showAddServerModal: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,

        init() {
            this.loadServerGroups();
        },

        async loadServerGroups() {
            try {
                const response = await fetch(`/api/admin/server_groups?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) {
                    throw new Error('Failed to load server groups');
                }
                const data = await response.json();
                this.serverGroups = data.groups;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading server groups:', error);
                ShowError('加载服务器组列表失败');
            }
        },

        async createServerGroup() {
            try {
                const response = await fetch('/api/admin/server_groups', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.serverGroup)
                });

                if (!response.ok) {
                    throw new Error('Failed to create server group');
                }

                this.resetServerGroupForm();
                this.loadServerGroups();
                ShowMessage('创建服务器组成功');
            } catch (error) {
                console.error('Error creating server group:', error);
                ShowError('创建服务器组失败');
            }
        },

        async updateServerGroup() {
            try {
                const response = await fetch(`/api/admin/server_groups/${this.editingServerGroup.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.editingServerGroup)
                });

                if (!response.ok) {
                    throw new Error('Failed to update server group');
                }

                this.showEditModal = false;
                this.loadServerGroups();
                ShowMessage('更新服务器组成功');
            } catch (error) {
                console.error('Error updating server group:', error);
                ShowError('更新服务器组失败');
            }
        },

        async deleteServerGroup(id) {
            if (!confirm('确定要删除这个服务器组吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/server_groups/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    throw new Error('Failed to delete server group');
                }

                this.loadServerGroups();
                ShowMessage('删除服务器组成功');
            } catch (error) {
                console.error('Error deleting server group:', error);
                ShowError('删除服务器组失败');
            }
        },

        async loadGroupServers(groupId) {
            try {
                const response = await fetch(`/api/admin/server_groups/${groupId}/servers`);
                if (!response.ok) {
                    throw new Error('Failed to load group servers');
                }
                const data = await response.json();
                this.groupServers = data.data;
            } catch (error) {
                console.error('Error loading group servers:', error);
                ShowError('加载服务器组中的服务器失败');
            }
        },

        async loadAvailableServers() {
            try {
                const response = await fetch('/api/admin/physical_servers');
                if (!response.ok) {
                    throw new Error('Failed to load available servers');
                }
                const data = await response.json();
                this.availableServers = data.data;
            } catch (error) {
                console.error('Error loading available servers:', error);
                ShowError('加载可用服务器列表失败');
            }
        },

        async addServerToGroup(groupId, serverId) {
            try {
                const response = await fetch(`/api/admin/server_groups/${groupId}/servers/${serverId}`, {
                    method: 'POST'
                });

                if (!response.ok) {
                    throw new Error('Failed to add server to group');
                }

                this.loadGroupServers(groupId);
                this.loadAvailableServers();
                ShowMessage('添加服务器成功');
            } catch (error) {
                console.error('Error adding server to group:', error);
                ShowError('添加服务器失败');
            }
        },

        async removeServerFromGroup(groupId, serverId) {
            if (!confirm('确定要从服务器组中移除这个服务器吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/server_groups/${groupId}/servers/${serverId}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    throw new Error('Failed to remove server from group');
                }

                this.loadGroupServers(groupId);
                ShowMessage('移除服务器成功');
            } catch (error) {
                console.error('Error removing server from group:', error);
                ShowError('移除服务器失败');
            }
        },

        editServerGroup(group) {
            this.editingServerGroup = { ...group };
            this.showEditModal = true;
        },

        manageServers(group) {
            this.editingServerGroup = { ...group };
            this.showManageServersModal = true;
            this.loadGroupServers(group.id);
        },

        showAddServer() {
            this.showAddServerModal = true;
            this.loadAvailableServers();
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) {
                return;
            }
            this.currentPage = page;
            this.loadServerGroups();
        },

        resetServerGroupForm() {
            this.serverGroup = {
                name: ''
            };
        },

        formatDate(dateString) {
            const date = new Date(dateString);
            return date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            });
        },

        getStatusText(status) {
            switch (status) {
                case 0:
                    return '维护中';
                case 1:
                    return '正常';
                case 2:
                    return '爆满';
                default:
                    return '未知';
            }
        }
    };
} 