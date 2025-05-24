function imagemanagerManagement() {
    return {
        images: [],
        showUploadModal: false,
        showPreviewModal: false,
        selectedImages: [],
        currentImage: null,

        init() {
            this.fetchImages();
        },

        async fetchImages() {
            try {
                const response = await fetch('/api/admin/imagemanager/list');
                if (!response.ok) throw new Error('获取图片列表失败');
                this.images = await response.json();
            } catch (error) {
                ShowError(error.message);
            }
        },

        uploadImage() {
            this.selectedImages = [];
            this.showUploadModal = true;
        },

        handleImageSelect(event) {
            const files = Array.from(event.target.files).filter(file => file.type.startsWith('image/'));
            
            files.forEach(file => {
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.selectedImages.push({
                        file: file,
                        preview: e.target.result
                    });
                };
                reader.readAsDataURL(file);
            });
        },

        removeImage(index) {
            this.selectedImages.splice(index, 1);
        },

        async submitUpload() {
            try {
                const formData = new FormData();
                this.selectedImages.forEach(image => {
                    formData.append('images', image.file);
                });

                const response = await fetch('/api/admin/imagemanager/upload', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '上传失败');
                }

                await this.fetchImages();
                this.showUploadModal = false;
                this.selectedImages = [];
                ShowMessage('图片上传成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        async deleteImage(image) {
            if (!confirm(`确定要删除图片 "${image.name}" 吗？`)) return;

            try {
                const response = await fetch(`/api/admin/imagemanager/delete/${image.id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除失败');
                }

                await this.fetchImages();
                ShowMessage('图片删除成功');
            } catch (error) {
                ShowError(error.message);
            }
        },

        viewImage(image) {
            this.currentImage = image;
            this.showPreviewModal = true;
        },

        async copyUrl(image) {
            try {
                await navigator.clipboard.writeText(image.url);
                ShowMessage('图片链接已复制到剪贴板');
            } catch (error) {
                ShowError('复制失败，请手动复制');
            }
        },

        formatSize(size) {
            if (size === 0) return '0 B';
            const units = ['B', 'KB', 'MB', 'GB', 'TB'];
            const k = 1024;
            const i = Math.floor(Math.log(size) / Math.log(k));
            return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + units[i];
        }
    };
} 