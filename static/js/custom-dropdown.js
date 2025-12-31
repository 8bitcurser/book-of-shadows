/**
 * Custom Dropdown Module - Retro 8-Bit Styled Dropdowns
 * Transforms native select elements into fully styled custom dropdowns
 * @module CustomDropdown
 */

const CustomDropdown = {
    // Track all dropdown instances
    instances: [],

    /**
     * Initialize custom dropdowns on select elements
     * @param {string} selector - CSS selector for select elements to transform
     */
    init(selector = 'select.form-control') {
        const selects = document.querySelectorAll(selector);
        selects.forEach(select => {
            // Skip if already transformed
            if (select.classList.contains('has-custom-dropdown')) return;
            this.create(select);
        });

        // Close dropdowns when clicking outside
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.custom-dropdown')) {
                this.closeAll();
            }
        });

        // Close on escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.closeAll();
            }
        });
    },

    /**
     * Create a custom dropdown from a select element
     * @param {HTMLSelectElement} select - The select element to transform
     */
    create(select) {
        // Mark the original select
        select.classList.add('has-custom-dropdown');

        // Create dropdown container
        const dropdown = document.createElement('div');
        dropdown.className = 'custom-dropdown';
        dropdown.dataset.selectId = select.id;

        // Create trigger button
        const trigger = document.createElement('button');
        trigger.type = 'button';
        trigger.className = 'custom-dropdown-trigger';

        const selectedText = document.createElement('span');
        selectedText.className = 'custom-dropdown-selected';

        const arrow = document.createElement('span');
        arrow.className = 'custom-dropdown-arrow';

        trigger.appendChild(selectedText);
        trigger.appendChild(arrow);

        // Create dropdown menu
        const menu = document.createElement('div');
        menu.className = 'custom-dropdown-menu';

        // Populate options
        this.populateOptions(select, menu, selectedText);

        // Assemble dropdown
        dropdown.appendChild(trigger);
        dropdown.appendChild(menu);

        // Insert after select
        select.parentNode.insertBefore(dropdown, select.nextSibling);

        // Event: Toggle dropdown
        trigger.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            this.toggle(dropdown);
        });

        // Store reference
        this.instances.push({
            select,
            dropdown,
            trigger,
            menu,
            selectedText
        });

        // Watch for select changes (e.g., from JavaScript)
        select.addEventListener('change', () => {
            this.syncFromSelect(select, dropdown);
        });
    },

    /**
     * Populate dropdown options from select
     * @param {HTMLSelectElement} select - Source select element
     * @param {HTMLElement} menu - Dropdown menu element
     * @param {HTMLElement} selectedText - Element to show selected value
     */
    populateOptions(select, menu, selectedText) {
        menu.innerHTML = '';
        let hasSelection = false;

        Array.from(select.options).forEach((option, index) => {
            const optionEl = document.createElement('div');
            optionEl.className = 'custom-dropdown-option';
            optionEl.dataset.value = option.value;
            optionEl.dataset.index = index;

            // Store description data attribute if exists
            if (option.dataset.description) {
                optionEl.dataset.description = option.dataset.description;
            }

            // Handle disabled options (separators)
            if (option.disabled) {
                optionEl.classList.add('disabled');
                optionEl.textContent = option.textContent;
                menu.appendChild(optionEl);
                return;
            }

            // Handle suggested occupations
            if (option.classList.contains('suggested-occupation') || option.textContent.includes('*')) {
                optionEl.classList.add('suggested');
            }

            // Clean star from text if present
            let displayText = option.textContent.replace(/^[\*\s]+/, '');
            optionEl.textContent = displayText;

            // Mark selected
            if (option.selected && option.value) {
                optionEl.classList.add('selected');
                selectedText.textContent = displayText;
                selectedText.classList.remove('custom-dropdown-placeholder');
                hasSelection = true;
            }

            // Click handler
            optionEl.addEventListener('click', (e) => {
                e.stopPropagation();
                this.selectOption(select, optionEl, menu, selectedText);
            });

            menu.appendChild(optionEl);
        });

        // Set placeholder if no selection
        if (!hasSelection) {
            const firstOption = select.options[0];
            if (firstOption && !firstOption.value) {
                selectedText.textContent = firstOption.textContent;
                selectedText.classList.add('custom-dropdown-placeholder');
            }
        }
    },

    /**
     * Select an option
     * @param {HTMLSelectElement} select - Original select element
     * @param {HTMLElement} optionEl - Clicked option element
     * @param {HTMLElement} menu - Dropdown menu
     * @param {HTMLElement} selectedText - Selected text display element
     */
    selectOption(select, optionEl, menu, selectedText) {
        const value = optionEl.dataset.value;
        const index = parseInt(optionEl.dataset.index, 10);

        // Update native select
        select.selectedIndex = index;
        select.value = value;

        // Update visual state
        menu.querySelectorAll('.custom-dropdown-option').forEach(opt => {
            opt.classList.remove('selected');
        });
        optionEl.classList.add('selected');

        // Update displayed text
        selectedText.textContent = optionEl.textContent;
        selectedText.classList.remove('custom-dropdown-placeholder');

        // Close dropdown
        this.closeAll();

        // Trigger change event on original select
        // This will invoke any onchange handlers (both attribute and addEventListener)
        const changeEvent = new Event('change', { bubbles: true });
        select.dispatchEvent(changeEvent);
    },

    /**
     * Sync dropdown state from select element
     * @param {HTMLSelectElement} select - Source select element
     * @param {HTMLElement} dropdown - Dropdown container
     */
    syncFromSelect(select, dropdown) {
        const menu = dropdown.querySelector('.custom-dropdown-menu');
        const selectedText = dropdown.querySelector('.custom-dropdown-selected');

        // Rebuild options in case they changed
        this.populateOptions(select, menu, selectedText);
    },

    /**
     * Toggle dropdown open/closed
     * @param {HTMLElement} dropdown - Dropdown container
     */
    toggle(dropdown) {
        const isOpen = dropdown.classList.contains('open');

        // Close all other dropdowns first
        this.closeAll();

        if (!isOpen) {
            dropdown.classList.add('open');

            // Position the menu using fixed positioning to escape overflow:hidden
            const trigger = dropdown.querySelector('.custom-dropdown-trigger');
            const menu = dropdown.querySelector('.custom-dropdown-menu');
            const rect = trigger.getBoundingClientRect();

            menu.style.position = 'fixed';
            menu.style.top = `${rect.bottom + 4}px`;
            menu.style.left = `${rect.left}px`;
            menu.style.width = `${rect.width}px`;

            // Scroll selected option into view
            const selected = menu.querySelector('.custom-dropdown-option.selected');
            if (selected) {
                selected.scrollIntoView({ block: 'nearest' });
            }
        }
    },

    /**
     * Close all open dropdowns
     */
    closeAll() {
        document.querySelectorAll('.custom-dropdown.open').forEach(dropdown => {
            dropdown.classList.remove('open');
        });
    },

    /**
     * Refresh a specific dropdown (call after options change)
     * @param {string} selectId - ID of the select element
     */
    refresh(selectId) {
        const instance = this.instances.find(i => i.select.id === selectId);
        if (instance) {
            this.syncFromSelect(instance.select, instance.dropdown);
        }
    },

    /**
     * Destroy a custom dropdown and restore native select
     * @param {string} selectId - ID of the select element
     */
    destroy(selectId) {
        const index = this.instances.findIndex(i => i.select.id === selectId);
        if (index > -1) {
            const instance = this.instances[index];
            instance.select.classList.remove('has-custom-dropdown');
            instance.dropdown.remove();
            this.instances.splice(index, 1);
        }
    },

    /**
     * Destroy all custom dropdowns
     */
    destroyAll() {
        this.instances.forEach(instance => {
            instance.select.classList.remove('has-custom-dropdown');
            instance.dropdown.remove();
        });
        this.instances = [];
    }
};

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CustomDropdown;
}

// Make available globally
window.CustomDropdown = CustomDropdown;

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    CustomDropdown.init();
});

// Re-initialize on HTMX content swap
document.addEventListener('htmx:afterSwap', () => {
    CustomDropdown.init();
});

// Re-initialize on HTMX content settle (for fully settled DOM)
document.addEventListener('htmx:afterSettle', () => {
    CustomDropdown.init();
});
