// Forge Admin - Inline Forms JavaScript

(function () {
    'use strict';

    // Initialize inline forms
    document.addEventListener('DOMContentLoaded', function () {
        initInlineForms();
    });

    // Initialize inline forms
    function initInlineForms() {
        const inlineGroups = document.querySelectorAll('.inline-group');
        inlineGroups.forEach(function (group) {
            const addButton = group.querySelector('.btn-add-inline');
            if (addButton) {
                addButton.addEventListener('click', function () {
                    addInlineForm(group);
                });
            }
        });
    }

    // Add inline form
    function addInlineForm(group) {
        const formsContainer = group.querySelector('.inline-forms, .inline-table-container tbody');
        if (!formsContainer) return;

        const emptyForm = formsContainer.querySelector('.empty-form');
        if (!emptyForm) return;

        // Clone empty form
        const newForm = emptyForm.cloneNode(true);
        newForm.classList.remove('empty-form');

        // Update form indices
        updateFormIndices(newForm, formsContainer.children.length);

        // Clear values
        const inputs = newForm.querySelectorAll('input, select, textarea');
        inputs.forEach(function (input) {
            if (input.type !== 'checkbox' && input.type !== 'hidden') {
                input.value = '';
            }
        });

        // Insert before empty form
        formsContainer.insertBefore(newForm, emptyForm);

        // Check max forms
        checkMaxForms(group);
    }

    // Update form indices in name attributes
    function updateFormIndices(form, index) {
        const inputs = form.querySelectorAll('input, select, textarea');
        inputs.forEach(function (input) {
            const name = input.name;
            if (name) {
                // Update index in name (e.g., "comment-0-title" -> "comment-1-title")
                input.name = name.replace(/-(\d+)-/, '-' + index + '-');
            }
        });
    }

    // Check maximum number of forms
    function checkMaxForms(group) {
        const maxNum = parseInt(group.dataset.maxNum || '0');
        if (maxNum === 0) return;

        const forms = group.querySelectorAll('.inline-form-row, .inline-form:not(.empty-form)');
        const addButton = group.querySelector('.btn-add-inline');

        if (forms.length >= maxNum && addButton) {
            addButton.style.display = 'none';
        } else if (addButton) {
            addButton.style.display = '';
        }
    }

    // Remove inline form
    function removeInlineForm(formRow) {
        const group = formRow.closest('.inline-group');
        formRow.remove();

        // Update indices
        const forms = group.querySelectorAll('.inline-form-row, .inline-form:not(.empty-form)');
        forms.forEach(function (form, index) {
            updateFormIndices(form, index);
        });

        // Show add button if needed
        checkMaxForms(group);
    }

    // Delete checkbox handler
    document.addEventListener('change', function (e) {
        if (e.target.type === 'checkbox' && e.target.name.includes('-DELETE')) {
            const formRow = e.target.closest('.inline-form-row, .inline-form');
            if (e.target.checked && formRow) {
                formRow.classList.add('deleted');
            } else if (formRow) {
                formRow.classList.remove('deleted');
            }
        }
    });

    // Make functions available globally
    window.addInlineForm = function (groupName) {
        const group = document.querySelector(`[data-inline-name="${groupName}"]`);
        if (group) {
            addInlineForm(group);
        }
    };

})();
