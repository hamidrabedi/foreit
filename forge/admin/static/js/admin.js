// Forge Admin - Modern Admin Interface JavaScript

(function() {
    'use strict';

    // Initialize admin interface
    document.addEventListener('DOMContentLoaded', function() {
        initSelectAll();
        initBulkActions();
        initFieldsetToggles();
        initFormValidation();
        initSearch();
    });

    // Select all checkbox functionality
    function initSelectAll() {
        const selectAllCheckboxes = document.querySelectorAll('.select-all');
        selectAllCheckboxes.forEach(function(checkbox) {
            checkbox.addEventListener('change', function() {
                const checked = this.checked;
                const rowCheckboxes = document.querySelectorAll('input[name="selected"]');
                rowCheckboxes.forEach(function(cb) {
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
            goButton.addEventListener('click', function() {
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

    // Get selected items
    function getSelectedItems() {
        const selected = [];
        const checkboxes = document.querySelectorAll('input[name="selected"]:checked');
        checkboxes.forEach(function(cb) {
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

        selected.forEach(function(id) {
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
        toggleButtons.forEach(function(button) {
            button.addEventListener('click', function() {
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
        forms.forEach(function(form) {
            form.addEventListener('submit', function(e) {
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
            searchInput.addEventListener('input', function() {
                clearTimeout(timeout);
                timeout = setTimeout(function() {
                    // Auto-submit on enter or after delay
                }, 500);
            });
        }
    }

    // Utility: Show message
    window.adminShowMessage = function(message, type) {
        type = type || 'info';
        const alertDiv = document.createElement('div');
        alertDiv.className = 'alert alert-' + type;
        alertDiv.textContent = message;
        
        const messagesContainer = document.querySelector('.admin-messages');
        if (messagesContainer) {
            messagesContainer.appendChild(alertDiv);
            setTimeout(function() {
                alertDiv.remove();
            }, 5000);
        }
    };

    // Utility: Confirm delete
    window.adminConfirmDelete = function(url) {
        if (confirm('Are you sure you want to delete this item? This action cannot be undone.')) {
            window.location.href = url;
        }
    };

})();
