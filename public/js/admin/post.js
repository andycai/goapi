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
                Alpine.store('notification').show(error.message, 'error');
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
                Alpine.store('notification').show(this.modalAction === 'add' ? '添加文章成功' : '更新文章成功', 'success');
                this.closePanel();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
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
                Alpine.store('notification').show('删除文章成功', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
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
        }
    };
} 