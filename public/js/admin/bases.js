// bases 模块
function basesManagement() {
    return {
        // 状态
        entities: [],
        fields: [],
        entityData: [],
        showEntityPanel: false,
        showFieldPanel: false,
        showFieldForm: false,
        showDataPanel: false,
        showDataForm: false,
        currentEntity: {},
        currentField: {},
        currentData: {},
        fieldTypes: ['string', 'text', 'int', 'float', 'bool', 'datetime', 'date', 'time', 'json', 'enum'],

        // 初始化
        init() {
            this.loadEntities();
        },

        // 加载实体列表
        async loadEntities() {
            try {
                const response = await fetch('/api/admin/bases/entities');
                const data = await response.json();
                this.entities = data.data;
            } catch (error) {
                console.error('加载实体列表失败:', error);
                alert('加载实体列表失败');
            }
        },

        // 加载字段列表
        async loadFields(entityId) {
            try {
                const response = await fetch(`/api/admin/bases/entities/${entityId}/fields`);
                this.fields = await response.json();
            } catch (error) {
                console.error('加载字段列表失败:', error);
                alert('加载字段列表失败');
            }
        },

        // 加载数据列表
        async loadEntityData(entityId) {
            try {
                const response = await fetch(`/api/admin/bases/entities/${entityId}/data`);
                const data = await response.json();
                this.entityData = data.data;
            } catch (error) {
                console.error('加载数据列表失败:', error);
                alert('加载数据列表失败');
            }
        },

        // 打开实体面板
        openEntityPanel() {
            this.currentEntity = {};
            this.showEntityPanel = true;
        },

        // 关闭实体面板
        closeEntityPanel() {
            this.showEntityPanel = false;
        },

        // 编辑实体
        editEntity(entity) {
            this.currentEntity = { ...entity };
            this.showEntityPanel = true;
        },

        // 保存实体
        async saveEntity() {
            try {
                const url = this.currentEntity.id ? `/api/admin/bases/entities/${this.currentEntity.id}` : '/api/admin/bases/entities';
                const method = this.currentEntity.id ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentEntity),
                });

                if (!response.ok) {
                    throw new Error('保存实体失败');
                }

                await this.loadEntities();
                this.closeEntityPanel();
            } catch (error) {
                console.error('保存实体失败:', error);
                alert('保存实体失败');
            }
        },

        // 删除实体
        async deleteEntity(id) {
            if (!confirm('确定要删除这个实体吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/bases/entities/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) {
                    throw new Error('删除实体失败');
                }

                await this.loadEntities();
            } catch (error) {
                console.error('删除实体失败:', error);
                alert('删除实体失败');
            }
        },

        // 打开字段面板
        async openFieldPanel(entity) {
            this.currentEntity = entity;
            await this.loadFields(entity.id);
            this.showFieldPanel = true;
        },

        // 关闭字段面板
        closeFieldPanel() {
            this.showFieldPanel = false;
        },

        // 打开字段表单
        openFieldForm() {
            this.currentField = {
                entity_id: this.currentEntity.id,
                is_nullable: false,
                is_unique: false,
            };
            this.showFieldForm = true;
        },

        // 关闭字段表单
        closeFieldForm() {
            this.showFieldForm = false;
        },

        // 编辑字段
        editField(field) {
            this.currentField = { ...field };
            this.showFieldForm = true;
        },

        // 保存字段
        async saveField() {
            try {
                const url = this.currentField.id ? `/api/admin/bases/fields/${this.currentField.id}` : '/api/admin/bases/fields';
                const method = this.currentField.id ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.currentField),
                });

                if (!response.ok) {
                    throw new Error('保存字段失败');
                }

                await this.loadFields(this.currentEntity.id);
                this.closeFieldForm();
            } catch (error) {
                console.error('保存字段失败:', error);
                alert('保存字段失败');
            }
        },

        // 删除字段
        async deleteField(id) {
            if (!confirm('确定要删除这个字段吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/bases/fields/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) {
                    throw new Error('删除字段失败');
                }

                await this.loadFields(this.currentEntity.id);
            } catch (error) {
                console.error('删除字段失败:', error);
                alert('删除字段失败');
            }
        },

        // 打开数据面板
        async openDataPanel(entity) {
            this.currentEntity = entity;
            await this.loadFields(entity.id);
            await this.loadEntityData(entity.id);
            this.showDataPanel = true;
        },

        // 关闭数据面板
        closeDataPanel() {
            this.showDataPanel = false;
        },

        // 打开数据表单
        openDataForm() {
            this.currentData = {
                entity_id: this.currentEntity.id,
            };
            // 初始化字段默认值
            this.fields.forEach(field => {
                if (field.default_value) {
                    this.currentData[field.name] = field.default_value;
                }
            });
            this.showDataForm = true;
        },

        // 关闭数据表单
        closeDataForm() {
            this.showDataForm = false;
        },

        // 编辑数据
        editData(data) {
            this.currentData = { ...data };
            this.showDataForm = true;
        },

        // 保存数据
        async saveData() {
            try {
                const url = this.currentData.id ? `/api/admin/bases/data/${this.currentData.id}` : '/api/admin/bases/data';
                const method = this.currentData.id ? 'PUT' : 'POST';
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        entity_id: this.currentEntity.id,
                        data: JSON.stringify(this.currentData),
                    }),
                });

                if (!response.ok) {
                    throw new Error('保存数据失败');
                }

                await this.loadEntityData(this.currentEntity.id);
                this.closeDataForm();
            } catch (error) {
                console.error('保存数据失败:', error);
                alert('保存数据失败');
            }
        },

        // 删除数据
        async deleteData(id) {
            if (!confirm('确定要删除这条数据吗？')) {
                return;
            }

            try {
                const response = await fetch(`/api/admin/bases/data/${id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) {
                    throw new Error('删除数据失败');
                }

                await this.loadEntityData(this.currentEntity.id);
            } catch (error) {
                console.error('删除数据失败:', error);
                alert('删除数据失败');
            }
        },

        // 格式化日期
        formatDate(date) {
            if (!date) return '';
            return new Date(date).toLocaleString();
        },
    };
}