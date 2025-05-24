/**
 * 页面管理
 */
function pageManagement() {
    return {
        // 页面数据
        pages: [],
        // 分页信息
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        totalPages: 1,
        // 表单数据
        form: {
            id: 0,
            title: '',
            content: '',
            slug: '',
            status: 'draft',
            author_id: 0
        },
        // 面板状态
        showPanel: false,
        modalAction: '',
        // 加载状态
        loading: false,

        init() {
            this.loadPages();
        },

        // 加载页面列表
        async loadPages() {
            try {
                const response = await fetch(`/api/admin/page/list?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) throw new Error('加载页面列表失败');
                const data = await response.json();
                this.pages = data.pages;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                ShowError(error.message);
            }
        },

        // 切换页码
        async changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            await this.loadPages();
        },

        // 打开添加页面面板
        openAddModal() {
            this.form = {
                id: 0,
                title: '',
                content: '',
                slug: '',
                status: 'draft',
                author_id: 0
            };
            this.modalAction = 'add';
            this.showPanel = true;
        },

        // 打开编辑页面面板
        openEditModal(page) {
            this.form = {
                id: page.id,
                title: page.title,
                content: page.content,
                slug: page.slug,
                status: page.status,
                author_id: page.author_id
            };
            this.modalAction = 'edit';
            this.showPanel = true;
        },

        // 关闭页面面板
        closePanel() {
            this.showPanel = false;
            this.form = {
                id: 0,
                title: '',
                content: '',
                slug: '',
                status: 'draft',
                author_id: 0
            };
        },

        // 保存页面
        async savePage() {
            if (this.loading) return;
            this.loading = true;
            
            try {
                const url = this.modalAction === 'add' ? '/api/admin/page/add' : '/api/admin/page/edit';
                
                // 确保id是整数
                const formData = {
                    ...this.form,
                    id: parseInt(this.form.id, 10) || 0,
                    author_id: parseInt(this.form.author_id, 10) || 0
                };
                
                const response = await fetch(url, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '保存页面失败');
                }

                await this.loadPages();
                ShowMessage(this.modalAction === 'add' ? '添加页面成功' : '更新页面成功');
                this.closePanel();
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 删除页面
        async deletePage(id) {
            if (!confirm('确定要删除这个页面吗？')) {
                return;
            }

            try {
                const response = await fetch('/api/admin/page/delete', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ id: parseInt(id, 10) })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除页面失败');
                }

                await this.loadPages();
                ShowMessage('页面删除成功');
            } catch (error) {
                ShowError(error.message);
            }
        },
        
        // 格式化日期
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

        // 处理Tab键缩进
        handleTab(e) {
            const textarea = e.target;
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;
            
            // 在当前光标位置插入Tab
            const text = textarea.value;
            textarea.value = text.substring(0, start) + '    ' + text.substring(end);
            
            // 重新设置光标位置
            textarea.selectionStart = textarea.selectionEnd = start + 4;
            
            // 更新表单内容
            this.form.content = textarea.value;
        },

        // 插入Markdown标记
        insertMarkdown(template) {
            const textarea = this.$refs.markdownEditor;
            if (!textarea) return;
            
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;
            const text = textarea.value;
            
            // 如果有选中文本，使用选中的文本替换模板中的占位符
            let insertion = template;
            if (start !== end) {
                const selection = text.substring(start, end);
                if (template === '**粗体**') {
                    insertion = `**${selection}**`;
                } else if (template === '*斜体*') {
                    insertion = `*${selection}*`;
                } else if (template === '# 标题') {
                    insertion = `# ${selection}`;
                } else if (template === '[链接文本](https://example.com)') {
                    insertion = `[${selection}](https://example.com)`;
                } else if (template === '![图片描述](https://example.com/image.jpg)') {
                    insertion = `![${selection}](https://example.com/image.jpg)`;
                } else if (template === '```\n代码块\n```') {
                    insertion = `\`\`\`\n${selection}\n\`\`\``;
                }
            }
            
            // 插入内容
            textarea.value = text.substring(0, start) + insertion + text.substring(end);
            
            // 更新表单内容
            this.form.content = textarea.value;
            
            // 设置新的光标位置
            const newPosition = start + insertion.length;
            textarea.selectionStart = textarea.selectionEnd = newPosition;
            
            // 聚焦回文本框
            textarea.focus();
        },
    };
} 