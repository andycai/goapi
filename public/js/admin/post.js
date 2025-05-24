/**
 * 文章管理
 */
function postManagement() {
    return {
        // 文章数据
        posts: [],
        // 分页信息
        currentPage: 1,
        pageSize: 10,
        totalRecords: 0,
        pages: 1,
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
        modalAction: 'add',
        // 加载状态
        loading: false,

        init() {
            this.loadPosts();
        },

        // 加载文章列表
        async loadPosts() {
            try {
                const response = await fetch(`/api/admin/post/list?page=${this.currentPage}&limit=${this.pageSize}`);
                if (!response.ok) throw new Error('加载文章列表失败');
                const data = await response.json();
                this.posts = data.posts;
                this.totalRecords = data.total;
                this.pages = Math.ceil(this.totalRecords / this.pageSize);
            } catch (error) {
                ShowError(error.message);
            }
        },

        // 切换页码
        async changePage(page) {
            if (page < 1 || page > this.pages) return;
            this.currentPage = page;
            await this.loadPosts();
        },

        // 打开添加文章面板
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

        // 打开编辑文章面板
        openEditModal(post) {
            this.form = {
                id: post.id,
                title: post.title,
                content: post.content,
                slug: post.slug,
                status: post.status,
                author_id: post.author_id
            };
            this.modalAction = 'edit';
            this.showPanel = true;
        },

        // 关闭文章面板
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

        // 保存文章
        async savePost() {
            if (this.loading) return;
            this.loading = true;
            
            try {
                const url = this.modalAction === 'add' ? '/api/admin/post/add' : '/api/admin/post/edit';
                
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
                    throw new Error(error.error || '保存文章失败');
                }

                await this.loadPosts();
                ShowMessage(this.modalAction === 'add' ? '文章添加成功' : '文章更新成功');
                this.closePanel();
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 删除文章
        async deletePost(id) {
            if (!confirm('确定要删除这篇文章吗？')) {
                return;
            }

            try {
                const response = await fetch('/api/admin/post/delete', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ id: parseInt(id, 10) })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除文章失败');
                }

                await this.loadPosts();
                ShowMessage('文章删除成功');
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