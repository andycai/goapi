$(document).ready(function() {
    let currentPath = './';

    // Load initial image list
    loadImageList(currentPath);

    // Handle breadcrumb navigation
    $('#pathBreadcrumb').on('click', 'a', function(e) {
        e.preventDefault();
        currentPath = $(this).data('path');
        loadImageList(currentPath);
    });

    // Upload button
    $('#btnUpload').click(function() {
        $('#uploadModal').modal('show');
    });

    // Confirm upload
    $('#btnConfirmUpload').click(function() {
        const formData = new FormData();
        const fileInput = $('#imageInput')[0];
        if (fileInput.files.length === 0) {
            toastr.error('请选择图片');
            return;
        }
        formData.append('file', fileInput.files[0]);
        formData.append('path', currentPath);

        $.ajax({
            url: '/api/imagemanager/upload',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function(response) {
                $('#uploadModal').modal('hide');
                $('#uploadForm')[0].reset();
                toastr.success('图片上传成功');
                loadImageList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '上传失败');
            }
        });
    });

    // Image grid event delegation
    $('#imageGrid').on('click', '.btn-preview', function() {
        const path = $(this).data('path');
        const name = $(this).data('name');
        showImagePreview(path, name);
    });

    $('#imageGrid').on('click', '.btn-rename', function() {
        const path = $(this).data('path');
        const name = $(this).data('name');
        $('#renameModal').data('path', path).modal('show');
        $('#newName').val(name);
    });

    $('#imageGrid').on('click', '.btn-move', function() {
        const path = $(this).data('path');
        $('#moveModal').data('path', path).modal('show');
    });

    $('#imageGrid').on('click', '.btn-delete', function() {
        const path = $(this).data('path');
        if (confirm('确定要删除这张图片吗？')) {
            deleteImage(path);
        }
    });

    // Confirm rename
    $('#btnConfirmRename').click(function() {
        const oldPath = $('#renameModal').data('path');
        const newName = $('#newName').val();
        if (!newName) {
            toastr.error('请输入新名称');
            return;
        }

        const pathParts = oldPath.split('/');
        pathParts.pop();
        const newPath = pathParts.join('/') + '/' + newName;

        $.ajax({
            url: '/api/imagemanager/rename',
            type: 'POST',
            data: {
                old_path: oldPath,
                new_path: newPath
            },
            success: function(response) {
                $('#renameModal').modal('hide');
                toastr.success('重命名成功');
                loadImageList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '重命名失败');
            }
        });
    });

    // Confirm move
    $('#btnConfirmMove').click(function() {
        const sourcePath = $('#moveModal').data('path');
        const destPath = $('#destPath').val();
        if (!destPath) {
            toastr.error('请输入目标路径');
            return;
        }

        $.ajax({
            url: '/api/imagemanager/move',
            type: 'POST',
            data: {
                source_path: sourcePath,
                dest_path: destPath
            },
            success: function(response) {
                $('#moveModal').modal('hide');
                $('#moveForm')[0].reset();
                toastr.success('移动成功');
                loadImageList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '移动失败');
            }
        });
    });

    // Confirm copy
    $('#btnConfirmCopy').click(function() {
        const sourcePath = $('#moveModal').data('path');
        const destPath = $('#destPath').val();
        if (!destPath) {
            toastr.error('请输入目标路径');
            return;
        }

        $.ajax({
            url: '/api/imagemanager/copy',
            type: 'POST',
            data: {
                source_path: sourcePath,
                dest_path: destPath
            },
            success: function(response) {
                $('#moveModal').modal('hide');
                $('#moveForm')[0].reset();
                toastr.success('复制成功');
                loadImageList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '复制失败');
            }
        });
    });

    // Load image list function
    function loadImageList(path) {
        $.ajax({
            url: '/api/imagemanager/list',
            type: 'GET',
            data: { path: path },
            success: function(response) {
                updateBreadcrumb(path);
                updateImageGrid(response.images);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '加载失败');
            }
        });
    }

    // Update breadcrumb function
    function updateBreadcrumb(path) {
        const parts = path.split('/').filter(p => p !== '');
        let html = '<li class="breadcrumb-item"><a href="#" data-path="./">根目录</a></li>';
        let currentPath = './';

        parts.forEach(part => {
            currentPath += part + '/';
            html += `<li class="breadcrumb-item"><a href="#" data-path="${currentPath}">${part}</a></li>`;
        });

        $('#pathBreadcrumb').html(html);
    }

    // Update image grid function
    function updateImageGrid(images) {
        let html = '';
        images.forEach(image => {
            html += `
                <div class="col-md-3">
                    <div class="card image-card">
                        <img src="/api/imagemanager/thumbnail?path=${image.path}" class="card-img-top" alt="${image.name}">
                        <div class="card-body">
                            <p class="card-text text-truncate">${image.name}</p>
                            <small class="text-muted">${formatFileSize(image.size)}</small>
                        </div>
                        <div class="card-actions">
                            <button class="btn btn-primary btn-sm btn-preview" data-path="${image.path}" data-name="${image.name}">
                                <i class="fas fa-eye"></i>
                            </button>
                            <button class="btn btn-info btn-sm btn-rename" data-path="${image.path}" data-name="${image.name}">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-warning btn-sm btn-move" data-path="${image.path}">
                                <i class="fas fa-arrows-alt"></i>
                            </button>
                            <button class="btn btn-danger btn-sm btn-delete" data-path="${image.path}">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            `;
        });

        $('#imageGrid').html(html || '<div class="col-12"><p class="text-center">暂无图片</p></div>');
    }

    // Show image preview function
    function showImagePreview(path, name) {
        $.ajax({
            url: '/api/imagemanager/info',
            type: 'GET',
            data: { path: path },
            success: function(response) {
                const info = response.info;
                $('#previewImage').attr('src', '/api/imagemanager/view?path=' + path);
                $('#imageInfo').html(`
                    <strong>文件名：</strong>${info.name}<br>
                    <strong>尺寸：</strong>${info.width} x ${info.height}<br>
                    <strong>大小：</strong>${formatFileSize(info.size)}<br>
                    <strong>格式：</strong>${info.format.toUpperCase()}<br>
                    <strong>修改时间：</strong>${new Date(info.mod_time).toLocaleString()}
                `);
                $('#btnDownload').attr('href', '/api/imagemanager/view?path=' + path);
                $('#previewModal').modal('show');
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '获取图片信息失败');
            }
        });
    }

    // Delete image function
    function deleteImage(path) {
        $.ajax({
            url: '/api/imagemanager/delete',
            type: 'POST',
            data: { path: path },
            success: function(response) {
                toastr.success('删除成功');
                loadImageList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '删除失败');
            }
        });
    }

    // Format file size function
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
}); 