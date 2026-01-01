// Forge Admin - Dark Mode Toggle

(function () {
    'use strict';

    const DarkMode = {
        storageKey: 'forge_admin_theme',
        theme: 'light',

        init: function () {
            // Load saved theme or detect system preference
            const saved = localStorage.getItem(this.storageKey);
            if (saved) {
                this.theme = saved;
            } else {
                // Detect system preference
                if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                    this.theme = 'dark';
                }
            }

            this.applyTheme();
            this.setupToggle();
        },

        applyTheme: function () {
            const root = document.documentElement;
            if (this.theme === 'dark') {
                root.setAttribute('data-theme', 'dark');
                root.classList.add('dark-mode');
            } else {
                root.removeAttribute('data-theme');
                root.classList.remove('dark-mode');
            }
        },

        toggle: function () {
            this.theme = this.theme === 'dark' ? 'light' : 'dark';
            localStorage.setItem(this.storageKey, this.theme);
            this.applyTheme();
        },

        setupToggle: function () {
            // Create toggle button if it doesn't exist
            let toggle = document.getElementById('dark-mode-toggle');
            if (!toggle) {
                const header = document.querySelector('.admin-header-content');
                if (header) {
                    toggle = document.createElement('button');
                    toggle.id = 'dark-mode-toggle';
                    toggle.className = 'admin-dark-mode-toggle';
                    toggle.setAttribute('aria-label', 'Toggle dark mode');
                    toggle.innerHTML = this.theme === 'dark' ? '☀' : '🌙';
                    toggle.addEventListener('click', () => {
                        this.toggle();
                        toggle.innerHTML = this.theme === 'dark' ? '☀' : '🌙';
                    });
                    header.appendChild(toggle);
                }
            }
        }
    };

    // Initialize on DOM ready
    document.addEventListener('DOMContentLoaded', function () {
        DarkMode.init();
    });

    // Watch for system theme changes
    if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
            if (!localStorage.getItem(DarkMode.storageKey)) {
                DarkMode.theme = e.matches ? 'dark' : 'light';
                DarkMode.applyTheme();
            }
        });
    }

    // Expose to global scope
    window.DarkMode = DarkMode;

})();
