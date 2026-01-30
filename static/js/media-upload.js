document.addEventListener('DOMContentLoaded', function () {
    var zone = document.getElementById('media-upload-zone');
    var fileInput = document.getElementById('media-file-input');
    var uploadBtn = document.getElementById('media-upload-btn');
    var progressWrap = document.getElementById('media-upload-progress');
    var progressFill = document.getElementById('media-progress-fill');
    var progressText = document.getElementById('media-progress-text');
    var uploadList = document.getElementById('media-upload-list');
    var contentTextarea = document.getElementById('content');

    if (!zone || !fileInput) return;

    var csrfInput = document.querySelector('input[name="csrf_token"]');
    var csrfToken = csrfInput ? csrfInput.value : '';

    var allowedTypes = [
        'image/jpeg', 'image/png', 'image/gif', 'image/webp',
        'video/mp4', 'video/webm'
    ];
    var maxSize = 10 * 1024 * 1024;

    // Click to browse
    uploadBtn.addEventListener('click', function (e) {
        e.preventDefault();
        fileInput.click();
    });

    zone.addEventListener('click', function (e) {
        if (e.target === uploadBtn) return;
        fileInput.click();
    });

    fileInput.addEventListener('change', function () {
        if (this.files.length > 0) {
            uploadFiles(this.files);
            this.value = '';
        }
    });

    // Drag and drop
    zone.addEventListener('dragover', function (e) {
        e.preventDefault();
        zone.classList.add('drag-over');
    });

    zone.addEventListener('dragleave', function (e) {
        e.preventDefault();
        zone.classList.remove('drag-over');
    });

    zone.addEventListener('drop', function (e) {
        e.preventDefault();
        zone.classList.remove('drag-over');
        if (e.dataTransfer.files.length > 0) {
            uploadFiles(e.dataTransfer.files);
        }
    });

    function uploadFiles(files) {
        for (var i = 0; i < files.length; i++) {
            uploadSingleFile(files[i]);
        }
    }

    function uploadSingleFile(file) {
        if (allowedTypes.indexOf(file.type) === -1) {
            alert('File type not allowed: ' + file.name + '\nAllowed: JPEG, PNG, GIF, WebP, MP4, WebM');
            return;
        }
        if (file.size > maxSize) {
            alert('File too large: ' + file.name + '\nMaximum size is 10MB.');
            return;
        }

        var formData = new FormData();
        formData.append('file', file);

        progressWrap.style.display = 'block';
        progressFill.style.width = '0%';
        progressText.textContent = 'Uploading ' + file.name + '...';

        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/media/upload', true);
        xhr.setRequestHeader('X-CSRF-Token', csrfToken);

        xhr.upload.addEventListener('progress', function (e) {
            if (e.lengthComputable) {
                var pct = Math.round((e.loaded / e.total) * 100);
                progressFill.style.width = pct + '%';
                progressText.textContent = 'Uploading ' + file.name + '... ' + pct + '%';
            }
        });

        xhr.addEventListener('load', function () {
            progressWrap.style.display = 'none';

            if (xhr.status === 200) {
                var data = JSON.parse(xhr.responseText);
                addMediaItem(data);
            } else {
                var errData;
                try {
                    errData = JSON.parse(xhr.responseText);
                } catch (e) {
                    errData = { error: 'Upload failed.' };
                }
                alert('Upload error: ' + errData.error);
            }
        });

        xhr.addEventListener('error', function () {
            progressWrap.style.display = 'none';
            alert('Upload failed. Please try again.');
        });

        xhr.send(formData);
    }

    function addMediaItem(data) {
        var item = document.createElement('div');
        item.className = 'media-upload-item';

        var isVideo = data.mime_type.indexOf('video') === 0;
        var previewHTML;
        if (isVideo) {
            previewHTML = '<div class="media-upload-item-preview" style="display:flex;align-items:center;justify-content:center;font-size:1.5rem;color:var(--text-muted);">&#9654;</div>';
        } else {
            previewHTML = '<img class="media-upload-item-preview" src="' + escapeAttr(data.preview_url) + '" alt="">';
        }

        item.innerHTML = previewHTML +
            '<div class="media-upload-item-info">' +
                '<div class="media-upload-item-name">' + escapeHtml(data.original_name) + '</div>' +
                '<div class="media-upload-item-size">' + formatFileSize(data.file_size) + '</div>' +
            '</div>' +
            '<span class="media-upload-item-tag" title="Click to copy embed tag">' + escapeHtml(data.embed_tag) + '</span>' +
            '<button type="button" class="media-insert-btn">Insert</button>' +
            '<span class="media-upload-item-copied">Copied!</span>';

        uploadList.appendChild(item);

        var tagEl = item.querySelector('.media-upload-item-tag');
        var copiedEl = item.querySelector('.media-upload-item-copied');
        tagEl.addEventListener('click', function () {
            navigator.clipboard.writeText(data.embed_tag).then(function () {
                copiedEl.style.display = 'inline';
                setTimeout(function () { copiedEl.style.display = 'none'; }, 2000);
            });
        });

        var insertBtn = item.querySelector('.media-insert-btn');
        insertBtn.addEventListener('click', function () {
            insertAtCursor(contentTextarea, data.embed_tag);
        });
    }

    function insertAtCursor(textarea, text) {
        if (!textarea) return;
        var start = textarea.selectionStart;
        var end = textarea.selectionEnd;
        var before = textarea.value.substring(0, start);
        var after = textarea.value.substring(end);
        var insert = '\n' + text + '\n';
        textarea.value = before + insert + after;
        textarea.selectionStart = textarea.selectionEnd = start + insert.length;
        textarea.focus();
    }

    function formatFileSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / 1048576).toFixed(1) + ' MB';
    }

    function escapeHtml(str) {
        var div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    function escapeAttr(str) {
        return str.replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/'/g, '&#39;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    }
});
