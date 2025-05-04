function physicalServerManagement() {
    return {
        physicalServers: [],
        physicalServer: {
            serverId: '',
            serverName: '',
            status: 0,
            available: true,
            mergeId: '',
            onlineCount: 0,
            serverPort: 0,
            serverIp: ''
        },
        editingPhysicalServer: {
            id: 0,
            serverId: '',
            serverName: '',
            status: 0,
            available: true,
            mergeId: '',
            onlineCount: 0,
            serverPort: 0,
            serverIp: ''
        },
        showEditModal: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,

        init() {
            this.loadPhysicalServers();
        },

        async loadPhysicalServers() {
            try {
                const response = await fetch(`/api/physical-servers?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) {
                    throw new Error('Failed to load physical servers');
                }
                const data = await response.json();
                this.physicalServers = data.data;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading physical servers:', error);
                alert('加载物理服务器列表失败');
            }
        },

        async createPhysicalServer() {
            try {
                const response = await fetch('/api/physical-servers', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.physicalServer)
                });

                if (!response.ok) {
                    throw new Error('Failed to create physical server');
                }

                this.resetPhysicalServerForm();
                this.loadPhysicalServers();
                alert('创建物理服务器成功');
            } catch (error) {
                console.error('Error creating physical server:', error);
                alert('创建物理服务器失败');
            }
        },

        async updatePhysicalServer() {
            try {
                const response = await fetch(`/api/physical-servers/${this.editingPhysicalServer.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.editingPhysicalServer)
                });

                if (!response.ok) {
                    throw new Error('Failed to update physical server');
                }

                this.showEditModal = false;
                this.loadPhysicalServers();
                alert('更新物理服务器成功');
            } catch (error) {
                console.error('Error updating physical server:', error);
                alert('更新物理服务器失败');
            }
        },

        async deletePhysicalServer(id) {
            if (!confirm('确定要删除这个物理服务器吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/physical-servers/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    throw new Error('Failed to delete physical server');
                }

                this.loadPhysicalServers();
                alert('删除物理服务器成功');
            } catch (error) {
                console.error('Error deleting physical server:', error);
                alert('删除物理服务器失败');
            }
        },

        editPhysicalServer(server) {
            this.editingPhysicalServer = { ...server };
            this.showEditModal = true;
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) {
                return;
            }
            this.currentPage = page;
            this.loadPhysicalServers();
        },

        resetPhysicalServerForm() {
            this.physicalServer = {
                serverId: '',
                serverName: '',
                status: 0,
                available: true,
                mergeId: '',
                onlineCount: 0,
                serverPort: 0,
                serverIp: ''
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