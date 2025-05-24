function unitoolManagement() {
    return {
        searchParams: {
            targetPath: '',
            notificationUrl: ''
        },
        logs: [],
        duplicateGuids: [],
        showDetailModal: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 1,

        init() {
            this.loadLogs();
        },

        async loadLogs() {
            try {
                const response = await fetch(`/api/admin/unitool/logs?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) throw new Error('加载日志列表失败');
                const data = await response.json();
                this.logs = data.logs;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                ShowError(error.message);
            }
        },

        async changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            await this.loadLogs();
        },

        async findDuplicateGuids() {
            if (!this.searchParams.targetPath || !this.searchParams.notificationUrl) {
                ShowError('目标路径和通知URL不能为空');
                return;
            }

            try {
                const response = await fetch('/unitool/findguid', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        targetPath: this.searchParams.targetPath,
                        notificationUrl: this.searchParams.notificationUrl
                    })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '查找重复GUID失败');
                }

                const data = await response.json();
                ShowMessage('查找任务已开始');
                
                // 一段时间后刷新列表
                setTimeout(() => this.loadLogs(), 2000);
            } catch (error) {
                ShowError(error.message);
            }
        },

        async viewDuplicateDetails(logId) {
            try {
                const response = await fetch(`/api/admin/unitool/duplicates/${logId}`);
                if (!response.ok) throw new Error('获取重复GUID详情失败');
                const data = await response.json();
                this.duplicateGuids = data.duplicates;
                this.showDetailModal = true;
            } catch (error) {
                ShowError(error.message);
            }
        },

        formatDate(timestamp) {
            if (!timestamp) return '';
            return new Date(timestamp).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        }
    };
} 