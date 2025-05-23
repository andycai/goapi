function reposyncManagement() {
    return {
        config: {
            repo_type1: 'svn',
            repo_url1: '',
            local_path1: '',
            username1: '',
            password1: '',
            repo_type2: 'svn',
            repo_url2: '',
            local_path2: '',
            username2: '',
            password2: ''
        },
        commits: [],
        selectedCommits: [],
        selectAll: false,
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 1,
        refreshLimit: 100,
        openFileLists: new Set(), // 跟踪哪些提交的文件列表是展开的

        init() {
            this.loadConfig();
            this.loadCommits();
        },

        async loadConfig() {
            try {
                const response = await fetch('/api/admin/reposync/config');
                if (!response.ok) throw new Error('加载配置失败');
                const data = await response.json();
                if (data) {
                    this.config = data;
                }
            } catch (error) {
                ShowError(error.message);
            }
        },

        async saveConfig() {
            try {
                const response = await fetch('/api/admin/reposync/config', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.config)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '保存配置失败');
                }

                ShowMessage('配置保存成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async checkoutRepos() {
            try {
                const response = await fetch('/api/admin/reposync/checkout', {
                    method: 'POST'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '检出仓库失败');
                }

                await this.loadCommits();
                ShowMessage('仓库检出成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async loadCommits() {
            try {
                const response = await fetch(`/api/admin/reposync/commits?page=${this.currentPage}&pageSize=${this.pageSize}`);
                if (!response.ok) throw new Error('加载提交记录失败');
                const data = await response.json();
                this.commits = data.commits;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
                this.selectedCommits = [];
                this.selectAll = false;
                this.openFileLists = new Set(); // 重置打开的文件列表
            } catch (error) {
                ShowError(error.message);
            }
        },

        async changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            await this.loadCommits();
        },

        toggleSelectAll() {
            this.selectAll = !this.selectAll;
            this.selectedCommits = this.selectAll ? this.commits.filter(commit => !commit.synced).map(commit => commit.revision) : [];
        },

        toggleSelect(revision) {
            const index = this.selectedCommits.indexOf(revision);
            if (index === -1) {
                this.selectedCommits.push(revision);
            } else {
                this.selectedCommits.splice(index, 1);
            }
            this.selectAll = this.selectedCommits.length === this.commits.filter(commit => !commit.synced).length;
        },

        async syncSelectedCommits() {
            if (this.selectedCommits.length === 0) {
                ShowError('请选择要同步的提交记录');
                return;
            }

            try {
                // 对版本号进行排序（从小到大）
                const sortedRevisions = [...this.selectedCommits].sort((a, b) => {
                    // 如果是数字版本号，按数字大小排序
                    const numA = parseInt(a);
                    const numB = parseInt(b);
                    if (!isNaN(numA) && !isNaN(numB)) {
                        return numA - numB;
                    }
                    // 如果不是数字，按字符串排序
                    return a.localeCompare(b);
                });

                const response = await fetch('/api/admin/reposync/sync', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ revisions: sortedRevisions })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '批量同步失败');
                }

                await this.loadCommits();
                ShowMessage('批量同步成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async syncCommit(revision) {
            try {
                const response = await fetch('/api/admin/reposync/sync', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ revisions: [revision] })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '同步失败');
                }

                await this.loadCommits();
                ShowMessage('同步成功');
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
        },

        getChangeTypeText(changeType) {
            const types = {
                'A': '新增',
                'M': '修改',
                'D': '删除'
            };
            return types[changeType] || changeType;
        },
        
        // 根据类型获取文件列表
        getFilesByType(files, type) {
            if (!files) return [];
            return files.filter(file => file.change_type === type);
        },
        
        // 切换文件列表的显示状态
        toggleFileList(revision) {
            if (this.openFileLists.has(revision)) {
                this.openFileLists.delete(revision);
            } else {
                this.openFileLists.add(revision);
            }
        },
        
        // 检查文件列表是否打开
        isFileListOpen(revision) {
            return this.openFileLists.has(revision);
        },

        async refreshCommits() {
            try {
                const response = await fetch('/api/admin/reposync/refresh', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ limit: this.refreshLimit })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '刷新提交记录失败');
                }

                await this.loadCommits();
                ShowMessage('刷新提交记录成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async clearSyncData() {
            if (!confirm('确定要清空所有同步数据吗？此操作不可恢复！')) {
                return;
            }

            try {
                const response = await fetch('/api/admin/reposync/clear', {
                    method: 'POST'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '清空数据失败');
                }

                await this.loadCommits();
                ShowMessage('清空数据成功');
            } catch (error) {
                ShowError(error.message);
            }
        }
    };
}