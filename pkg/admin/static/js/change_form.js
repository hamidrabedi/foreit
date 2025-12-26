'use strict';
{
    const inputTags = ['BUTTON', 'INPUT', 'SELECT', 'TEXTAREA'];
    const formAddConstants = document.getElementById('forge-admin-form-add-constants');
    
    if (formAddConstants) {
        const modelName = formAddConstants.dataset.modelName;
        if (modelName) {
            const form = document.getElementById(modelName + '_form');
            if (form) {
                for (const element of form.elements) {
                    // HTMLElement.offsetParent returns null when the element is not
                    // rendered.
                    if (inputTags.includes(element.tagName) && !element.disabled && element.offsetParent) {
                        element.focus();
                        break;
                    }
                }
            }
        }
    }
}

