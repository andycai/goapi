function physicalServerManagement() {
    return {
        physicalServers: [],
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,
        showPanel: false,
        isEditing: false,
        panelTitle: '',
        currentServer: {
            id: 0,
            serverId: 0,
            name: '',
            serverStatus: 1,
            available: true,
            mergeId: 0,
            online: 0,
            serverPort: 0,
            serverIp: ''
        },

        init() {
            this.loadPhysicalServers();
        },

        async loadPhysicalServers() {
            try {
                const response = await fetch(`/api/admin/physical_servers?page=${this.currentPage}&pageSize=${this.pageSize}`);
                const data = await response.json();
                this.physicalServers = data.servers;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading physical servers:', error);
            }
        },

        changePage(page) {
            if (page >= 1 && page <= this.totalPages) {
                this.currentPage = page;
                this.loadPhysicalServers();
            }
        },

        getServerStatusText(status) {
            const statusMap = {
                0: '维护中',
                1: '正常',
                2: '爆满'
            };
            return statusMap[status] || '未知';
        },

        openCreatePanel() {
            this.isEditing = false;
            this.panelTitle = '创建物理服务器';
            this.currentServer = {
                id: 0,
                serverId: 0,
                name: '',
                serverStatus: 1,
                available: true,
                mergeId: 0,
                online: 0,
                serverPort: 0,
                serverIp: ''
            };
            this.showPanel = true;
        },

        editPhysicalServer(server) {
            this.isEditing = true;
            this.panelTitle = '编辑物理服务器';
            this.currentServer = { ...server };
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
        },

        async createPhysicalServer() {
            try {
                // Convert string values to numbers where needed
                const serverData = {
                    ...this.currentServer,
                    id: parseInt(this.currentServer.id) || 0,
                    serverId: parseInt(this.currentServer.serverId) || 0,
                    serverStatus: parseInt(this.currentServer.serverStatus) || 1,
                    mergeId: parseInt(this.currentServer.mergeId) || 0,
                    online: parseInt(this.currentServer.online) || 0,
                    serverPort: parseInt(this.currentServer.serverPort) || 0,
                    available: this.currentServer.available === 'true' || this.currentServer.available === true
                };

                const response = await fetch('/api/admin/physical_servers', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(serverData)
                });

                if (response.ok) {
                    this.closePanel();
                    this.loadPhysicalServers();
                } else {
                    console.error('Failed to create physical server');
                }
            } catch (error) {
                console.error('Error creating physical server:', error);
            }
        },

        async updatePhysicalServer() {
            try {
                // Convert string values to numbers where needed
                const serverData = {
                    ...this.currentServer,
                    serverId: parseInt(this.currentServer.serverId) || 0,
                    serverStatus: parseInt(this.currentServer.serverStatus) || 1,
                    mergeId: parseInt(this.currentServer.mergeId) || 0,
                    online: parseInt(this.currentServer.online) || 0,
                    serverPort: parseInt(this.currentServer.serverPort) || 0,
                    available: this.currentServer.available === 'true' || this.currentServer.available === true
                };

                const response = await fetch(`/api/admin/physical_servers/${this.currentServer.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(serverData)
                });

                if (response.ok) {
                    this.closePanel();
                    this.loadPhysicalServers();
                } else {
                    console.error('Failed to update physical server');
                }
            } catch (error) {
                console.error('Error updating physical server:', error);
            }
        },

        async deletePhysicalServer(id) {
            if (confirm('确定要删除这个物理服务器吗？')) {
                try {
                    const response = await fetch(`/api/admin/physical_servers/${id}`, {
                        method: 'DELETE'
                    });

                    if (response.ok) {
                        this.loadPhysicalServers();
                    } else {
                        console.error('Failed to delete physical server');
                    }
                } catch (error) {
                    console.error('Error deleting physical server:', error);
                }
            }
        }
    }
} 