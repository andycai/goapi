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

        init() {
            this.loadConfig();
            this.loadCommits();
        },

        async loadConfig() {
            try {
                const response = await fetch('/api/reposync/config');
                if (!response.ok) throw new Error('加载配置失败');
                const data = await response.json();
                if (data) {
                    this.config = data;
                }
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async saveConfig() {
            try {
                const response = await fetch('/api/reposync/config', {
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

                Alpine.store('notification').show('配置保存成功', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async checkoutRepos() {
            try {
                const response = await fetch('/api/reposync/checkout', {
                    method: 'POST'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '检出仓库失败');
                }

                await this.loadCommits();
                Alpine.store('notification').show('仓库检出成功', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async loadCommits() {
            try {
                const response = await fetch('/api/reposync/commits');
                if (!response.ok) throw new Error('加载提交记录失败');
                this.commits = await response.json();
                this.selectedCommits = [];
                this.selectAll = false;
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
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
                Alpine.store('notification').show('请选择要同步的提交记录', 'error');
                return;
            }

            try {
                const response = await fetch('/api/reposync/sync', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ revisions: this.selectedCommits })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '批量同步失败');
                }

                await this.loadCommits();
                Alpine.store('notification').show('批量同步成功', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        async syncCommit(revision) {
            try {
                const response = await fetch('/api/reposync/sync', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ revision })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '同步失败');
                }

                await this.loadCommits();
                Alpine.store('notification').show('同步成功', 'success');
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
        }
    };
}