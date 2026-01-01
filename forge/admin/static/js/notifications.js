// Forge Admin - Toast Notification System

(function () {
    'use strict';

    const NotificationManager = {
        container: null,
        notifications: [],

        init: function () {
            // Create notification container if it doesn't exist
            if (!this.container) {
                this.container = document.createElement('div');
                this.container.id = 'admin-notifications';
                this.container.className = 'admin-notifications';
                document.body.appendChild(this.container);
            }
        },

        show: function (message, type, duration) {
            this.init();

            type = type || 'info';
            duration = duration || 5000; // Default 5 seconds

            const notification = this.createNotification(message, type, duration);
            this.container.appendChild(notification);
            this.notifications.push(notification);

            // Trigger animation
            setTimeout(function () {
                notification.classList.add('show');
            }, 10);

            // Auto-dismiss
            if (duration > 0) {
                setTimeout(function () {
                    NotificationManager.dismiss(notification);
                }, duration);
            }

            return notification;
        },

        createNotification: function (message, type, duration) {
            const notification = document.createElement('div');
            notification.className = 'admin-notification admin-notification-' + type;
            notification.setAttribute('role', 'alert');

            const icon = this.getIcon(type);
            const closeBtn = '<button type="button" class="notification-close" aria-label="Close">&times;</button>';

            notification.innerHTML = `
                <div class="notification-content">
                    <span class="notification-icon">${icon}</span>
                    <span class="notification-message">${this.escapeHtml(message)}</span>
                </div>
                ${closeBtn}
            `;

            // Close button handler
            const closeButton = notification.querySelector('.notification-close');
            if (closeButton) {
                closeButton.addEventListener('click', function () {
                    NotificationManager.dismiss(notification);
                });
            }

            return notification;
        },

        getIcon: function (type) {
            const icons = {
                success: '✓',
                error: '✕',
                warning: '⚠',
                info: 'ℹ'
            };
            return icons[type] || icons.info;
        },

        dismiss: function (notification) {
            if (!notification || !notification.parentNode) {
                return;
            }

            notification.classList.add('dismissing');
            setTimeout(function () {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
                // Remove from array
                const index = NotificationManager.notifications.indexOf(notification);
                if (index > -1) {
                    NotificationManager.notifications.splice(index, 1);
                }
            }, 300);
        },

        success: function (message, duration) {
            return this.show(message, 'success', duration);
        },

        error: function (message, duration) {
            return this.show(message, 'error', duration);
        },

        warning: function (message, duration) {
            return this.show(message, 'warning', duration);
        },

        info: function (message, duration) {
            return this.show(message, 'info', duration);
        },

        escapeHtml: function (text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        },

        clear: function () {
            this.notifications.forEach(function (notification) {
                NotificationManager.dismiss(notification);
            });
        }
    };

    // Initialize on DOM ready
    document.addEventListener('DOMContentLoaded', function () {
        NotificationManager.init();

        // Check for flash messages from server
        const flashMessages = document.querySelectorAll('[data-flash-message]');
        flashMessages.forEach(function (element) {
            const message = element.getAttribute('data-flash-message');
            const type = element.getAttribute('data-flash-type') || 'info';
            NotificationManager.show(message, type);
            element.remove();
        });
    });

    // Expose to global scope
    window.AdminNotifications = NotificationManager;

})();
