// Bases management functionality
function basesManagement() {
    return {
        // 状态变量
        entities: [],
        fields: [],
        entityData: [],
        currentPage: 1,
        pageSize: 10,
        totalPages: 1,
        totalRecords: 0,
        searchKeyword: '',
        loading: false,
        editMode: false,
        currentEntityId: 0,

        // 面板状态
        showPanel: false,
        showFieldPanel: false,
        showFieldForm: false,
        showDataPanel: false,
        showDataForm: false,

        // 表单数据
        form: {
            name: '',
            table_name: '',
            description: ''
        },
        fieldForm: {
            name: '',
            type: '',
            length: '',
            is_nullable: false,
            is_unique: false,
            default: '',
            description: ''
        },
        dataForm: {},

        // 字段类型选项
        fieldTypes: [
            'string',
            'text',
            'int',
            'float',
            'bool',
            'datetime',
            'date',
            'time',
            'json'
        ],

        // 初始化
        init() {
            this.loadEntities();
        },

        // 加载实体列表
        async loadEntities() {
            try {
                this.loading = true;
                let url = `/api/admin/bases/entities?page=${this.currentPage}&limit=${this.pageSize}`;
                if (this.searchKeyword) {
                    url += `&search=${encodeURIComponent(this.searchKeyword)}`;
                }
                
                const response = await fetch(url);
                if (!response.ok) throw new Error('获取实体列表失败');
                const data = await response.json();
                
                this.entities = data.entities;
                this.totalRecords = data.total;
                this.totalPages = Math.ceil(data.total / this.pageSize);
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 搜索
        search() {
            this.currentPage = 1;
            this.loadEntities();
        },

        // 分页
        goToPage(page) {
            if (page < 1 || page > this.totalPages || page === this.currentPage) return;
            this.currentPage = page;
            this.loadEntities();
        },

        changePage(page) {
            if (page < 1 || page > this.totalPages) return;
            this.currentPage = page;
            this.loadEntities();
        },

        // 计算分页按钮
        get paginationPages() {
            const pages = [];
            const maxPages = 5;
            let start = Math.max(1, this.currentPage - Math.floor(maxPages / 2));
            let end = Math.min(this.totalPages, start + maxPages - 1);
            if (end - start + 1 < maxPages) {
                start = Math.max(1, end - maxPages + 1);
            }
            for (let i = start; i <= end; i++) {
                pages.push(i);
            }
            return pages;
        },

        // 创建实体
        createEntity() {
            this.editMode = false;
            this.form = {
                name: '',
                table_name: '',
                description: ''
            };
            this.showPanel = true;
        },

        // 编辑实体
        editEntity(entity) {
            this.editMode = true;
            this.form = { ...entity };
            this.showPanel = true;
        },

        // 删除实体
        async deleteEntity(id) {
            if (!confirm('确定要删除这个实体吗？此操作不可恢复。')) return;
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/bases/entities/${id}`, {
                    method: 'DELETE'
                });
                if (!response.ok) throw new Error('删除实体失败');
                
                ShowMessage('实体删除成功');
                this.loadEntities();
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 提交实体表单
        async submitForm() {
            if (this.loading) return;
            
            // 验证表单
            if (!this.validateEntityForm()) return;
            
            try {
                this.loading = true;
                const url = this.editMode ? `/api/admin/bases/entities/${this.form.id}` : '/api/admin/bases/entities';
                const method = this.editMode ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.form)
                });
                if (!response.ok) throw new Error('保存实体失败');
                
                ShowMessage(this.editMode ? '实体更新成功' : '实体创建成功');
                this.closePanel();
                this.loadEntities();
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 验证实体表单
        validateEntityForm() {
            if (!this.form.name.trim()) {
                ShowError('实体名称不能为空');
                return false;
            }
            if (!this.form.table_name.trim()) {
                ShowError('表名不能为空');
                return false;
            }
            return true;
        },

        // 关闭面板
        closePanel() {
            this.showPanel = false;
            this.form = {
                name: '',
                table_name: '',
                description: ''
            };
        },

        // 打开字段面板
        async openFieldPanel(entity) {
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/bases/fields?entity_id=${entity.id}`);
                if (!response.ok) throw new Error('获取字段列表失败');
                const data = await response.json();
                
                this.fields = data.fields;
                this.showFieldPanel = true;
                this.currentEntityId = entity.id;
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 关闭字段面板
        closeFieldPanel() {
            this.showFieldPanel = false;
            this.fields = [];
        },

        // 创建字段
        createField() {
            this.editMode = false;
            this.fieldForm = {
                name: '',
                type: 'string',
                length: 0,
                is_nullable: false,
                is_unique: false,
                default: '',
                description: ''
            };
            this.showFieldForm = true;
        },

        // 编辑字段
        editField(field) {
            this.editMode = true;
            this.fieldForm = { ...field };
            this.showFieldForm = true;
        },

        // 删除字段
        async deleteField(id) {
            if (!confirm('确定要删除这个字段吗？此操作不可恢复。')) return;
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/bases/fields/${id}`, {
                    method: 'DELETE'
                });
                if (!response.ok) throw new Error('删除字段失败');
                
                ShowMessage('字段删除成功');
                this.openFieldPanel({ id: this.currentEntityId });
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 提交字段表单
        async submitFieldForm() {
            if (this.loading) return;
            
            // 验证表单
            if (!this.validateFieldForm()) return;

            const fieldData = {
                ...this.fieldForm,
                id: parseInt(this.fieldForm.id) || 0,
                entity_id: parseInt(this.currentEntityId) || 0,
                length: parseInt(this.fieldForm.length) || 0,
                is_nullable: this.fieldForm.is_nullable === 'true' || this.fieldForm.is_nullable === true,
            }
            
            try {
                this.loading = true;
                const url = this.editMode ? `/api/admin/bases/fields/${this.fieldForm.id}` : '/api/admin/bases/fields';
                const method = this.editMode ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(fieldData)
                });
                if (!response.ok) throw new Error('保存字段失败');
                
                ShowMessage(this.editMode ? '字段更新成功' : '字段创建成功');
                this.closeFieldForm();
                // this.openFieldPanel({ id: this.currentEntityId });
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 验证字段表单
        validateFieldForm() {
            if (!this.fieldForm.name.trim()) {
                ShowError('字段名称不能为空');
                return false;
            }
            if (!this.fieldForm.type) {
                ShowError('字段类型不能为空');
                return false;
            }
            return true;
        },

        // 关闭字段表单
        closeFieldForm() {
            this.showFieldForm = false;
            this.fieldForm = {
                name: '',
                type: '',
                length: '',
                is_nullable: false,
                is_unique: false,
                default: '',
                description: ''
            };
            this.openFieldPanel({ id: this.currentEntityId });
        },

        // 打开数据面板
        async openDataPanel(entity) {
            try {
                this.loading = true;
                const [fieldsResponse, dataResponse] = await Promise.all([
                    fetch(`/api/admin/bases/fields?entity_id=${entity.id}`),
                    fetch(`/api/admin/bases/data?entity_id=${entity.id}`)
                ]);
                if (!fieldsResponse.ok) throw new Error('获取字段列表失败');
                if (!dataResponse.ok) throw new Error('获取数据列表失败');
                
                const [fieldsData, dataData] = await Promise.all([
                    fieldsResponse.json(),
                    dataResponse.json()
                ]);
                
                this.fields = fieldsData.fields;
                this.entityData = dataData.data;
                this.showDataPanel = true;
                this.currentEntityId = entity.id;
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 关闭数据面板
        closeDataPanel() {
            this.showDataPanel = false;
            // this.fields = [];
            this.entityData = [];
        },

        // 创建数据
        createData() {
            this.editMode = false;
            this.dataForm = {};
            this.fields.forEach(field => {
                this.dataForm[field.name] = field.default || '';
            });
            this.showDataForm = true;
        },

        // 编辑数据
        editData(data) {
            this.editMode = true;
            this.dataForm = { ...data.data };
            this.showDataForm = true;
        },

        // 删除数据
        async deleteData(id) {
            if (!confirm('确定要删除这条数据吗？此操作不可恢复。')) return;
            try {
                this.loading = true;
                const response = await fetch(`/api/admin/bases/data/${id}`, {
                    method: 'DELETE'
                });
                if (!response.ok) throw new Error('删除数据失败');
                
                ShowMessage('数据删除成功');
                this.openDataPanel({ id: this.currentEntityId });
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 提交数据表单
        async submitDataForm() {
            if (this.loading) return;
            
            // 验证表单
            if (!this.validateDataForm()) return;

            // 根据字段类型转换数据
            const convertedData = {};
            for (const field of this.fields) {
                let value = this.dataForm[field.name];
                
                // 如果值为空且字段可为空，则跳过
                if (value === '' && field.is_nullable) {
                    continue;
                }

                // 根据字段类型转换值
                switch (field.type) {
                    case 'int':
                        value = value === '' ? 0 : parseInt(value);
                        break;
                    case 'float':
                        value = value === '' ? 0 : parseFloat(value);
                        break;
                    case 'bool':
                        value = value === 'true' || value === true;
                        break;
                    case 'json':
                        try {
                            value = value === '' ? {} : JSON.parse(value);
                        } catch (e) {
                            ShowError(`字段 ${field.name} 的 JSON 格式不正确`);
                            return;
                        }
                        break;
                    case 'datetime':
                    case 'date':
                    case 'time':
                        // 保持日期时间格式不变，由服务器处理
                        break;
                    default:
                        // string, text 等类型保持原样
                        break;
                }
                
                convertedData[field.name] = value;
            }
            
            const fieldDataForm = {
                data: convertedData,
                entity_id: parseInt(this.currentEntityId) || 0,
            }
            
            try {
                this.loading = true;
                const url = this.editMode ? `/api/admin/bases/data/${this.dataForm.id}` : '/api/admin/bases/data';
                const method = this.editMode ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(fieldDataForm)
                });
                if (!response.ok) throw new Error('保存数据失败');
                
                ShowMessage(this.editMode ? '数据更新成功' : '数据创建成功');
                this.closeDataForm();
                // this.openDataPanel({ id: this.currentEntityId });
            } catch (error) {
                ShowError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // 验证数据表单
        validateDataForm() {
            for (const field of this.fields) {
                if (!field.is_nullable && !this.dataForm[field.name]) {
                    ShowError(`字段 ${field.name} 不能为空`);
                    return false;
                }
            }
            return true;
        },

        // 关闭数据表单
        closeDataForm() {
            this.showDataForm = false;
            this.dataForm = {};
            this.openDataPanel({ id: this.currentEntityId });
        },

        // 格式化日期
        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
            });
        }
    };
} 