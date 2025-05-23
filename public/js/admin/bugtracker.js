function bugtrackerManagement() {
    return {
        projects: [],
        currentProject: null,
        currentProjectIssues: [],
        currentIssue: null,
        currentIssueComments: [],
        showProjectModal: false,
        showProjectDetailModal: false,
        showIssueModal: false,
        showIssueDetailModal: false,
        editMode: false,
        form: {
            id: '',
            name: '',
            description: '',
            status: 'active'
        },
        issueForm: {
            id: '',
            title: '',
            description: '',
            status: 'open',
            priority: 'medium'
        },
        commentForm: {
            content: ''
        },

        init() {
            this.fetchProjects();
        },

        async fetchProjects() {
            try {
                const response = await fetch('/api/admin/bugtracker/projects');
                if (!response.ok) throw new Error('获取项目列表失败');
                this.projects = await response.json();
            } catch (error) {
                ShowError(error.message);
            }
        },

        createProject() {
            this.editMode = false;
            this.form = {
                id: '',
                name: '',
                description: '',
                status: 'active'
            };
            this.showProjectModal = true;
        },

        editProject(project) {
            this.editMode = true;
            this.form = { ...project };
            this.showProjectModal = true;
        },

        async submitProject() {
            try {
                const url = this.editMode ? `/api/admin/bugtracker/projects/${this.form.id}` : '/api/admin/bugtracker/projects';
                const method = this.editMode ? 'PUT' : 'POST';
                
                const formData = { ...this.form };
                if (!this.editMode) {
                    delete formData.id;
                }
                
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                await this.fetchProjects();
                this.showProjectModal = false;
                ShowMessage(this.editMode ? '项目更新成功' : '项目创建成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async deleteProject(id) {
            if (!confirm('确定要删除这个项目吗？')) return;

            try {
                const response = await fetch(`/api/admin/bugtracker/projects/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除项目失败');

                await this.fetchProjects();
                ShowMessage('项目删除成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async viewProject(project) {
            this.currentProject = project;
            await this.fetchProjectIssues(project.id);
            this.showProjectDetailModal = true;
        },

        async fetchProjectIssues(projectId) {
            try {
                const response = await fetch(`/api/admin/bugtracker/projects/${projectId}/issues`);
                if (!response.ok) throw new Error('获取问题列表失败');
                this.currentProjectIssues = await response.json();
            } catch (error) {
                ShowError(error.message);
            }
        },

        createIssue() {
            this.editMode = false;
            this.issueForm = {
                id: '',
                title: '',
                description: '',
                status: 'open',
                priority: 'medium'
            };
            this.showIssueModal = true;
        },

        editIssue(issue) {
            this.editMode = true;
            this.issueForm = { ...issue };
            this.showIssueModal = true;
        },

        async submitIssue() {
            try {
                const url = this.editMode 
                    ? `/api/admin/bugtracker/projects/${this.currentProject.id}/issues/${this.issueForm.id}` 
                    : `/api/admin/bugtracker/projects/${this.currentProject.id}/issues`;
                const method = this.editMode ? 'PUT' : 'POST';
                
                const formData = { ...this.issueForm };
                if (!this.editMode) {
                    delete formData.id;
                }
                
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                await this.fetchProjectIssues(this.currentProject.id);
                this.showIssueModal = false;
                ShowMessage(this.editMode ? '问题更新成功' : '问题创建成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async deleteIssue(id) {
            if (!confirm('确定要删除这个问题吗？')) return;

            try {
                const response = await fetch(`/api/admin/bugtracker/projects/${this.currentProject.id}/issues/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除问题失败');

                await this.fetchProjectIssues(this.currentProject.id);
                ShowMessage('问题删除成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async viewIssue(issue) {
            this.currentIssue = issue;
            await this.fetchIssueComments(issue.id);
            this.showIssueDetailModal = true;
        },

        async fetchIssueComments(issueId) {
            try {
                const response = await fetch(`/api/admin/bugtracker/issues/${issueId}/comments`);
                if (!response.ok) throw new Error('获取评论列表失败');
                this.currentIssueComments = await response.json();
            } catch (error) {
                ShowError(error.message);
            }
        },

        async submitComment() {
            try {
                const response = await fetch(`/api/admin/bugtracker/issues/${this.currentIssue.id}/comments`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.commentForm)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                await this.fetchIssueComments(this.currentIssue.id);
                this.commentForm.content = '';
                ShowMessage('评论添加成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        getIssueStatusText(status) {
            const statusMap = {
                'open': '待处理',
                'in_progress': '处理中',
                'resolved': '已解决',
                'closed': '已关闭'
            };
            return statusMap[status] || status;
        },

        getIssuePriorityText(priority) {
            const priorityMap = {
                'low': '低',
                'medium': '中',
                'high': '高'
            };
            return priorityMap[priority] || priority;
        },

        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        }
    };
} 