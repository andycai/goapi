function lubanManagement() {
    return {
        projects: [],
        currentProject: null,
        currentProjectTables: [],
        showProjectModal: false,
        showProjectDetailModal: false,
        showTableModal: false,
        showExportModal: false,
        editMode: false,
        form: {
            id: '',
            name: '',
            description: '',
            root_path: '',
            output_path: '',
            status: 'active'
        },
        tableForm: {
            id: '',
            project_id: '',
            name: '',
            description: '',
            file_path: '',
            file_type: 'excel',
            sheet_name: '',
            validators: '',
            status: 'active'
        },
        exportForm: {
            project_id: '',
            table_id: '',
            format: 'json',
            language: 'go'
        },
        exportProgress: null,
        progressTimer: null,

        init() {
            this.fetchProjects();
        },

        async fetchProjects() {
            try {
                const response = await fetch('/api/admin/luban/projects');
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
                root_path: '',
                output_path: '',
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
                const url = this.editMode ? `/api/admin/luban/projects/${this.form.id}` : '/api/admin/luban/projects';
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
                const response = await fetch(`/api/admin/luban/projects/${id}`, {
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
            await this.fetchProjectTables(project.id);
            this.showProjectDetailModal = true;
        },

        async fetchProjectTables(projectId) {
            try {
                const response = await fetch(`/api/admin/luban/tables?project_id=${projectId}`);
                if (!response.ok) throw new Error('获取配置表列表失败');
                this.currentProjectTables = await response.json();
            } catch (error) {
                ShowError(error.message);
            }
        },

        createTable() {
            this.editMode = false;
            this.tableForm = {
                id: '',
                project_id: this.currentProject.id,
                name: '',
                description: '',
                file_path: '',
                file_type: 'excel',
                sheet_name: '',
                validators: '',
                status: 'active'
            };
            this.showTableModal = true;
        },

        editTable(table) {
            this.editMode = true;
            this.tableForm = { ...table };
            this.showTableModal = true;
        },

        async submitTable() {
            try {
                const url = this.editMode 
                    ? `/api/admin/luban/tables/${this.tableForm.id}` 
                    : '/api/admin/luban/tables';
                const method = this.editMode ? 'PUT' : 'POST';
                
                const formData = { ...this.tableForm };
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

                await this.fetchProjectTables(this.currentProject.id);
                this.showTableModal = false;
                ShowMessage(this.editMode ? '配置表更新成功' : '配置表创建成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async deleteTable(id) {
            if (!confirm('确定要删除这个配置表吗？')) return;

            try {
                const response = await fetch(`/api/admin/luban/tables/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除配置表失败');

                await this.fetchProjectTables(this.currentProject.id);
                ShowMessage('配置表删除成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        exportTable(table) {
            this.exportForm = {
                project_id: this.currentProject.id,
                table_id: table.id,
                format: 'json',
                language: 'go'
            };
            this.exportProgress = null;
            this.showExportModal = true;
        },

        async submitExport() {
            try {
                const response = await fetch('/api/admin/luban/export', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.exportForm)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                const result = await response.json();
                this.startProgressPolling(result.id);
                ShowMessage('导出任务已开始');
            } catch (error) {
                ShowError(error.message);
            }
        },

        startProgressPolling(exportId) {
            if (this.progressTimer) {
                clearInterval(this.progressTimer);
            }

            const pollProgress = async () => {
                try {
                    const response = await fetch(`/api/admin/luban/exports/progress/${exportId}`);
                    if (!response.ok) throw new Error('获取导出进度失败');
                    
                    this.exportProgress = await response.json();
                    
                    if (this.exportProgress.status !== 'running') {
                        clearInterval(this.progressTimer);
                        this.progressTimer = null;

                        if (this.exportProgress.status === 'success') {
                            ShowMessage('导出成功');
                        } else if (this.exportProgress.status === 'failed') {
                            ShowError('导出失败: ' + this.exportProgress.error);
                        }
                    }
                } catch (error) {
                    clearInterval(this.progressTimer);
                    this.progressTimer = null;
                    ShowError(error.message);
                }
            };

            this.progressTimer = setInterval(pollProgress, 1000);
            pollProgress(); // 立即执行一次
        }
    };
} 