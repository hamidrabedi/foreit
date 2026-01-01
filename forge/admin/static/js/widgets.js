// Forge Admin - Widget Initialization JavaScript

(function () {
    'use strict';

    // Rich Text Editor initialization
    window.initRichTextEditor = function (id) {
        const textarea = document.getElementById(id);
        if (!textarea) return;

        const toolbar = textarea.getAttribute('data-toolbar') || 'full';
        const height = parseInt(textarea.getAttribute('data-height')) || 300;

        // Simple rich text editor implementation
        // In production, you might use TinyMCE, CKEditor, Quill, etc.
        const container = document.createElement('div');
        container.className = 'rich-text-editor-container';
        container.style.height = height + 'px';

        // Create toolbar
        const toolbarDiv = createToolbar(toolbar);
        container.appendChild(toolbarDiv);

        // Create editor div
        const editorDiv = document.createElement('div');
        editorDiv.className = 'rich-text-editor';
        editorDiv.contentEditable = true;
        editorDiv.innerHTML = textarea.value || '';
        editorDiv.style.minHeight = (height - 50) + 'px';
        editorDiv.style.padding = '10px';
        editorDiv.style.border = '1px solid #ccc';
        editorDiv.style.borderTop = 'none';
        editorDiv.style.overflow = 'auto';
        container.appendChild(editorDiv);

        // Hide textarea and insert editor
        textarea.style.display = 'none';
        textarea.parentNode.insertBefore(container, textarea);

        // Sync content back to textarea
        editorDiv.addEventListener('input', function () {
            textarea.value = editorDiv.innerHTML;
        });

        // Handle toolbar buttons
        toolbarDiv.addEventListener('click', function (e) {
            if (e.target.tagName === 'BUTTON') {
                e.preventDefault();
                const command = e.target.getAttribute('data-command');
                if (command) {
                    document.execCommand(command, false, null);
                    editorDiv.focus();
                }
            }
        });
    };

    // Create toolbar for rich text editor
    function createToolbar(type) {
        const toolbar = document.createElement('div');
        toolbar.className = 'rich-text-toolbar';
        toolbar.style.padding = '5px';
        toolbar.style.background = '#f5f5f5';
        toolbar.style.border = '1px solid #ccc';
        toolbar.style.borderBottom = 'none';

        const buttons = getToolbarButtons(type);
        buttons.forEach(function (btn) {
            const button = document.createElement('button');
            button.type = 'button';
            button.className = 'btn btn-sm btn-secondary';
            button.style.marginRight = '5px';
            button.setAttribute('data-command', btn.command);
            button.textContent = btn.label;
            button.title = btn.title;
            toolbar.appendChild(button);
        });

        return toolbar;
    }

    // Get toolbar buttons based on type
    function getToolbarButtons(type) {
        const basicButtons = [
            { command: 'bold', label: 'B', title: 'Bold' },
            { command: 'italic', label: 'I', title: 'Italic' },
            { command: 'underline', label: 'U', title: 'Underline' }
        ];

        const fullButtons = [
            { command: 'bold', label: 'B', title: 'Bold' },
            { command: 'italic', label: 'I', title: 'Italic' },
            { command: 'underline', label: 'U', title: 'Underline' },
            { command: 'strikeThrough', label: 'S', title: 'Strike' },
            { command: 'insertUnorderedList', label: '• List', title: 'Bullet List' },
            { command: 'insertOrderedList', label: '1. List', title: 'Numbered List' },
            { command: 'formatBlock', label: 'H1', title: 'Heading 1' },
            { command: 'formatBlock', label: 'H2', title: 'Heading 2' },
            { command: 'justifyLeft', label: '←', title: 'Align Left' },
            { command: 'justifyCenter', label: '↔', title: 'Align Center' },
            { command: 'justifyRight', label: '→', title: 'Align Right' },
            { command: 'removeFormat', label: '✕', title: 'Clear Formatting' }
        ];

        if (type === 'full') return fullButtons;
        if (type === 'basic') return basicButtons;
        return basicButtons;
    }

    // Select Search (Autocomplete) initialization
    window.initSelectSearch = function (id) {
        const select = document.getElementById(id);
        if (!select) return;

        // Simple autocomplete implementation
        // In production, you might use Select2, Choices.js, etc.
        const placeholder = select.getAttribute('data-placeholder') || 'Search...';
        const allowClear = select.getAttribute('data-allow-clear') === 'true';

        // Create wrapper
        const wrapper = document.createElement('div');
        wrapper.className = 'select-search-wrapper';
        wrapper.style.position = 'relative';

        // Create search input
        const searchInput = document.createElement('input');
        searchInput.type = 'text';
        searchInput.className = 'form-control';
        searchInput.placeholder = placeholder;

        // Set initial value
        const selectedOption = select.options[select.selectedIndex];
        if (selectedOption && selectedOption.value) {
            searchInput.value = selectedOption.text;
        }

        // Create dropdown
        const dropdown = document.createElement('div');
        dropdown.className = 'select-search-dropdown';
        dropdown.style.display = 'none';
        dropdown.style.position = 'absolute';
        dropdown.style.top = '100%';
        dropdown.style.left = '0';
        dropdown.style.right = '0';
        dropdown.style.maxHeight = '200px';
        dropdown.style.overflow = 'auto';
        dropdown.style.border = '1px solid #ccc';
        dropdown.style.background = 'white';
        dropdown.style.zIndex = '1000';

        // Hide original select
        select.style.display = 'none';
        select.parentNode.insertBefore(wrapper, select);
        wrapper.appendChild(searchInput);
        wrapper.appendChild(dropdown);

        // Populate dropdown options
        function updateDropdown(filter) {
            dropdown.innerHTML = '';
            let hasResults = false;

            if (allowClear) {
                const clearOption = createDropdownOption('', '---------', '', filter);
                if (clearOption) {
                    dropdown.appendChild(clearOption);
                    hasResults = true;
                }
            }

            Array.from(select.options).forEach(function (option) {
                if (!option.value && !allowClear) return;
                
                const dropdownOption = createDropdownOption(option.value, option.text, option.text, filter);
                if (dropdownOption) {
                    dropdown.appendChild(dropdownOption);
                    hasResults = true;
                }
            });

            if (!hasResults) {
                const noResults = document.createElement('div');
                noResults.className = 'select-search-option';
                noResults.style.padding = '10px';
                noResults.textContent = 'No results found';
                dropdown.appendChild(noResults);
            }
        }

        // Create dropdown option
        function createDropdownOption(value, text, fullText, filter) {
            if (filter && fullText.toLowerCase().indexOf(filter.toLowerCase()) === -1) {
                return null;
            }

            const option = document.createElement('div');
            option.className = 'select-search-option';
            option.style.padding = '10px';
            option.style.cursor = 'pointer';
            option.setAttribute('data-value', value);
            option.textContent = text;

            option.addEventListener('mouseenter', function () {
                option.style.background = '#f0f0f0';
            });

            option.addEventListener('mouseleave', function () {
                option.style.background = 'white';
            });

            option.addEventListener('click', function () {
                select.value = value;
                searchInput.value = text;
                dropdown.style.display = 'none';
                select.dispatchEvent(new Event('change'));
            });

            return option;
        }

        // Show dropdown on focus
        searchInput.addEventListener('focus', function () {
            updateDropdown('');
            dropdown.style.display = 'block';
        });

        // Filter on input
        searchInput.addEventListener('input', function () {
            updateDropdown(this.value);
            dropdown.style.display = 'block';
        });

        // Hide dropdown on blur (with delay for click handling)
        searchInput.addEventListener('blur', function () {
            setTimeout(function () {
                dropdown.style.display = 'none';
            }, 200);
        });

        // Clear button
        if (allowClear && searchInput.value) {
            addClearButton(wrapper, searchInput, select);
        }

        searchInput.addEventListener('input', function () {
            if (allowClear) {
                const clearBtn = wrapper.querySelector('.select-search-clear');
                if (this.value && !clearBtn) {
                    addClearButton(wrapper, searchInput, select);
                } else if (!this.value && clearBtn) {
                    clearBtn.remove();
                }
            }
        });
    };

    // Add clear button to select search
    function addClearButton(wrapper, searchInput, select) {
        const clearBtn = document.createElement('button');
        clearBtn.type = 'button';
        clearBtn.className = 'select-search-clear';
        clearBtn.innerHTML = '×';
        clearBtn.style.position = 'absolute';
        clearBtn.style.right = '10px';
        clearBtn.style.top = '50%';
        clearBtn.style.transform = 'translateY(-50%)';
        clearBtn.style.border = 'none';
        clearBtn.style.background = 'none';
        clearBtn.style.fontSize = '20px';
        clearBtn.style.cursor = 'pointer';
        clearBtn.style.color = '#999';

        clearBtn.addEventListener('click', function (e) {
            e.preventDefault();
            searchInput.value = '';
            select.value = '';
            select.dispatchEvent(new Event('change'));
            clearBtn.remove();
        });

        wrapper.style.position = 'relative';
        wrapper.appendChild(clearBtn);
    }

    // File upload preview
    window.initFileUploadPreview = function (id) {
        const input = document.getElementById(id);
        if (!input) return;

        input.addEventListener('change', function () {
            const files = this.files;
            if (files.length === 0) return;

            const maxSize = parseInt(this.getAttribute('data-max-size')) || (10 * 1024 * 1024);
            const maxFiles = parseInt(this.getAttribute('data-max-files')) || 1;

            if (files.length > maxFiles) {
                alert('Maximum ' + maxFiles + ' file(s) allowed');
                this.value = '';
                return;
            }

            for (let i = 0; i < files.length; i++) {
                if (files[i].size > maxSize) {
                    alert('File "' + files[i].name + '" is too large. Maximum size: ' + (maxSize / 1024 / 1024) + 'MB');
                    this.value = '';
                    return;
                }
            }

            // Show preview for images
            const accept = this.getAttribute('accept') || '';
            if (accept.includes('image')) {
                showImagePreview(this, files);
            } else {
                showFileInfo(this, files);
            }
        });
    };

    // Show image preview
    function showImagePreview(input, files) {
        const preview = document.getElementById(input.id + '_preview') ||
            createPreviewContainer(input);

        preview.innerHTML = '';

        Array.from(files).forEach(function (file) {
            const reader = new FileReader();
            reader.onload = function (e) {
                const img = document.createElement('img');
                img.src = e.target.result;
                img.style.maxWidth = '200px';
                img.style.maxHeight = '200px';
                img.style.margin = '5px';
                preview.appendChild(img);
            };
            reader.readAsDataURL(file);
        });
    }

    // Show file info
    function showFileInfo(input, files) {
        const preview = document.getElementById(input.id + '_preview') ||
            createPreviewContainer(input);

        preview.innerHTML = '';

        Array.from(files).forEach(function (file) {
            const info = document.createElement('div');
            info.textContent = file.name + ' (' + formatFileSize(file.size) + ')';
            preview.appendChild(info);
        });
    }

    // Create preview container
    function createPreviewContainer(input) {
        const preview = document.createElement('div');
        preview.id = input.id + '_preview';
        preview.className = 'file-upload-preview';
        preview.style.marginTop = '10px';
        input.parentNode.insertBefore(preview, input.nextSibling);
        return preview;
    }

    // Format file size
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    }

})();
