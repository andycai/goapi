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
                Alpine.store('notification').show(error.message, 'error');
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
                Alpine.store('notification').show(this.modalAction === 'add' ? '添加页面成功' : '更新页面成功', 'success');
                this.closePanel();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
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
                Alpine.store('notification').show('删除页面成功', 'success');
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