// Forge Admin - Column Management

(function () {
    'use strict';

    const ColumnManager = {
        storageKey: 'forge_admin_column_prefs_',
        currentModel: null,

        init: function (modelName) {
            this.currentModel = modelName;
            this.loadPreferences();
            this.setupColumnToggles();
            this.setupColumnReordering();
        },

        loadPreferences: function () {
            if (!this.currentModel) return;

            const key = this.storageKey + this.currentModel;
            const prefs = localStorage.getItem(key);
            if (!prefs) return;

            try {
                const preferences = JSON.parse(prefs);
                
                // Apply visibility
                if (preferences.visible) {
                    this.applyVisibility(preferences.visible);
                }

                // Apply order
                if (preferences.order) {
                    this.applyOrder(preferences.order);
                }
            } catch (e) {
                console.error('Failed to load column preferences:', e);
            }
        },

        savePreferences: function () {
            if (!this.currentModel) return;

            const key = this.storageKey + this.currentModel;
            const preferences = {
                visible: this.getVisibleColumns(),
                order: this.getColumnOrder()
            };

            try {
                localStorage.setItem(key, JSON.stringify(preferences));
            } catch (e) {
                console.error('Failed to save column preferences:', e);
            }
        },

        getVisibleColumns: function () {
            const visible = [];
            const headers = document.querySelectorAll('.admin-table thead th:not(.select-column):not(.actions-column)');
            headers.forEach(function (th, index) {
                if (th.style.display !== 'none') {
                    const fieldName = th.getAttribute('data-field') || th.textContent.trim();
                    visible.push(fieldName);
                }
            });
            return visible;
        },

        getColumnOrder: function () {
            const order = [];
            const headers = document.querySelectorAll('.admin-table thead th:not(.select-column):not(.actions-column)');
            headers.forEach(function (th) {
                const fieldName = th.getAttribute('data-field') || th.textContent.trim();
                order.push(fieldName);
            });
            return order;
        },

        applyVisibility: function (visibleColumns) {
            const headers = document.querySelectorAll('.admin-table thead th:not(.select-column):not(.actions-column)');
            const cells = document.querySelectorAll('.admin-table tbody td:not(.select-column):not(.actions-column)');
            
            headers.forEach(function (th, index) {
                const fieldName = th.getAttribute('data-field') || th.textContent.trim();
                const isVisible = visibleColumns.includes(fieldName);
                
                // Toggle header
                th.style.display = isVisible ? '' : 'none';
                
                // Toggle corresponding cells in each row
                const rowCount = document.querySelectorAll('.admin-table tbody tr').length;
                for (let i = 0; i < rowCount; i++) {
                    const cellIndex = i * (headers.length) + index;
                    if (cells[cellIndex]) {
                        cells[cellIndex].style.display = isVisible ? '' : 'none';
                    }
                }
            });
        },

        applyOrder: function (columnOrder) {
            // This is complex - would need to reorder DOM elements
            // For now, we'll just track the order
            // Full implementation would require reordering table cells
            console.log('Column order:', columnOrder);
        },

        setupColumnToggles: function () {
            // Create column management button
            const listHeader = document.querySelector('.admin-list-header');
            if (!listHeader) return;

            // Check if button already exists
            if (document.getElementById('column-manager-toggle')) return;

            const toggleBtn = document.createElement('button');
            toggleBtn.id = 'column-manager-toggle';
            toggleBtn.className = 'btn btn-secondary';
            toggleBtn.innerHTML = '⚙ Columns';
            toggleBtn.type = 'button';
            toggleBtn.addEventListener('click', function () {
                ColumnManager.showColumnManager();
            });

            listHeader.appendChild(toggleBtn);
        },

        showColumnManager: function () {
            // Create or show column manager modal
            let modal = document.getElementById('column-manager-modal');
            if (!modal) {
                modal = this.createColumnManagerModal();
                document.body.appendChild(modal);
            }

            this.populateColumnManager();
            modal.style.display = 'block';
        },

        createColumnManagerModal: function () {
            const modal = document.createElement('div');
            modal.id = 'column-manager-modal';
            modal.className = 'admin-modal';
            modal.style.display = 'none';
            modal.innerHTML = `
                <div class="admin-modal-overlay"></div>
                <div class="admin-modal-content">
                    <div class="admin-modal-header">
                        <h3 class="admin-modal-title">Manage Columns</h3>
                        <button type="button" class="admin-modal-close" aria-label="Close">&times;</button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="column-manager-list"></div>
                    </div>
                    <div class="admin-modal-footer">
                        <button type="button" class="btn btn-secondary column-manager-reset">Reset to Default</button>
                        <button type="button" class="btn btn-primary column-manager-save">Save</button>
                    </div>
                </div>
            `;

            // Setup close handlers
            modal.querySelector('.admin-modal-close').addEventListener('click', function () {
                modal.style.display = 'none';
            });
            modal.querySelector('.admin-modal-overlay').addEventListener('click', function () {
                modal.style.display = 'none';
            });

            // Setup save handler
            modal.querySelector('.column-manager-save').addEventListener('click', function () {
                ColumnManager.saveColumnPreferences();
                modal.style.display = 'none';
            });

            // Setup reset handler
            modal.querySelector('.column-manager-reset').addEventListener('click', function () {
                ColumnManager.resetToDefault();
                modal.style.display = 'none';
            });

            return modal;
        },

        populateColumnManager: function () {
            const list = document.querySelector('#column-manager-modal .column-manager-list');
            if (!list) return;

            list.innerHTML = '';

            const headers = document.querySelectorAll('.admin-table thead th:not(.select-column):not(.actions-column)');
            headers.forEach(function (th, index) {
                const fieldName = th.getAttribute('data-field') || th.textContent.trim();
                const label = th.textContent.trim();
                const isVisible = th.style.display !== 'none';

                const item = document.createElement('div');
                item.className = 'column-manager-item';
                item.innerHTML = `
                    <label class="column-toggle-label">
                        <input type="checkbox" class="column-toggle" data-field="${this.escapeHtml(fieldName)}" ${isVisible ? 'checked' : ''}>
                        <span class="column-label">${this.escapeHtml(label)}</span>
                        <span class="column-drag-handle" data-field="${this.escapeHtml(fieldName)}">☰</span>
                    </label>
                `;
                list.appendChild(item);
            }.bind(this));
        },

        saveColumnPreferences: function () {
            const checkboxes = document.querySelectorAll('#column-manager-modal .column-toggle');
            const visible = [];
            
            checkboxes.forEach(function (cb) {
                if (cb.checked) {
                    visible.push(cb.getAttribute('data-field'));
                }
            });

            this.applyVisibility(visible);
            this.savePreferences();
        },

        resetToDefault: function () {
            if (!this.currentModel) return;

            const key = this.storageKey + this.currentModel;
            localStorage.removeItem(key);

            // Show all columns
            const headers = document.querySelectorAll('.admin-table thead th:not(.select-column):not(.actions-column)');
            const cells = document.querySelectorAll('.admin-table tbody td:not(.select-column):not(.actions-column)');
            
            headers.forEach(function (th) {
                th.style.display = '';
            });
            
            cells.forEach(function (td) {
                td.style.display = '';
            });
        },

        setupColumnReordering: function () {
            // This would use a drag-and-drop library like SortableJS
            // For now, we'll just track the functionality
            // Full implementation would require:
            // 1. Making column headers draggable
            // 2. Reordering table cells when headers are moved
            // 3. Saving the new order
        },

        escapeHtml: function (text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    };

    // Initialize on DOM ready
    document.addEventListener('DOMContentLoaded', function () {
        // Extract model name from URL or page
        const path = window.location.pathname;
        const match = path.match(/\/admin\/([^\/]+)/);
        if (match) {
            ColumnManager.init(match[1]);
        }
    });

    // Expose to global scope
    window.ColumnManager = ColumnManager;

})();
