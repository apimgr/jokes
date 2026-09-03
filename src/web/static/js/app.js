// ========================================
// Jokes API - Professional UI with Modals
// ========================================

// Modal System
class Modal {
    constructor() {
        this.modals = new Map();
    }

    create(id, title, content, options = {}) {
        const modal = document.createElement('div');
        modal.className = 'modal';
        modal.id = id;
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h2>${title}</h2>
                    <button class="modal-close" aria-label="Close">&times;</button>
                </div>
                <div class="modal-body">${content}</div>
                ${options.footer ? `<div class="modal-footer">${options.footer}</div>` : ''}
            </div>
        `;

        document.body.appendChild(modal);
        this.modals.set(id, modal);

        // Close handlers
        modal.querySelector('.modal-close').addEventListener('click', () => this.close(id));
        modal.addEventListener('click', (e) => {
            if (e.target === modal) this.close(id);
        });

        return modal;
    }

    show(id) {
        const modal = this.modals.get(id);
        if (modal) {
            modal.classList.add('active');
            document.body.style.overflow = 'hidden';
        }
    }

    close(id) {
        const modal = this.modals.get(id);
        if (modal) {
            modal.classList.remove('active');
            document.body.style.overflow = '';
        }
    }

    destroy(id) {
        const modal = this.modals.get(id);
        if (modal) {
            modal.remove();
            this.modals.delete(id);
        }
    }
}

// Toast Notification System
class Toast {
    constructor() {
        this.toasts = [];
    }

    show(message, type = 'info', duration = 3000) {
        const icons = {
            success: '✅',
            error: '❌',
            warning: '⚠️',
            info: 'ℹ️'
        };

        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.innerHTML = `
            <span class="toast-icon">${icons[type] || icons.info}</span>
            <span class="toast-message">${message}</span>
            <button class="toast-close">&times;</button>
        `;

        document.body.appendChild(toast);
        this.toasts.push(toast);

        const closeBtn = toast.querySelector('.toast-close');
        closeBtn.addEventListener('click', () => this.hide(toast));

        if (duration > 0) {
            setTimeout(() => this.hide(toast), duration);
        }

        return toast;
    }

    hide(toast) {
        toast.style.animation = 'slideOutRight 0.3s ease';
        setTimeout(() => {
            toast.remove();
            this.toasts = this.toasts.filter(t => t !== toast);
        }, 300);
    }
}

// Initialize systems
const modalSystem = new Modal();
const toastSystem = new Toast();

// ========================================
// Theme Management
// ========================================
const themeToggle = document.getElementById('theme-toggle');
const themeIcon = document.querySelector('.theme-icon');
const body = document.body;

const savedTheme = localStorage.getItem('theme') || 'dark';
body.setAttribute('data-theme', savedTheme);
updateThemeIcon(savedTheme);

themeToggle?.addEventListener('click', () => {
    const currentTheme = body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

    body.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);

    toastSystem.show(`Switched to ${newTheme} theme`, 'success', 2000);
});

function updateThemeIcon(theme) {
    if (themeIcon) {
        themeIcon.textContent = theme === 'dark' ? '☀️' : '🌙';
    }
}

// ========================================
// Mobile Menu
// ========================================
const mobileMenuToggle = document.querySelector('.mobile-menu-toggle');
const navLinks = document.querySelector('.nav-links');

mobileMenuToggle?.addEventListener('click', () => {
    navLinks?.classList.toggle('active');
});

// Close mobile menu when clicking outside
document.addEventListener('click', (e) => {
    if (navLinks?.classList.contains('active') &&
        !navLinks.contains(e.target) &&
        !mobileMenuToggle?.contains(e.target)) {
        navLinks.classList.remove('active');
    }
});

// ========================================
// Get New Joke
// ========================================
const newJokeBtn = document.getElementById('new-joke-btn');
const jokeContent = document.getElementById('joke-content');

newJokeBtn?.addEventListener('click', async () => {
    try {
        newJokeBtn.disabled = true;
        newJokeBtn.innerHTML = '<span class="loading"></span> Loading...';

        const response = await fetch('/api/v1/jokes/random');
        const data = await response.json();

        if (data.type === 'success' && data.value) {
            const joke = data.value;
            jokeContent.innerHTML = `
                <p class="joke-text">${escapeHtml(joke.joke)}</p>
                <div class="joke-meta">
                    <span class="joke-id">ID: ${joke.id}</span>
                    <div class="joke-categories">
                        ${joke.categories.map(cat => `<span class="category-badge">${escapeHtml(cat)}</span>`).join('')}
                    </div>
                </div>
            `;
            toastSystem.show('New joke loaded!', 'success', 2000);
        } else {
            throw new Error('Invalid response');
        }
    } catch (error) {
        console.error('Error fetching joke:', error);
        toastSystem.show('Failed to load joke. Please try again.', 'error');
        jokeContent.innerHTML = `
            <div class="message error">
                <span>❌</span>
                <span>Failed to load joke. Please try again.</span>
            </div>
        `;
    } finally {
        newJokeBtn.disabled = false;
        newJokeBtn.textContent = 'Get New Joke';
    }
});

// ========================================
// Copy Code Example
// ========================================
const copyBtns = document.querySelectorAll('.copy-btn');

copyBtns.forEach(btn => {
    btn.addEventListener('click', async () => {
        const targetId = btn.getAttribute('data-clipboard-target');
        const target = document.querySelector(targetId);

        if (target) {
            try {
                await navigator.clipboard.writeText(target.textContent);
                const originalText = btn.textContent;
                btn.textContent = '✅ Copied!';
                toastSystem.show('Code copied to clipboard!', 'success', 2000);

                setTimeout(() => {
                    btn.textContent = originalText;
                }, 2000);
            } catch (error) {
                console.error('Failed to copy:', error);
                toastSystem.show('Failed to copy code', 'error');
            }
        }
    });
});

// ========================================
// Notification System
// ========================================
const notificationBtn = document.getElementById('notification-btn');
const notificationBadge = document.getElementById('notification-badge');
let notificationPanel = null;

notificationBtn?.addEventListener('click', () => {
    if (!notificationPanel) {
        createNotificationPanel();
    }
    notificationPanel.classList.toggle('active');
});

function createNotificationPanel() {
    notificationPanel = document.createElement('div');
    notificationPanel.className = 'notification-panel';
    notificationPanel.innerHTML = `
        <div class="notification-panel-header">
            <h3>🔔 Notifications</h3>
            <button class="modal-close" aria-label="Close">&times;</button>
        </div>
        <div class="notification-list">
            <div class="notification-item">
                <strong>Welcome!</strong>
                <p>Enjoy browsing 5,160+ jokes across 16 categories.</p>
            </div>
        </div>
    `;

    document.body.appendChild(notificationPanel);

    notificationPanel.querySelector('.modal-close').addEventListener('click', () => {
        notificationPanel.classList.remove('active');
    });
}

function showNotification(count) {
    if (notificationBadge && count > 0) {
        notificationBadge.textContent = count;
        notificationBadge.style.display = 'flex';
    } else if (notificationBadge) {
        notificationBadge.style.display = 'none';
    }
}

// ========================================
// Service Worker (PWA)
// ========================================
if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/sw.js')
            .then(registration => {
                console.log('✅ ServiceWorker registered:', registration.scope);
            })
            .catch(error => {
                console.log('❌ ServiceWorker registration failed:', error);
            });
    });
}

// ========================================
// API Health Check
// ========================================
async function checkAPIHealth() {
    try {
        const response = await fetch('/healthz');
        if (response.ok) {
            console.log('✅ API is healthy');
        }
    } catch (error) {
        console.error('❌ API health check failed:', error);
        toastSystem.show('API health check failed', 'warning', 5000);
    }
}

checkAPIHealth();

// ========================================
// Online/Offline Status
// ========================================
window.addEventListener('online', () => {
    console.log('✅ Back online');
    toastSystem.show('You are back online!', 'success');
    showNotification(0);
});

window.addEventListener('offline', () => {
    console.log('⚠️ Offline mode');
    toastSystem.show('You are offline. Some features may be limited.', 'warning');
    showNotification(1);
});

// ========================================
// Smooth Scrolling
// ========================================
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
    });
});

// ========================================
// Utility Functions
// ========================================
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showLoader() {
    const overlay = document.createElement('div');
    overlay.className = 'spinner-overlay';
    overlay.innerHTML = '<div class="spinner"></div>';
    overlay.id = 'global-loader';
    document.body.appendChild(overlay);
    return overlay;
}

function hideLoader() {
    const loader = document.getElementById('global-loader');
    if (loader) {
        loader.remove();
    }
}

// ========================================
// Keyboard Shortcuts
// ========================================
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + K: Focus search (if exists)
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        // Add search functionality
    }

    // Escape: Close modals/panels
    if (e.key === 'Escape') {
        modalSystem.modals.forEach((modal, id) => {
            if (modal.classList.contains('active')) {
                modalSystem.close(id);
            }
        });
        if (notificationPanel?.classList.contains('active')) {
            notificationPanel.classList.remove('active');
        }
        if (navLinks?.classList.contains('active')) {
            navLinks.classList.remove('active');
        }
    }
});

// ========================================
// Initialize Tooltips
// ========================================
document.querySelectorAll('[data-tooltip]').forEach(element => {
    const tooltipText = element.getAttribute('data-tooltip');
    element.classList.add('tooltip');
    const tooltip = document.createElement('span');
    tooltip.className = 'tooltip-text';
    tooltip.textContent = tooltipText;
    element.appendChild(tooltip);
});

console.log('🎭 Jokes API - Ready!');
