/**
 * Utilities Module - Shared helper functions
 * @module utils
 */

const Utils = {
    // =========================================================================
    // DOM Utilities
    // =========================================================================

    /**
     * Get element by ID with null check
     * @param {string} id - Element ID
     * @returns {HTMLElement|null}
     */
    $(id) {
        return document.getElementById(id);
    },

    /**
     * Query selector shorthand
     * @param {string} selector - CSS selector
     * @param {HTMLElement} parent - Parent element (default: document)
     * @returns {HTMLElement|null}
     */
    qs(selector, parent = document) {
        return parent.querySelector(selector);
    },

    /**
     * Query selector all shorthand
     * @param {string} selector - CSS selector
     * @param {HTMLElement} parent - Parent element (default: document)
     * @returns {NodeList}
     */
    qsa(selector, parent = document) {
        return parent.querySelectorAll(selector);
    },

    /**
     * Safely get element value
     * @param {string} id - Element ID
     * @returns {string}
     */
    getValue(id) {
        const el = this.$(id);
        return el ? el.value : '';
    },

    /**
     * Safely set element text content
     * @param {string} id - Element ID
     * @param {string} text - Text content
     */
    setText(id, text) {
        const el = this.$(id);
        if (el) el.textContent = text;
    },

    /**
     * Safely set element innerHTML
     * @param {string} id - Element ID
     * @param {string} html - HTML content
     */
    setHTML(id, html) {
        const el = this.$(id);
        if (el) el.innerHTML = html;
    },

    // =========================================================================
    // Dice Rolling
    // =========================================================================

    /**
     * Roll a single die
     * @param {number} sides - Number of sides
     * @returns {number}
     */
    rollDie(sides) {
        return Math.floor(Math.random() * sides) + 1;
    },

    /**
     * Roll multiple dice and sum
     * @param {number} count - Number of dice
     * @param {number} sides - Number of sides per die
     * @returns {number}
     */
    rollDice(count, sides) {
        let total = 0;
        for (let i = 0; i < count; i++) {
            total += this.rollDie(sides);
        }
        return total;
    },

    /**
     * Roll attribute using formula
     * @param {string} formula - Formula like '3d6x5' or '2d6p6x5'
     * @returns {number}
     */
    rollAttribute(formula) {
        if (formula === '3d6x5') {
            return this.rollDice(3, 6) * 5;
        } else if (formula === '2d6p6x5') {
            return (this.rollDice(2, 6) + 6) * 5;
        }
        return 0;
    },

    // =========================================================================
    // Value Calculations
    // =========================================================================

    /**
     * Calculate half value (floor)
     * @param {number} value - Original value
     * @returns {number}
     */
    half(value) {
        return Math.floor(value / 2);
    },

    /**
     * Calculate fifth value (floor)
     * @param {number} value - Original value
     * @returns {number}
     */
    fifth(value) {
        return Math.floor(value / 5);
    },

    /**
     * Update derived values (half/fifth) in a container
     * @param {HTMLElement} container - Container element
     * @param {number} value - Base value
     */
    updateDerivedValues(container, value) {
        const halfSpan = container.querySelector('[data-half], .half-value, .attr-half');
        const fifthSpan = container.querySelector('[data-fifth], .fifth-value, .attr-fifth');

        if (halfSpan) halfSpan.textContent = this.half(value);
        if (fifthSpan) fifthSpan.textContent = this.fifth(value);
    },

    // =========================================================================
    // Visual Feedback
    // =========================================================================

    /**
     * Show success feedback on element
     * @param {HTMLElement} element - Target element
     * @param {number} duration - Duration in ms (default: 300)
     */
    showSuccess(element, duration = 300) {
        element.classList.add('bg-success', 'bg-opacity-10');
        setTimeout(() => {
            element.classList.remove('bg-success', 'bg-opacity-10');
        }, duration);
    },

    /**
     * Show error feedback on element
     * @param {HTMLElement} element - Target element
     * @param {number} duration - Duration in ms (default: 300)
     */
    showError(element, duration = 300) {
        element.classList.add('bg-danger', 'bg-opacity-10');
        setTimeout(() => {
            element.classList.remove('bg-danger', 'bg-opacity-10');
        }, duration);
    },

    /**
     * Show invalid feedback on element
     * @param {HTMLElement} element - Target element
     * @param {number} duration - Duration in ms (default: 800)
     */
    showInvalid(element, duration = 800) {
        element.classList.add('is-invalid');
        setTimeout(() => {
            element.classList.remove('is-invalid');
        }, duration);
    },

    /**
     * Flash highlight effect
     * @param {HTMLElement} element - Target element
     * @param {number} duration - Duration in ms (default: 800)
     */
    flashHighlight(element, duration = 800) {
        element.classList.add('flash-highlight');
        setTimeout(() => {
            element.classList.remove('flash-highlight');
        }, duration);
    },

    /**
     * Show toast notification
     * @param {string} title - Toast title
     * @param {string} message - Toast message
     * @param {string} icon - Icon (emoji)
     * @param {number} duration - Duration in ms (default: 1500)
     */
    showToast(title, message, icon = '', duration = 1500) {
        const toast = document.createElement('div');
        toast.className = 'position-fixed bottom-0 end-0 p-3';
        toast.style.zIndex = '1050';
        toast.innerHTML = `
            <div class="toast show" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">${icon} ${title}</strong>
                    <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close" onclick="this.closest('.position-fixed').remove()"></button>
                </div>
                <div class="toast-body">${message}</div>
            </div>
        `;
        document.body.appendChild(toast);

        setTimeout(() => toast.remove(), duration);
    },

    // =========================================================================
    // Form Utilities
    // =========================================================================

    /**
     * Disable/enable a button with loading state
     * @param {HTMLElement} button - Button element
     * @param {boolean} loading - Loading state
     * @param {string} loadingText - Text to show while loading
     */
    setButtonLoading(button, loading, loadingText = 'Processing...') {
        if (loading) {
            button.disabled = true;
            button.dataset.originalText = button.innerHTML;
            button.innerHTML = `<span class="spinner-border spinner-border-sm me-2"></span>${loadingText}`;
        } else {
            button.disabled = false;
            button.innerHTML = button.dataset.originalText || button.innerHTML;
        }
    },

    /**
     * Update button style based on enabled state
     * @param {HTMLElement} button - Button element
     * @param {boolean} enabled - Enabled state
     */
    updateButtonState(button, enabled) {
        button.disabled = !enabled;
        if (enabled) {
            button.style.background = 'linear-gradient(135deg, #6d6875 0%, #b5838d 100%)';
            button.classList.add('pulse-button');
        } else {
            button.style.background = '#e5e5e5';
            button.classList.remove('pulse-button');
        }
    },

    // =========================================================================
    // Data Utilities
    // =========================================================================

    /**
     * Parse integer with fallback
     * @param {*} value - Value to parse
     * @param {number} fallback - Fallback value (default: 0)
     * @returns {number}
     */
    parseInt(value, fallback = 0) {
        const parsed = parseInt(value, 10);
        return isNaN(parsed) ? fallback : parsed;
    },

    /**
     * Get current character data from hidden input
     * @returns {object|null}
     */
    getCurrentCharacter() {
        const hiddenInput = this.$('currentCharacter');
        if (!hiddenInput || !hiddenInput.value) {
            return null;
        }
        try {
            return JSON.parse(hiddenInput.value);
        } catch (e) {
            console.error('Failed to parse character data:', e);
            return null;
        }
    },

    /**
     * Get current character ID from various sources
     * @returns {string}
     */
    getCurrentCharacterId() {
        // Try name field (character sheet context)
        const nameInput = this.qs('input[data-field="Name"]');
        if (nameInput && nameInput.id) {
            return nameInput.id;
        }

        // Try hidden investigatorId field (wizard context)
        const hiddenInput = this.$('investigatorId');
        if (hiddenInput && hiddenInput.value) {
            return hiddenInput.value;
        }

        return '';
    },

    /**
     * Download blob as file
     * @param {Blob} blob - File blob
     * @param {string} filename - File name
     */
    downloadBlob(blob, filename) {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    },

    // =========================================================================
    // Attribute Name Mapping
    // =========================================================================

    /** Mapping from abbreviation to full name */
    ATTRIBUTE_NAMES: {
        'POW': 'Power',
        'STR': 'Strength',
        'LCK': 'Luck',
        'APP': 'Appearance',
        'DEX': 'Dexterity',
        'INT': 'Intelligence',
        'EDU': 'Education',
        'SIZ': 'Size',
        'CON': 'Constitution',
    },

    /**
     * Get full attribute name from abbreviation
     * @param {string} abbrev - Abbreviation
     * @returns {string}
     */
    getAttributeName(abbrev) {
        return this.ATTRIBUTE_NAMES[abbrev] || abbrev;
    },
};

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = Utils;
}

// Make available globally
window.Utils = Utils;