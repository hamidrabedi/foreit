// Forge Admin - Rich Text Editor JavaScript

(function() {
    'use strict';

    // Initialize rich text editors
    document.addEventListener('DOMContentLoaded', function() {
        initRichTextEditors();
    });

    // Initialize rich text editors
    function initRichTextEditors() {
        const editors = document.querySelectorAll('.richtext-editor');
        editors.forEach(function(textarea) {
            initRichTextEditor(textarea);
        });
    }

    // Initialize a single rich text editor
    function initRichTextEditor(textarea) {
        const editorId = textarea.id;
        const toolbar = textarea.dataset.toolbar || 'basic';
        const height = parseInt(textarea.dataset.height || '300');

        // Simple rich text editor implementation
        // In production, would integrate with a library like TinyMCE, CKEditor, or Quill
        
        // Create toolbar
        const toolbarDiv = document.createElement('div');
        toolbarDiv.className = 'richtext-toolbar';
        toolbarDiv.innerHTML = getToolbarHTML(toolbar);
        
        // Insert toolbar before textarea
        textarea.parentNode.insertBefore(toolbarDiv, textarea);
        
        // Add toolbar event listeners
        setupToolbarEvents(toolbarDiv, textarea);
        
        // Set height
        textarea.style.height = height + 'px';
    }

    // Get toolbar HTML based on configuration
    function getToolbarHTML(toolbar) {
        const toolbars = {
            'full': `
                <button type="button" class="toolbar-btn" data-command="bold" title="Bold">B</button>
                <button type="button" class="toolbar-btn" data-command="italic" title="Italic">I</button>
                <button type="button" class="toolbar-btn" data-command="underline" title="Underline">U</button>
                <span class="toolbar-separator"></span>
                <button type="button" class="toolbar-btn" data-command="formatBlock" data-value="h1" title="Heading 1">H1</button>
                <button type="button" class="toolbar-btn" data-command="formatBlock" data-value="h2" title="Heading 2">H2</button>
                <button type="button" class="toolbar-btn" data-command="formatBlock" data-value="p" title="Paragraph">P</button>
                <span class="toolbar-separator"></span>
                <button type="button" class="toolbar-btn" data-command="insertUnorderedList" title="Bullet List">•</button>
                <button type="button" class="toolbar-btn" data-command="insertOrderedList" title="Numbered List">1.</button>
                <span class="toolbar-separator"></span>
                <button type="button" class="toolbar-btn" data-command="createLink" title="Link">🔗</button>
            `,
            'basic': `
                <button type="button" class="toolbar-btn" data-command="bold" title="Bold">B</button>
                <button type="button" class="toolbar-btn" data-command="italic" title="Italic">I</button>
                <button type="button" class="toolbar-btn" data-command="underline" title="Underline">U</button>
            `,
            'minimal': `
                <button type="button" class="toolbar-btn" data-command="bold" title="Bold">B</button>
                <button type="button" class="toolbar-btn" data-command="italic" title="Italic">I</button>
            `
        };
        
        return toolbars[toolbar] || toolbars['basic'];
    }

    // Setup toolbar event listeners
    function setupToolbarEvents(toolbar, textarea) {
        const buttons = toolbar.querySelectorAll('.toolbar-btn');
        buttons.forEach(function(button) {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                const command = button.dataset.command;
                const value = button.dataset.value;
                
                // Focus textarea
                textarea.focus();
                
                // Execute command
                if (command === 'createLink') {
                    const url = prompt('Enter URL:');
                    if (url) {
                        document.execCommand('createLink', false, url);
                    }
                } else if (value) {
                    document.execCommand(command, false, value);
                } else {
                    document.execCommand(command, false, null);
                }
            });
        });
    }

    // Make function available globally
    window.initRichTextEditor = function(editorId) {
        const textarea = document.getElementById(editorId);
        if (textarea) {
            initRichTextEditor(textarea);
        }
    };

})();
