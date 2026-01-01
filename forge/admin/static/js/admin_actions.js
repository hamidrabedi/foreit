// Forge Admin - Action Confirmation and Forms

(function () {
    'use strict';

    let currentAction = null;
    let currentSelected = [];
    let actionFormFields = null;

    // Initialize action system
    document.addEventListener('DOMContentLoaded', function () {
        initActionModals();
        initBulkActionConfirmations();
    });

    // Initialize action modals
    function initActionModals() {
        // Create confirmation modal if it doesn't exist
        if (!document.getElementById('action-confirmation-modal')) {
            createConfirmationModal();
        }

        // Create form modal if it doesn't exist
        if (!document.getElementById('action-form-modal')) {
            createFormModal();
        }

        // Setup modal close handlers
        setupModalCloseHandlers();
    }

    // Create confirmation modal
    function createConfirmationModal() {
        const modal = document.createElement('div');
        modal.id = 'action-confirmation-modal';
        modal.className = 'admin-modal';
        modal.style.display = 'none';
        modal.innerHTML = `
            <div class="admin-modal-overlay"></div>
            <div class="admin-modal-content">
                <div class="admin-modal-header">
                    <h3 class="admin-modal-title">Confirm Action</h3>
                    <button type="button" class="admin-modal-close" aria-label="Close">&times;</button>
                </div>
                <div class="admin-modal-body">
                    <p class="action-confirmation-message"></p>
                    <div class="action-confirmation-details">
                        <p><strong>Action:</strong> <span class="action-name"></span></p>
                        <p><strong>Items selected:</strong> <span class="action-count"></span></p>
                    </div>
                </div>
                <div class="admin-modal-footer">
                    <button type="button" class="btn btn-secondary action-cancel">Cancel</button>
                    <button type="button" class="btn btn-primary action-confirm">Confirm</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }

    // Create form modal
    function createFormModal() {
        const modal = document.createElement('div');
        modal.id = 'action-form-modal';
        modal.className = 'admin-modal';
        modal.style.display = 'none';
        modal.innerHTML = `
            <div class="admin-modal-overlay"></div>
            <div class="admin-modal-content admin-modal-large">
                <div class="admin-modal-header">
                    <h3 class="admin-modal-title">Action Parameters</h3>
                    <button type="button" class="admin-modal-close" aria-label="Close">&times;</button>
                </div>
                <form id="action-form" class="admin-modal-body">
                    <div class="action-form-fields"></div>
                    <div class="action-form-info">
                        <p><strong>Action:</strong> <span class="action-name"></span></p>
                        <p><strong>Items selected:</strong> <span class="action-count"></span></p>
                    </div>
                </form>
                <div class="admin-modal-footer">
                    <button type="button" class="btn btn-secondary action-cancel">Cancel</button>
                    <button type="submit" form="action-form" class="btn btn-primary action-submit">Execute</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }

    // Setup modal close handlers
    function setupModalCloseHandlers() {
        // Close buttons
        document.querySelectorAll('.admin-modal-close, .action-cancel').forEach(function (btn) {
            btn.addEventListener('click', function () {
                closeModals();
            });
        });

        // Overlay click
        document.querySelectorAll('.admin-modal-overlay').forEach(function (overlay) {
            overlay.addEventListener('click', function () {
                closeModals();
            });
        });

        // ESC key
        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') {
                closeModals();
            }
        });
    }

    // Initialize bulk action confirmations
    function initBulkActionConfirmations() {
        const bulkActionForm = document.querySelector('.admin-bulk-actions');
        if (!bulkActionForm) return;

        const actionSelect = bulkActionForm.querySelector('.action-select');
        const goButton = bulkActionForm.querySelector('button[type="button"]');

        if (goButton && actionSelect) {
            goButton.addEventListener('click', function (e) {
                e.preventDefault();

                const action = actionSelect.value;
                if (!action) {
                    showNotification('Please select an action.', 'error');
                    return;
                }

                const selected = getSelectedItems();
                if (selected.length === 0) {
                    showNotification('Please select at least one item.', 'error');
                    return;
                }

                // Get action metadata
                const actionOption = actionSelect.options[actionSelect.selectedIndex];
                const actionLabel = actionOption.text;
                const requiresForm = actionOption.dataset.requiresForm === 'true';
                const requiresConfirmation = actionOption.dataset.requiresConfirmation !== 'false';

                currentAction = action;
                currentSelected = selected;

                if (requiresForm) {
                    showActionForm(action, actionLabel, selected);
                } else if (requiresConfirmation) {
                    showActionConfirmation(action, actionLabel, selected);
                } else {
                    executeBulkAction(action, selected);
                }
            });
        }
    }

    // Show action confirmation
    function showActionConfirmation(action, actionLabel, selected) {
        const modal = document.getElementById('action-confirmation-modal');
        if (!modal) return;

        const message = modal.querySelector('.action-confirmation-message');
        const actionName = modal.querySelector('.action-name');
        const actionCount = modal.querySelector('.action-count');
        const confirmBtn = modal.querySelector('.action-confirm');

        message.textContent = `Are you sure you want to perform "${actionLabel}" on ${selected.length} item(s)?`;
        actionName.textContent = actionLabel;
        actionCount.textContent = selected.length;

        // Remove old handler and add new one
        const newConfirmBtn = confirmBtn.cloneNode(true);
        confirmBtn.parentNode.replaceChild(newConfirmBtn, confirmBtn);
        newConfirmBtn.addEventListener('click', function () {
            closeModals();
            executeBulkAction(action, selected);
        });

        modal.style.display = 'block';
    }

    // Show action form
    function showActionForm(action, actionLabel, selected) {
        const modal = document.getElementById('action-form-modal');
        if (!modal) return;

        const actionName = modal.querySelector('.action-form-info .action-name');
        const actionCount = modal.querySelector('.action-form-info .action-count');
        const formFields = modal.querySelector('.action-form-fields');
        const form = document.getElementById('action-form');

        actionName.textContent = actionLabel;
        actionCount.textContent = selected.length;

        // Load form fields (this would come from the server)
        // For now, we'll use a simple example
        formFields.innerHTML = `
            <div class="form-field">
                <label class="form-label">Reason (optional)</label>
                <input type="text" name="reason" class="form-control" placeholder="Enter reason for this action">
            </div>
        `;

        // Setup form submission
        form.onsubmit = function (e) {
            e.preventDefault();
            const formData = new FormData(form);
            closeModals();
            executeBulkActionWithForm(action, selected, formData);
        };

        modal.style.display = 'block';
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

    // Execute bulk action with form data
    function executeBulkActionWithForm(action, selected, formData) {
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

        // Add form data
        for (const [key, value] of formData.entries()) {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = key;
            input.value = value;
            form.appendChild(input);
        }

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

    // Close all modals
    function closeModals() {
        document.querySelectorAll('.admin-modal').forEach(function (modal) {
            modal.style.display = 'none';
        });
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

    // Show notification using toast system
    function showNotification(message, type) {
        if (window.AdminNotifications) {
            window.AdminNotifications.show(message, type);
        } else {
            alert(message); // Fallback
        }
    }

})();
