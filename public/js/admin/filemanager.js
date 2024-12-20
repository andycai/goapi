$(document).ready(function() {
    let currentPath = './';

    // Load initial file list
    loadFileList(currentPath);

    // Handle breadcrumb navigation
    $('#pathBreadcrumb').on('click', 'a', function(e) {
        e.preventDefault();
        currentPath = $(this).data('path');
        loadFileList(currentPath);
    });

    // Upload file button
    $('#btnUpload').click(function() {
        $('#uploadModal').modal('show');
    });

    // Create directory button
    $('#btnCreateDir').click(function() {
        $('#createDirModal').modal('show');
    });

    // Create file button
    $('#btnCreateFile').click(function() {
        $('#createFileModal').modal('show');
    });

    // Confirm upload
    $('#btnConfirmUpload').click(function() {
        const formData = new FormData();
        const fileInput = $('#fileInput')[0];
        if (fileInput.files.length === 0) {
            toastr.error('请选择文件');
            return;
        }
        formData.append('file', fileInput.files[0]);
        formData.append('path', currentPath);

        $.ajax({
            url: '/api/filemanager/upload',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function(response) {
                $('#uploadModal').modal('hide');
                $('#uploadForm')[0].reset();
                toastr.success('文件上传成功');
                loadFileList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '上传失败');
            }
        });
    });

    // Confirm create directory
    $('#btnConfirmCreateDir').click(function() {
        const dirName = $('#dirName').val();
        if (!dirName) {
            toastr.error('请输入文件夹名称');
            return;
        }

        $.ajax({
            url: '/api/filemanager/create',
            type: 'POST',
            data: {
                path: currentPath + dirName,
                is_dir: 'true'
            },
            success: function(response) {
                $('#createDirModal').modal('hide');
                $('#createDirForm')[0].reset();
                toastr.success('文件夹创建成功');
                loadFileList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '创建失败');
            }
        });
    });

    // Confirm create file
    $('#btnConfirmCreateFile').click(function() {
        const fileName = $('#fileName').val();
        if (!fileName) {
            toastr.error('请输入文件名称');
            return;
        }

        $.ajax({
            url: '/api/filemanager/create',
            type: 'POST',
            data: {
                path: currentPath + fileName,
                is_dir: 'false'
            },
            success: function(response) {
                $('#createFileModal').modal('hide');
                $('#createFileForm')[0].reset();
                toastr.success('文件创建成功');
                loadFileList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '创建失败');
            }
        });
    });

    // File list event delegation
    $('#fileList').on('click', '.btn-rename', function() {
        const path = $(this).data('path');
        const name = $(this).data('name');
        $('#renameModal').data('path', path).modal('show');
        $('#newName').val(name);
    });

    $('#fileList').on('click', '.btn-move', function() {
        const path = $(this).data('path');
        $('#moveModal').data('path', path).modal('show');
    });

    $('#fileList').on('click', '.btn-delete', function() {
        const path = $(this).data('path');
        if (confirm('确定要删除吗？')) {
            deleteFile(path);
        }
    });

    $('#fileList').on('click', '.btn-download', function() {
        const path = $(this).data('path');
        window.location.href = '/api/filemanager/download?path=' + encodeURIComponent(path);
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
            url: '/api/filemanager/rename',
            type: 'POST',
            data: {
                old_path: oldPath,
                new_path: newPath
            },
            success: function(response) {
                $('#renameModal').modal('hide');
                toastr.success('重命名成功');
                loadFileList(currentPath);
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
            url: '/api/filemanager/move',
            type: 'POST',
            data: {
                source_path: sourcePath,
                dest_path: destPath
            },
            success: function(response) {
                $('#moveModal').modal('hide');
                $('#moveForm')[0].reset();
                toastr.success('移动成功');
                loadFileList(currentPath);
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
            url: '/api/filemanager/copy',
            type: 'POST',
            data: {
                source_path: sourcePath,
                dest_path: destPath
            },
            success: function(response) {
                $('#moveModal').modal('hide');
                $('#moveForm')[0].reset();
                toastr.success('复制成功');
                loadFileList(currentPath);
            },
            error: function(xhr) {
                toastr.error(xhr.responseJSON?.error || '复制失败');
            }
        });
    });

    // Load file list function
    function loadFileList(path) {
        $.ajax({
            url: '/api/filemanager/list',
            type: 'GET',
            data: { path: path },
            success: function(response) {
                updateBreadcrumb(path);
                updateFileList(response.files);
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

    // Update file list function
    function updateFileList(files) {
        let html = '';
        files.forEach(file => {
            const icon = file.is_dir ? 'fa-folder' : 'fa-file';
            const size = file.is_dir ? '-' : formatFileSize(file.size);
            const downloadBtn = file.is_dir ? '' : `<button class="btn btn-success btn-sm btn-download" data-path="${file.path}"><i class="fas fa-download"></i></button>`;

            html += `
                <tr>
                    <td><i class="fas ${icon} mr-2"></i>${file.name}</td>
                    <td>${size}</td>
                    <td>${new Date(file.mod_time).toLocaleString()}</td>
                    <td>${file.mode}</td>
                    <td>
                        ${downloadBtn}
                        <button class="btn btn-primary btn-sm btn-rename" data-path="${file.path}" data-name="${file.name}">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="btn btn-info btn-sm btn-move" data-path="${file.path}">
                            <i class="fas fa-arrows-alt"></i>
                        </button>
                        <button class="btn btn-danger btn-sm btn-delete" data-path="${file.path}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
        });

        $('#fileList').html(html);
    }

    // Delete file function
    function deleteFile(path) {
        $.ajax({
            url: '/api/filemanager/delete',
            type: 'POST',
            data: { path: path },
            success: function(response) {
                toastr.success('删除成功');
                loadFileList(currentPath);
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