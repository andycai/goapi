function announcementManagement() {
    return {
        announcements: [],
        currentAnnouncement: {
            id: 0,
            title: '',
            content: '',
            status: 1
        },
        showPanel: false,
        isEditing: false,
        panelTitle: '',
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 0,

        init() {
            this.loadAnnouncements();
        },

        async loadAnnouncements() {
            try {
                const response = await fetch(`/api/admin/announcements?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) {
                    throw new Error('Failed to load announcements');
                }
                const data = await response.json();
                this.announcements = data.announcements;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                console.error('Error loading announcements:', error);
                alert('加载公告列表失败');
            }
        },

        openCreatePanel() {
            this.isEditing = false;
            this.panelTitle = '创建公告';
            this.currentAnnouncement = {
                id: 0,
                title: '',
                content: '',
                status: 1
            };
            this.showPanel = true;
        },

        closePanel() {
            this.showPanel = false;
            this.isEditing = false;
        },

        async createAnnouncement() {
            try {
                const response = await fetch('/api/admin/announcements', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentAnnouncement)
                });

                if (!response.ok) {
                    throw new Error('Failed to create announcement');
                }

                this.closePanel();
                this.loadAnnouncements();
                alert('创建公告成功');
            } catch (error) {
                console.error('Error creating announcement:', error);
                alert('创建公告失败');
            }
        },

        async updateAnnouncement() {
            try {
                const response = await fetch(`/api/admin/announcements/${this.currentAnnouncement.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentAnnouncement)
                });

                if (!response.ok) {
                    throw new Error('Failed to update announcement');
                }

                this.closePanel();
                this.loadAnnouncements();
                alert('更新公告成功');
            } catch (error) {
                console.error('Error updating announcement:', error);
                alert('更新公告失败');
            }
        },

        async deleteAnnouncement(id) {
            if (!confirm('确定要删除这条公告吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/announcements/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    throw new Error('Failed to delete announcement');
                }

                this.loadAnnouncements();
                alert('删除公告成功');
            } catch (error) {
                console.error('Error deleting announcement:', error);
                alert('删除公告失败');
            }
        },

        editAnnouncement(announcement) {
            this.isEditing = true;
            this.panelTitle = '编辑公告';
            this.currentAnnouncement = { ...announcement };
            this.showPanel = true;
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) {
                return;
            }
            this.currentPage = page;
            this.loadAnnouncements();
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