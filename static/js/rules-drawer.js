/**
 * Rules Drawer Module - Slide-out quick reference panel
 * @module rules-drawer
 */

const RulesDrawer = {
    /** @type {HTMLElement|null} */
    drawer: null,
    /** @type {HTMLElement|null} */
    overlay: null,

    /**
     * Initialize the rules drawer
     */
    init() {
        this.drawer = document.getElementById('rules-drawer');
        this.overlay = document.getElementById('rules-drawer-overlay');

        // Close on escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isOpen()) {
                this.close();
            }
        });
    },

    /**
     * Check if drawer is open
     * @returns {boolean}
     */
    isOpen() {
        return this.drawer && this.drawer.classList.contains('open');
    },

    /**
     * Open the drawer
     */
    open() {
        if (this.drawer && this.overlay) {
            this.drawer.classList.add('open');
            this.overlay.classList.add('active');
            document.body.style.overflow = 'hidden';
        }
    },

    /**
     * Close the drawer
     */
    close() {
        if (this.drawer && this.overlay) {
            this.drawer.classList.remove('open');
            this.overlay.classList.remove('active');
            document.body.style.overflow = '';
        }
    },

    /**
     * Toggle the drawer open/closed
     */
    toggle() {
        if (this.isOpen()) {
            this.close();
        } else {
            this.open();
        }
    },

    /**
     * Switch to a specific tab
     * @param {string} tabName - The tab to switch to (insanity, combat, weapons)
     */
    switchTab(tabName) {
        // Update tab buttons
        const tabs = document.querySelectorAll('.rules-tab');
        tabs.forEach(tab => {
            if (tab.dataset.tab === tabName) {
                tab.classList.add('active');
            } else {
                tab.classList.remove('active');
            }
        });

        // Update tab panels
        const panels = document.querySelectorAll('.rules-panel');
        panels.forEach(panel => {
            if (panel.id === `tab-${tabName}`) {
                panel.classList.add('active');
            } else {
                panel.classList.remove('active');
            }
        });
    }
};

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    RulesDrawer.init();
});

// Re-initialize after HTMX swaps (in case drawer is in swapped content)
document.addEventListener('htmx:afterSettle', () => {
    RulesDrawer.init();
});

// Export globally
window.RulesDrawer = RulesDrawer;
