function serverGroupManagement() {
    return {
        serverGroups: [],
        currentServerGroup: {
            id: 0,
            name: '',
            status: 1
        },
        groupServers: [],
        availableServers: [],
        showPanel: false,
        isEditing: false,
        panelTitle: '',
        showManageServersPanel: false,
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

        openCreatePanel() {
            this.isEditing = false;
            this.panelTitle = '创建服务器组';
            this.currentServerGroup = {
                id: 0,
                name: '',
                status: 1
            };
            this.showPanel = true;
            this.showManageServersPanel = false;
        },

        closePanel() {
            this.showPanel = false;
            this.showManageServersPanel = false;
            this.isEditing = false;
        },

        async createServerGroup() {
            try {
                const response = await fetch('/api/admin/server_groups', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentServerGroup)
                });

                if (!response.ok) {
                    throw new Error('Failed to create server group');
                }

                this.closePanel();
                this.loadServerGroups();
                ShowMessage('创建服务器组成功');
            } catch (error) {
                console.error('Error creating server group:', error);
                ShowError('创建服务器组失败');
            }
        },

        async updateServerGroup() {
            try {
                const response = await fetch(`/api/admin/server_groups/${this.currentServerGroup.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentServerGroup)
                });

                if (!response.ok) {
                    throw new Error('Failed to update server group');
                }

                this.closePanel();
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

        editServerGroup(group) {
            this.isEditing = true;
            this.panelTitle = '编辑服务器组';
            this.currentServerGroup = { ...group };
            this.showPanel = true;
            this.showManageServersPanel = false;
        },

        async loadGroupServers(groupId) {
            try {
                const response = await fetch(`/api/admin/server_groups/${groupId}/servers`);
                if (!response.ok) {
                    throw new Error('Failed to load group servers');
                }
                const data = await response.json();
                this.groupServers = data.data.map(server => ({
                    ...server,
                    isEditing: false,
                    editData: {
                        server_id: server.server_id,
                        name: server.name,
                        server_status: server.server_status,
                        available: server.available,
                        merge_id: server.merge_id
                    }
                }));
            } catch (error) {
                console.error('Error loading group servers:', error);
                ShowError('加载服务器组中的服务器失败');
            }
        },

        startEditServer(server) {
            server.isEditing = true;
            server.editData = {
                physical_server_id: server.physical_server_id,
                server_id: server.server_id,
                name: server.name,
                server_status: server.server_status,
                available: server.available,
                merge_id: server.merge_id
            };
        },

        cancelEditServer(server) {
            server.isEditing = false;
        },

        async saveServerEdit(server) {
            try {
                const editData = {
                    physical_server_id: parseInt(server.editData.physical_server_id),
                    server_id: parseInt(server.editData.server_id),
                    name: server.editData.name,
                    server_status: parseInt(server.editData.server_status),
                    available: server.editData.available === 'true' || server.editData.available === true,
                    merge_id: parseInt(server.editData.merge_id)
                };

                const response = await fetch(`/api/admin/server_groups/${this.currentServerGroup.id}/servers/${server.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(editData)
                });

                if (!response.ok) {
                    throw new Error('Failed to update server');
                }

                // Update local data
                server.physical_server_id = editData.physical_server_id;
                server.server_id = editData.server_id;
                server.name = editData.name;
                server.server_status = editData.server_status;
                server.available = editData.available;
                server.merge_id = editData.merge_id;
                server.isEditing = false;
                ShowMessage('更新服务器成功');
            } catch (error) {
                console.error('Error updating server:', error);
                ShowError('更新服务器失败');
            }
        },

        async loadAvailableServers() {
            try {
                const response = await fetch('/api/admin/physical_servers');
                if (!response.ok) {
                    throw new Error('Failed to load available servers');
                }
                const data = await response.json();
                this.availableServers = data.servers;
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

        manageServers(group) {
            this.currentServerGroup = { ...group };
            this.showManageServersPanel = true;
            this.showPanel = false;
            this.loadGroupServers(group.id);
            this.loadAvailableServers();
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) {
                return;
            }
            this.currentPage = page;
            this.loadServerGroups();
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

        getServerStatusText(status) {
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