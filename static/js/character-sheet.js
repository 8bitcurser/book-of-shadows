/**
 * Character Sheet Module - Handles character sheet viewing and editing
 * @module character-sheet
 */

const CharacterSheet = {
    // =========================================================================
    // Initialization
    // =========================================================================

    /**
     * Initialize character sheet functionality
     */
    init() {
        this.initHoverEffects();
        this.initSkillNameAdjustment();
    },

    /**
     * Initialize hover effects for interactive elements
     * Note: Hover effects are now handled purely in CSS for better theme support
     */
    initHoverEffects() {
        // Hover effects are handled via CSS :hover selectors
        // This method is kept for backwards compatibility
    },

    /**
     * Adjust skill names based on container width
     */
    initSkillNameAdjustment() {
        const adjustSkillNames = () => {
            Utils.qsa('.skill-item').forEach(item => {
                const container = item.querySelector('.skill-name-container');
                const name = item.querySelector('.skill-name');

                if (container && name) {
                    const containerWidth = container.offsetWidth;
                    if (containerWidth < 120) {
                        name.style.maxWidth = (containerWidth - 30) + 'px';
                    }
                }
            });
        };

        adjustSkillNames();
        window.addEventListener('resize', adjustSkillNames);
    },

    // =========================================================================
    // Core Attribute Adjustment
    // =========================================================================

    /**
     * Adjust attribute value by delta (triggered by +/- buttons)
     * @param {HTMLButtonElement} button - The +/- button clicked
     * @param {number} delta - Amount to change (+1 or -1)
     */
    async adjustAttribute(button, delta) {
        const charBox = button.closest('.characteristic-box');
        if (!charBox) return;

        const badge = charBox.querySelector('.characteristic-value');
        if (!badge) return;

        const abbrev = badge.dataset.attr;
        const currentValue = Utils.parseInt(badge.textContent);
        const newValue = Math.max(1, Math.min(99, currentValue + delta));

        if (newValue === currentValue) return;

        // Update display immediately
        badge.textContent = newValue;
        Utils.updateDerivedValues(charBox, newValue);

        const investigatorId = Utils.getCurrentCharacterId();
        if (!investigatorId) return;

        const attrName = Utils.getAttributeName(abbrev);

        try {
            await API.updateInvestigator(
                investigatorId,
                'attributes',
                attrName,
                newValue
            );

            // Refresh character sheet to show updated derivatives (HP, MP, Move, etc.)
            this.refreshCombatStats(investigatorId);
        } catch (error) {
            console.error('Error updating core attribute:', error);
            // Revert on error
            badge.textContent = currentValue;
            Utils.updateDerivedValues(charBox, currentValue);
        }
    },

    /**
     * Refresh the character sheet to show updated derivative values
     * @param {string} investigatorId - The investigator ID
     */
    async refreshCombatStats(investigatorId) {
        try {
            // Save scroll position before reload
            const scrollY = window.scrollY;

            // Use same approach as togglePinSkill for consistency
            const html = await API.getInvestigator(investigatorId);
            Utils.setHTML('character-sheet', html);
            window.scrollTo(0, scrollY);

            // Re-initialize skills manager
            if (window.SkillsManager) {
                SkillsManager.init();
            }
        } catch (error) {
            console.error('Error refreshing character sheet:', error);
        }
    },

    // =========================================================================
    // Lock/Unlock Functionality
    // =========================================================================

    /**
     * Toggle lock state of character sheet
     * @param {HTMLInputElement} checkbox - Lock checkbox element
     */
    toggleLock(checkbox) {
        const isLocked = checkbox.checked;

        Utils.qsa('.editable').forEach(element => {
            element.disabled = isLocked;
        });

        const title = isLocked ? 'Locked' : 'Unlocked';
        const icon = isLocked ? '\uD83D\uDD12' : '\uD83D\uDD13';
        const message = isLocked
            ? 'Character sheet is now locked. Unlock to make changes.'
            : 'Character sheet is now editable.';

        Utils.showToast(title, message, icon);
    },

    // =========================================================================
    // Value Updates
    // =========================================================================

    /**
     * Recalculate and update sheet values
     * @param {HTMLInputElement} input - Input element
     * @param {string} type - Value type (attribute, skill, stat)
     */
    async recalculateValues(input, type) {
        const value = Utils.parseInt(input.value);

        // Update derived values
        const container = input.closest('.characteristic-box') || input.parentElement;
        Utils.updateDerivedValues(container, value);

        const investigatorId = Utils.getCurrentCharacterId();
        if (!investigatorId) return;

        try {
            if (type === 'attribute') {
                await this.updateAttribute(input, value);
            } else if (type === 'skill') {
                await this.updateSkill(input, value);
            } else {
                await this.updateStat(input, value);
            }
            Utils.showSuccess(input);
        } catch (error) {
            console.error(`Error updating ${type}:`, error);
            Utils.showError(input);
        }
    },

    /**
     * Update attribute value
     * @param {HTMLInputElement} input - Attribute input
     * @param {number} value - New value
     */
    async updateAttribute(input, value) {
        const attrName = input.dataset.attr;
        await API.updateInvestigator(
            Utils.getCurrentCharacterId(),
            'combat',
            attrName,
            value
        );
    },

    /**
     * Update skill value
     * @param {HTMLInputElement} input - Skill input
     * @param {number} value - New value
     */
    async updateSkill(input, value) {
        const skillName = input.dataset.skill;

        // Update derived values in skill item
        const skillItem = input.closest('.skill-item');
        if (skillItem) {
            const derivedValues = skillItem.querySelector('.derived-values');
            if (derivedValues) {
                const spans = derivedValues.querySelectorAll('span');
                if (spans.length >= 3) {
                    spans[0].textContent = Utils.half(value);
                    spans[2].textContent = Utils.fifth(value);
                }
            }
        }

        await API.updateInvestigator(
            Utils.getCurrentCharacterId(),
            'skills',
            skillName,
            value
        );

        input.dataset.skillvalue = value.toString();
    },

    /**
     * Update stat value
     * @param {HTMLInputElement} input - Stat input
     * @param {number} value - New value
     */
    async updateStat(input, value) {
        const statName = input.dataset.stat;
        await API.updateInvestigator(
            Utils.getCurrentCharacterId(),
            'stats',
            statName,
            value
        );
    },

    // =========================================================================
    // Personal Info
    // =========================================================================

    /**
     * Update personal info field
     * @param {HTMLInputElement} input - Personal info input
     */
    async updatePersonalInfo(input) {
        const field = input.dataset.field;
        const value = field === 'age' ? Utils.parseInt(input.value) : input.value;
        const investigatorId = Utils.getCurrentCharacterId();

        if (!investigatorId) return;

        try {
            await API.updateInvestigator(investigatorId, 'personalInfo', field, value);
            Utils.showSuccess(input);
        } catch (error) {
            console.error('Error updating personal info:', error);
            Utils.showError(input);
        }
    },

    /** Timer for typewriter effect */
    _typingTimer: null,

    /**
     * Update header name and avatar in real-time as user types
     * @param {HTMLInputElement} input - Name input field
     */
    updateHeaderName(input) {
        const newName = input.value.trim();

        // Update the header name
        const headerName = Utils.$('header-name');
        if (headerName) {
            headerName.textContent = newName || 'Unnamed';

            // Add typing class for typewriter cursor effect
            headerName.classList.add('typing');

            // Clear existing timer
            if (this._typingTimer) {
                clearTimeout(this._typingTimer);
            }

            // Remove typing class after user stops typing
            this._typingTimer = setTimeout(() => {
                headerName.classList.remove('typing');
            }, 800);
        }

        // Update the avatar with the first character
        const headerAvatar = Utils.$('header-avatar');
        if (headerAvatar) {
            headerAvatar.textContent = newName ? newName.charAt(0).toUpperCase() : '?';
        }
    },

    // =========================================================================
    // Skill Management
    // =========================================================================

    /**
     * Toggle skill check mark
     * @param {HTMLInputElement} input - Checkbox input
     */
    async handleSkillToggleCheck(input) {
        const skillName = input.dataset.skill;
        try {
            await API.updateInvestigator(
                Utils.getCurrentCharacterId(),
                'skill_check',
                skillName,
                true
            );
        } catch (error) {
            console.error('Error toggling skill check:', error);
        }
    },

    /**
     * Handle skill name change (for custom skills)
     * @param {HTMLInputElement} input - Name input
     */
    async handleSkillNameChange(input) {
        const skillName = input.dataset.skill;
        try {
            await API.updateInvestigator(
                Utils.getCurrentCharacterId(),
                'skill_name',
                skillName,
                input.value
            );
        } catch (error) {
            console.error('Error updating skill name:', error);
        }
    },

    /**
     * Toggle skill pin/priority status
     * @param {HTMLButtonElement} button - Pin button
     */
    async togglePinSkill(button) {
        const skillName = button.dataset.skill;
        const isPinned = button.dataset.pinned === 'true';
        const newPinnedStatus = !isPinned;

        // Update button appearance
        button.dataset.pinned = newPinnedStatus.toString();

        const icon = button.querySelector('i');
        if (icon) {
            if (newPinnedStatus) {
                icon.classList.remove('bi-pin');
                icon.classList.add('bi-pin-fill');
                icon.style.color = '#C97700'; // Occult Amber
            } else {
                icon.classList.remove('bi-pin-fill');
                icon.classList.add('bi-pin');
                icon.style.color = '#B0B0B0'; // Phantom Gray
            }
        }

        const skillItem = button.closest('.skill-item');
        if (skillItem) {
            skillItem.dataset.priority = newPinnedStatus.toString();
        }

        try {
            await API.updateInvestigator(
                Utils.getCurrentCharacterId(),
                'skill_prio',
                skillName,
                newPinnedStatus
            );

            // Add visual feedback
            button.classList.add('scale-up');
            setTimeout(() => button.classList.remove('scale-up'), 300);

            // Save scroll position and reload skills section
            const scrollY = window.scrollY;
            const html = await API.getInvestigator(Utils.getCurrentCharacterId());
            Utils.setHTML('character-sheet', html);
            window.scrollTo(0, scrollY);

            // Re-initialize skills manager to restore collapsed state
            if (window.SkillsManager) {
                SkillsManager.init();
            }
        } catch (error) {
            console.error('Error updating skill priority:', error);
            // Revert on error
            button.dataset.pinned = isPinned.toString();
            if (skillItem) {
                skillItem.dataset.priority = isPinned.toString();
            }
            if (icon) {
                if (isPinned) {
                    icon.classList.remove('bi-pin');
                    icon.classList.add('bi-pin-fill');
                    icon.style.color = '#C97700';
                } else {
                    icon.classList.remove('bi-pin-fill');
                    icon.classList.add('bi-pin');
                    icon.style.color = '#B0B0B0';
                }
            }
        }
    },

    // =========================================================================
    // Export/Import
    // =========================================================================

    /**
     * Export character as PDF
     * @param {Event} evt - Click event
     * @param {string} key - Character key/ID
     */
    async exportPDF(evt, key) {
        try {
            const blob = await API.exportPDF(key);
            Utils.downloadBlob(blob, key + '.pdf');
        } catch (error) {
            console.error('Error exporting PDF:', error);
            Utils.showToast('Error', 'Failed to export PDF. Please try again.', '\u274C');
        }
    },

    /**
     * Import investigators from code
     */
    async importInvestigators() {
        const importCode = Utils.getValue('importCode');

        if (!importCode) {
            Utils.showToast('Error', 'Please enter an import code.', '\u274C');
            return;
        }

        try {
            await API.importInvestigators(importCode);

            // Close modal and trigger HTMX refresh
            const modal = Utils.$('importModal');
            if (modal) {
                modal.classList.add('hidden');
            }
            htmx.trigger('body', 'import');

            Utils.showToast('Success', 'Investigators imported successfully!', '\u2705');
        } catch (error) {
            console.error('Error importing investigators:', error);
            Utils.showToast('Error', 'Failed to import investigators. Please try again.', '\u274C');
        }
    },

    /**
     * Get export code for all investigators
     */
    async getExportCode() {
        try {
            return await API.getExportCode();
        } catch (error) {
            console.error('Error getting export code:', error);
            Utils.showToast('Error', 'Failed to get export code.', '\u274C');
            return null;
        }
    },

    // =========================================================================
    // Phobias & Manias Management
    // =========================================================================

    /**
     * Preview a condition's description when selected from dropdown
     * @param {HTMLSelectElement} select - The select element
     * @param {string} type - 'phobia' or 'mania'
     */
    previewCondition(select, type) {
        const previewId = type === 'phobia' ? 'phobia-preview' : 'mania-preview';
        const preview = Utils.$(previewId);

        if (!preview) return;

        const selectedOption = select.options[select.selectedIndex];
        const name = selectedOption.value;
        const description = selectedOption.dataset.description;

        if (name && description) {
            preview.innerHTML = `
                <p class="condition-preview-name mb-1">${name}</p>
                <p class="condition-preview-description mb-0">${description}</p>
            `;
        } else {
            const typeLabel = type === 'phobia' ? 'phobia' : 'mania';
            preview.innerHTML = `<p class="condition-preview-empty mb-0">Select a ${typeLabel} to see its description</p>`;
        }
    },

    /**
     * Add a condition (phobia or mania) to the investigator
     * @param {string} type - 'phobia' or 'mania'
     */
    async addCondition(type) {
        const selectId = type === 'phobia' ? 'phobia-select' : 'mania-select';
        const select = Utils.$(selectId);

        if (!select || !select.value) {
            Utils.showToast('Error', `Please select a ${type} to add.`, '\u274C');
            return;
        }

        const conditionName = select.value;
        const section = type === 'phobia' ? 'phobias' : 'manias';
        const investigatorId = Utils.getCurrentCharacterId();

        if (!investigatorId) return;

        try {
            await API.updateInvestigator(investigatorId, section, conditionName, true);

            // Save scroll position and reload the character sheet
            const scrollY = window.scrollY;
            const html = await API.getInvestigator(investigatorId);
            Utils.setHTML('character-sheet', html);
            window.scrollTo(0, scrollY);
            if (window.SkillsManager) SkillsManager.init();

            Utils.showToast('Added', `${conditionName} has been added.`, '\uD83D\uDCA5');
        } catch (error) {
            console.error(`Error adding ${type}:`, error);
            Utils.showToast('Error', `Failed to add ${type}.`, '\u274C');
        }
    },

    /**
     * Remove a condition (phobia or mania) from the investigator
     * @param {HTMLButtonElement} button - The remove button clicked
     */
    async removeCondition(button) {
        const conditionName = button.dataset.condition;
        const type = button.dataset.type;
        const section = type === 'phobia' ? 'phobias' : 'manias';
        const investigatorId = Utils.getCurrentCharacterId();

        if (!investigatorId) return;

        try {
            await API.updateInvestigator(investigatorId, section, conditionName, false);

            // Save scroll position and reload the character sheet
            const scrollY = window.scrollY;
            const html = await API.getInvestigator(investigatorId);
            Utils.setHTML('character-sheet', html);
            window.scrollTo(0, scrollY);
            if (window.SkillsManager) SkillsManager.init();

            Utils.showToast('Removed', `${conditionName} has been removed.`, '\u2705');
        } catch (error) {
            console.error(`Error removing ${type}:`, error);
            Utils.showToast('Error', `Failed to remove ${type}.`, '\u274C');
        }
    },
};

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CharacterSheet;
}

// Make available globally
window.CharacterSheet = CharacterSheet;
