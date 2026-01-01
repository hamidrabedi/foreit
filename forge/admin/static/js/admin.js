// Forge Admin - Modern Admin Interface JavaScript

(function () {
    'use strict';

    // Initialize admin interface
    document.addEventListener('DOMContentLoaded', function () {
        initSelectAll();
        initBulkActions();
        initFieldsetToggles();
        initFormValidation();
        initSearch();
        initPrepopulatedFields();
        initColumnSorting();
        initMobileMenu();
    });

    // Select all checkbox functionality
    function initSelectAll() {
        const selectAllCheckboxes = document.querySelectorAll('.select-all');
        selectAllCheckboxes.forEach(function (checkbox) {
            checkbox.addEventListener('change', function () {
                const checked = this.checked;
                const rowCheckboxes = document.querySelectorAll('input[name="selected"]');
                rowCheckboxes.forEach(function (cb) {
                    cb.checked = checked;
                });
            });
        });
    }

    // Bulk actions functionality
    function initBulkActions() {
        const bulkActionForm = document.querySelector('.admin-bulk-actions');
        if (!bulkActionForm) return;

        const actionSelect = bulkActionForm.querySelector('.action-select');
        const goButton = bulkActionForm.querySelector('button');

        if (goButton) {
            goButton.addEventListener('click', function () {
                const action = actionSelect.value;
                if (!action) {
                    alert('Please select an action.');
                    return;
                }

                const selected = getSelectedItems();
                if (selected.length === 0) {
                    alert('Please select at least one item.');
                    return;
                }

                if (confirm('Are you sure you want to perform this action on ' + selected.length + ' item(s)?')) {
                    executeBulkAction(action, selected);
                }
            });
        }
    }

    // Initialize column sorting
    function initColumnSorting() {
        const sortableHeaders = document.querySelectorAll('.admin-table th.sortable');
        sortableHeaders.forEach(function (header) {
            const sortLink = header.querySelector('.sort-link');
            if (sortLink) {
                // Add click handler to the entire header for better UX
                header.addEventListener('click', function (e) {
                    // Don't navigate if clicking directly on the link (it will handle it)
                    if (e.target === sortLink || sortLink.contains(e.target)) {
                        return;
                    }
                    // Otherwise, navigate to sort URL
                    window.location.href = sortLink.href;
                });
            }
        });
    }

    // Initialize mobile menu
    function initMobileMenu() {
        const menuToggle = document.querySelector('.admin-menu-toggle');
        const sidebar = document.querySelector('.admin-sidebar');
        
        if (menuToggle && sidebar) {
            // Create overlay if it doesn't exist
            let overlay = document.querySelector('.admin-sidebar-overlay');
            if (!overlay) {
                overlay = document.createElement('div');
                overlay.className = 'admin-sidebar-overlay';
                const container = document.querySelector('.admin-container');
                if (container) {
                    container.insertBefore(overlay, sidebar);
                }
            }

            menuToggle.addEventListener('click', function () {
                sidebar.classList.toggle('open');
                overlay.classList.toggle('active');
            });

            overlay.addEventListener('click', function () {
                sidebar.classList.remove('open');
                overlay.classList.remove('active');
            });
        }
    }

    // Get selected items
    function getSelectedItems() {
        const selected = [];
        const checkboxes = document.querySelectorAll('input[name="selected"]:checked');
        checkboxes.forEach(function (cb) {
            selected.push(cb.value);
        });
        return selected;
    }

    // Execute bulk action
    function executeBulkAction(action, selected) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = window.location.pathname + 'bulk-action/';

        const actionInput = document.createElement('input');
        actionInput.type = 'hidden';
        actionInput.name = 'action';
        actionInput.value = action;
        form.appendChild(actionInput);

        selected.forEach(function (id) {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = 'selected';
            input.value = id;
            form.appendChild(input);
        });

        // Add CSRF token if available
        const csrfToken = document.querySelector('input[name="csrf_token"]');
        if (csrfToken) {
            const csrfInput = document.createElement('input');
            csrfInput.type = 'hidden';
            csrfInput.name = 'csrf_token';
            csrfInput.value = csrfToken.value;
            form.appendChild(csrfInput);
        }

        document.body.appendChild(form);
        form.submit();
    }

    // Fieldset toggle functionality
    function initFieldsetToggles() {
        const toggleButtons = document.querySelectorAll('.toggle-fieldset');
        toggleButtons.forEach(function (button) {
            button.addEventListener('click', function () {
                const fieldset = this.closest('.form-fieldset');
                if (fieldset) {
                    fieldset.classList.toggle('collapsed');
                    this.textContent = fieldset.classList.contains('collapsed') ? 'Expand' : 'Collapse';
                }
            });
        });
    }

    // Form validation
    function initFormValidation() {
        const forms = document.querySelectorAll('.admin-form');
        forms.forEach(function (form) {
            form.addEventListener('submit', function (e) {
                if (!form.checkValidity()) {
                    e.preventDefault();
                    e.stopPropagation();
                }
                form.classList.add('was-validated');
            });
        });
    }

    // Search functionality
    function initSearch() {
        const searchForm = document.querySelector('.admin-search-form');
        if (!searchForm) return;

        const searchInput = searchForm.querySelector('.search-input');
        if (searchInput) {
            // Debounce search
            let timeout;
            searchInput.addEventListener('input', function () {
                clearTimeout(timeout);
                timeout = setTimeout(function () {
                    // Auto-submit on enter or after delay
                }, 500);
            });
        }
    }

    // Utility: Show message
    window.adminShowMessage = function (message, type) {
        type = type || 'info';
        const alertDiv = document.createElement('div');
        alertDiv.className = 'alert alert-' + type;
        alertDiv.textContent = message;

        const messagesContainer = document.querySelector('.admin-messages');
        if (messagesContainer) {
            messagesContainer.appendChild(alertDiv);
            setTimeout(function () {
                alertDiv.remove();
            }, 5000);
        }
    };

    // Utility: Confirm delete
    window.adminConfirmDelete = function (url) {
        if (confirm('Are you sure you want to delete this item? This action cannot be undone.')) {
            window.location.href = url;
        }
    };

    // Prepopulated fields
    function initPrepopulatedFields() {
        const forms = document.querySelectorAll('form[data-prepopulated-fields]');
        forms.forEach(function (form) {
            const attr = form.getAttribute('data-prepopulated-fields');
            if (!attr || attr === 'null') return;

            let prepopulated;
            try {
                prepopulated = JSON.parse(attr);
            } catch (e) {
                console.error('Failed to parse prepopulated fields:', e);
                return;
            }

            for (const targetField in prepopulated) {
                const sourceFields = prepopulated[targetField];
                const targetInput = form.querySelector('[name="' + targetField + '"]');
                if (!targetInput) continue;

                // Track if user has manually edited the target field
                let isManuallyEdited = targetInput.value !== '';

                targetInput.addEventListener('input', function () {
                    isManuallyEdited = true;
                });

                const sourceInputs = [];
                sourceFields.forEach(function (f) {
                    const i = form.querySelector('[name="' + f + '"]');
                    if (i) sourceInputs.push(i);
                });

                sourceInputs.forEach(function (input) {
                    input.addEventListener('input', function () {
                        if (isManuallyEdited) return;

                        const values = sourceInputs.map(function (i) { return i.value; });
                        const slug = values.join(' ').toLowerCase()
                            .replace(/[^\w\s-]/g, '')
                            .replace(/[\s_-]+/g, '-')
                            .replace(/^-+|-+$/g, '');
                        targetInput.value = slug;
                    });
                });
            }
        });
    }

})();
