/**
 * 字典管理
 */
function dictManagement() {
    return {
        // 字典类型数据
        dictTypes: [],
        // 字典数据
        dictData: [],
        // 分页信息
        currentTypePage: 1,
        typesPageSize: 10,
        totalTypeRecords: 0,
        typePages: 1,
        // 字典数据分页
        currentDataPage: 1,
        dataPageSize: 10,
        totalDataRecords: 0,
        dataPages: 1,
        // 当前选中的字典类型
        currentDictType: '',
        currentDictTypeName: '',
        // 表单数据
        typeForm: {
            id: 0,
            name: '',
            type: '',
            remark: ''
        },
        dataForm: {
            id: 0,
            type: '',
            label: '',
            value: '',
            sort: 0,
            remark: ''
        },
        // 模态框状态
        typeModalAction: 'add',
        dataModalAction: 'add',
        // 显示字典数据列表
        showDictData: false,

        init() {
            this.loadDictTypes();
        },

        // 加载字典类型列表
        async loadDictTypes() {
            try {
                const response = await fetch(`/api/admin/dict/type/list?page=${this.currentTypePage}&limit=${this.typesPageSize}`);
                if (!response.ok) throw new Error('加载字典类型失败');
                const data = await response.json();
                this.dictTypes = data.dictTypes;
                this.totalTypeRecords = data.total;
                this.typePages = Math.ceil(this.totalTypeRecords / this.typesPageSize);
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        // 加载字典数据列表
        async loadDictData() {
            if (!this.currentDictType) return;
            
            try {
                const response = await fetch(`/api/admin/dict/data/list?type=${this.currentDictType}&page=${this.currentDataPage}&limit=${this.dataPageSize}`);
                if (!response.ok) throw new Error('加载字典数据失败');
                const data = await response.json();
                this.dictData = data.dictData;
                this.totalDataRecords = data.total;
                this.dataPages = Math.ceil(this.totalDataRecords / this.dataPageSize);
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        // 切换字典类型页码
        async changeTypePage(page) {
            if (page < 1 || page > this.typePages) return;
            this.currentTypePage = page;
            await this.loadDictTypes();
        },

        // 切换字典数据页码
        async changeDataPage(page) {
            if (page < 1 || page > this.dataPages) return;
            this.currentDataPage = page;
            await this.loadDictData();
        },

        // 查看字典数据
        async viewDictData(type, name) {
            this.currentDictType = type;
            this.currentDictTypeName = name;
            this.showDictData = true;
            this.currentDataPage = 1;
            await this.loadDictData();
        },

        // 打开添加字典类型模态框
        openAddTypeModal() {
            this.typeForm = {
                id: 0,
                name: '',
                type: '',
                remark: ''
            };
            this.typeModalAction = 'add';
        },

        // 打开编辑字典类型模态框
        openEditTypeModal(dictType) {
            this.typeForm = {
                id: dictType.id,
                name: dictType.name,
                type: dictType.type,
                remark: dictType.remark || ''
            };
            this.typeModalAction = 'edit';
        },

        // 保存字典类型
        async saveType() {
            try {
                const url = this.typeModalAction === 'add' ? '/api/admin/dict/type/add' : '/api/admin/dict/type/edit';
                
                // 确保id是整数
                const formData = {
                    ...this.typeForm,
                    id: parseInt(this.typeForm.id, 10) || 0
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
                    throw new Error(error.error || '保存字典类型失败');
                }

                await this.loadDictTypes();
                Alpine.store('notification').show(this.typeModalAction === 'add' ? '添加字典类型成功' : '更新字典类型成功', 'success');
                
                // 关闭模态框
                document.querySelector('#typeModal').querySelector('[x-ref="close"]').click();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        // 删除字典类型
        async deleteType(id) {
            if (!confirm('确定要删除这个字典类型吗？这将同时删除所有关联的字典数据！')) {
                return;
            }

            try {
                const response = await fetch('/api/admin/dict/type/delete', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ id: parseInt(id, 10) })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除字典类型失败');
                }

                await this.loadDictTypes();
                
                // 如果当前显示的是被删除的字典数据，隐藏
                if (this.showDictData) {
                    const dictType = this.dictTypes.find(t => t.id === id);
                    if (dictType && dictType.type === this.currentDictType) {
                        this.showDictData = false;
                    }
                }
                
                Alpine.store('notification').show('删除字典类型成功', 'success');
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        // 打开添加字典数据模态框
        openAddDataModal() {
            if (!this.currentDictType) {
                Alpine.store('notification').show('请先选择一个字典类型', 'warning');
                return;
            }
            
            this.dataForm = {
                id: 0,
                type: this.currentDictType,
                label: '',
                value: '',
                sort: 0,
                remark: ''
            };
            this.dataModalAction = 'add';
        },

        // 打开编辑字典数据模态框
        openEditDataModal(dictData) {
            this.dataForm = {
                id: dictData.id,
                type: dictData.type,
                label: dictData.label,
                value: dictData.value,
                sort: dictData.sort || 0,
                remark: dictData.remark || ''
            };
            this.dataModalAction = 'edit';
        },

        // 保存字典数据
        async saveData() {
            try {
                const url = this.dataModalAction === 'add' ? '/api/admin/dict/data/add' : '/api/admin/dict/data/edit';
                
                // 确保id和sort是整数
                const formData = {
                    ...this.dataForm,
                    id: parseInt(this.dataForm.id, 10) || 0,
                    sort: parseInt(this.dataForm.sort, 10) || 0
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
                    throw new Error(error.error || '保存字典数据失败');
                }

                await this.loadDictData();
                Alpine.store('notification').show(this.dataModalAction === 'add' ? '添加字典数据成功' : '更新字典数据成功', 'success');
                
                // 关闭模态框
                document.querySelector('#dataModal').querySelector('[x-ref="close"]').click();
            } catch (error) {
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        // 删除字典数据
        async deleteData(id) {
            if (!confirm('确定要删除这个字典数据吗？')) {
                return;
            }

            try {
                const response = await fetch('/api/admin/dict/data/delete', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ id: parseInt(id, 10) })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除字典数据失败');
                }

                await this.loadDictData();
                Alpine.store('notification').show('删除字典数据成功', 'success');
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