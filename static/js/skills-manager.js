/**
 * Skills Manager Module - Handles collapsible skills section and search
 * @module skills-manager
 */

const SkillsManager = {
    /** @type {boolean} */
    isCollapsed: true,

    /**
     * Toggle the collapsed state of the skills section
     */
    toggleCollapse() {
        const content = document.getElementById('all-skills-content');
        const chevron = document.querySelector('.all-skills-chevron');

        if (!content) return;

        this.isCollapsed = !this.isCollapsed;

        // Remember state for when page reloads (e.g., after pinning a skill)
        sessionStorage.setItem('skillsCollapsed', this.isCollapsed ? 'true' : 'false');

        if (this.isCollapsed) {
            content.classList.add('collapsed');
            if (chevron) chevron.classList.remove('rotated');
        } else {
            content.classList.remove('collapsed');
            if (chevron) chevron.classList.add('rotated');
        }
    },

    /**
     * Expand the skills section (used when clicking search)
     */
    expand() {
        if (this.isCollapsed) {
            this.toggleCollapse();
        }
    },

    /**
     * Collapse the skills section
     */
    collapse() {
        if (!this.isCollapsed) {
            this.toggleCollapse();
        }
    },

    /**
     * Filter skills based on search query
     * @param {string} query - Search query
     */
    filterSkills(query) {
        const skillsGrid = document.getElementById('skills-grid');
        const noResults = document.getElementById('skills-no-results');
        const clearBtn = document.querySelector('.skills-search-clear');

        if (!skillsGrid) return;

        const normalizedQuery = query.toLowerCase().trim();
        const skillItems = skillsGrid.querySelectorAll('.skill-item');
        let visibleCount = 0;

        skillItems.forEach(item => {
            const skillName = item.dataset.skillName || '';
            const matches = skillName.toLowerCase().includes(normalizedQuery);

            if (matches || normalizedQuery === '') {
                item.style.display = '';
                visibleCount++;
            } else {
                item.style.display = 'none';
            }
        });

        // Show/hide no results message
        if (noResults) {
            noResults.style.display = visibleCount === 0 && normalizedQuery !== '' ? 'flex' : 'none';
        }

        // Show/hide clear button
        if (clearBtn) {
            clearBtn.style.display = normalizedQuery !== '' ? 'flex' : 'none';
        }

        // Auto-expand when searching
        if (normalizedQuery !== '' && this.isCollapsed) {
            this.expand();
        }
    },

    /**
     * Clear the search input and reset filter
     */
    clearSearch() {
        const searchInput = document.getElementById('skills-search');
        if (searchInput) {
            searchInput.value = '';
            this.filterSkills('');
            searchInput.focus();
        }
    },

    /**
     * Handle click on skill item to roll dice
     * @param {Event} event - Click event
     * @param {HTMLElement} skillItem - The skill item element
     */
    handleSkillClick(event, skillItem) {
        // For custom skills with text input, get the current name from the input
        const nameInput = skillItem.querySelector('.skill-name-input');
        const skillName = nameInput ? nameInput.value : skillItem.dataset.skillName;

        // Get skill value from the value input (most current)
        const valueInput = skillItem.querySelector('.skill-value-field');
        const skillValue = valueInput
            ? parseInt(valueInput.value, 10)
            : parseInt(skillItem.dataset.skillValue, 10);

        if (skillName && !isNaN(skillValue) && window.HelperPanel) {
            HelperPanel.quickRollSkill(skillName, skillValue);
        }
    },

    /**
     * Initialize the skills manager
     */
    init() {
        // Restore collapsed state from sessionStorage (persists across HTMX reloads)
        const savedState = sessionStorage.getItem('skillsCollapsed');
        this.isCollapsed = savedState === null ? true : savedState === 'true';

        // Apply the restored state
        const content = document.getElementById('all-skills-content');
        const chevron = document.querySelector('.all-skills-chevron');

        if (content && chevron) {
            if (this.isCollapsed) {
                content.classList.add('collapsed');
                chevron.classList.remove('rotated');
            } else {
                content.classList.remove('collapsed');
                chevron.classList.add('rotated');
            }
        }

        // Hide clear button initially
        const clearBtn = document.querySelector('.skills-search-clear');
        if (clearBtn) {
            clearBtn.style.display = 'none';
        }
    }
};

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    SkillsManager.init();
});

// Re-initialize after HTMX swaps
document.addEventListener('htmx:afterSettle', () => {
    SkillsManager.init();
});

// Export globally
window.SkillsManager = SkillsManager;
