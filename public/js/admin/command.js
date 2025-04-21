function commandManagement() {
    window.Alpine = window.Alpine || {};
    if (!Alpine.store('notification')) {
        Alpine.store('notification', {
            show: (message, type) => {
                console.error(message);
            },
            after: () => {}
        });
    }
    return {
        command: {
            name: '',
            script: ''
        },
        commands: [],
        executions: [],
        showExecutionModal: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 1,

        init() {
            this.loadCommands();
        },

        async loadCommands() {
            try {
                const response = await fetch(`/api/commands?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) throw new Error('加载命令列表失败');
                const data = await response.json();
                this.commands = data.commands;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            await this.loadCommands();
        },

        async createCommand() {
            if (!this.command.name || !this.command.script) {
                Alpine.store('notification').show('命令名称和脚本不能为空', 'error');
                return;
            }

            try {
                const response = await fetch('/api/commands', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.command)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '创建命令失败');
                }

                Alpine.store('notification').show('命令创建成功', 'success');
                this.command = { name: '', script: '' };
                await this.loadCommands();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async executeCommand(id) {
            try {
                const response = await fetch(`/api/commands/${id}/execute`, {
                    method: 'POST'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '执行命令失败');
                }

                Alpine.store('notification').show('命令执行成功', 'success');
                await this.loadCommands();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async viewExecutions(id) {
            try {
                const response = await fetch(`/api/commands/${id}/executions`);
                if (!response.ok) throw new Error('获取执行记录失败');
                const data = await response.json();
                this.executions = data;
                this.showExecutionModal = true;
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
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
        },

        getStatusText(status) {
            switch (status) {
                case 0: return '等待中';
                case 1: return '成功';
                case 2: return '失败';
                default: return '未知';
            }
        },

        getStatusClass(status) {
            switch (status) {
                case 0: return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-700 dark:text-yellow-100';
                case 1: return 'bg-green-100 text-green-800 dark:bg-green-700 dark:text-green-100';
                case 2: return 'bg-red-100 text-red-800 dark:bg-red-700 dark:text-red-100';
                default: return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-100';
            }
        }
    };
} 