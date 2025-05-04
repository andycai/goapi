function channelManagement() {
    return {
        channels: [],
        channel: {
            name: '',
            serverList: '',
            cdnVersion: '',
            cdnUrl: '',
            cdnUrl2: '',
            openPatch: '',
            loginApi: '',
            loginUrl: '',
            pkgVersion: '',
            serverListUrl: '',
            noticeUrl: '',
            noticeNumUrl: ''
        },
        editingChannel: {
            id: 0,
            name: '',
            serverList: '',
            cdnVersion: '',
            cdnUrl: '',
            cdnUrl2: '',
            openPatch: '',
            loginApi: '',
            loginUrl: '',
            pkgVersion: '',
            serverListUrl: '',
            noticeUrl: '',
            noticeNumUrl: ''
        },
        showEditModal: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,

        init() {
            this.loadChannels();
        },

        async loadChannels() {
            try {
                const response = await fetch(`/api/channel/list?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) {
                    throw new Error('Failed to load channels');
                }
                const data = await response.json();
                this.channels = data.channels;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading channels:', error);
                alert('加载渠道列表失败');
            }
        },

        async createChannel() {
            try {
                const response = await fetch('/api/channel', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.channel)
                });

                if (!response.ok) {
                    throw new Error('Failed to create channel');
                }

                this.resetChannelForm();
                this.loadChannels();
                alert('创建渠道成功');
            } catch (error) {
                console.error('Error creating channel:', error);
                alert('创建渠道失败');
            }
        },

        async updateChannel() {
            try {
                const response = await fetch(`/api/channel/${this.editingChannel.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.editingChannel)
                });

                if (!response.ok) {
                    throw new Error('Failed to update channel');
                }

                this.showEditModal = false;
                this.loadChannels();
                alert('更新渠道成功');
            } catch (error) {
                console.error('Error updating channel:', error);
                alert('更新渠道失败');
            }
        },

        async deleteChannel(id) {
            if (!confirm('确定要删除这个渠道吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/channel/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    throw new Error('Failed to delete channel');
                }

                this.loadChannels();
                alert('删除渠道成功');
            } catch (error) {
                console.error('Error deleting channel:', error);
                alert('删除渠道失败');
            }
        },

        editChannel(channel) {
            this.editingChannel = { ...channel };
            this.showEditModal = true;
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) {
                return;
            }
            this.currentPage = page;
            this.loadChannels();
        },

        resetChannelForm() {
            this.channel = {
                name: '',
                serverList: '',
                cdnVersion: '',
                cdnUrl: '',
                cdnUrl2: '',
                openPatch: '',
                loginApi: '',
                loginUrl: '',
                pkgVersion: '',
                serverListUrl: '',
                noticeUrl: '',
                noticeNumUrl: ''
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
        }
    };
} 