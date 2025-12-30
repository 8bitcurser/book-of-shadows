/**
 * Wizard Module - Handles character creation wizard steps
 * @module wizard
 */

const Wizard = {
    // =========================================================================
    // Description Formatting
    // =========================================================================

    /**
     * Format archetype/occupation description into structured HTML
     * @param {string} text - Raw description text
     * @returns {string} - Formatted HTML
     */
    formatDescription(text) {
        if (!text) return '';

        // Parse the description into sections
        const sections = {
            intro: '',
            coreChar: '',
            bonusPoints: '',
            skills: '',
            talents: '',
            traits: '',
            occupations: ''
        };

        // Extract intro (text before first "Core Characteristics:" or other keywords)
        const keywordPattern = /(Core Characteristics:|Bonus Points:|Archetype Skills:|Occupation Skills:|Number of Talents:|Suggested Traits:|Suggested Occupations:)/;
        const firstKeyword = text.search(keywordPattern);

        if (firstKeyword > 0) {
            sections.intro = text.substring(0, firstKeyword).trim();
        }

        // Extract each section
        const extractSection = (label) => {
            const regex = new RegExp(label + '\\s*([^]*?)(?=Core Characteristics:|Bonus Points:|Archetype Skills:|Occupation Skills:|Number of Talents:|Suggested Traits:|Suggested Occupations:|$)', 'i');
            const match = text.match(regex);
            return match ? match[1].trim() : '';
        };

        sections.coreChar = extractSection('Core Characteristics:');
        sections.bonusPoints = extractSection('Bonus Points:');
        sections.skills = extractSection('Archetype Skills:') || extractSection('Occupation Skills:');
        sections.talents = extractSection('Number of Talents:');
        sections.traits = extractSection('Suggested Traits:');
        sections.occupations = extractSection('Suggested Occupations:');

        // Build HTML
        let html = '';

        // Intro
        if (sections.intro) {
            html += `<div class="desc-intro">${sections.intro}</div>`;
        }

        // Stats grid
        const stats = [];
        if (sections.coreChar) {
            stats.push({ label: 'Core Stats', value: sections.coreChar });
        }
        if (sections.bonusPoints) {
            stats.push({ label: 'Bonus Points', value: sections.bonusPoints });
        }
        if (sections.talents) {
            stats.push({ label: 'Talents', value: sections.talents });
        }

        if (stats.length > 0) {
            html += '<div class="desc-stats">';
            stats.forEach(stat => {
                html += `<div class="desc-stat">
                    <span class="desc-stat-label">${stat.label}</span>
                    <span class="desc-stat-value">${stat.value}</span>
                </div>`;
            });
            html += '</div>';
        }

        // Skills
        if (sections.skills) {
            const skillList = sections.skills.split(',').map(s => s.trim()).filter(s => s);
            html += `<div class="desc-skills">
                <span class="desc-skills-label">Skills</span>
                <div class="desc-skills-list">
                    ${skillList.map(skill => `<span class="desc-skill-tag">${skill}</span>`).join('')}
                </div>
            </div>`;
        }

        // Suggestions
        const suggestions = [];
        if (sections.traits) {
            suggestions.push({ label: 'Suggested Traits', value: sections.traits });
        }
        if (sections.occupations) {
            suggestions.push({ label: 'Suggested Occupations', value: sections.occupations });
        }

        if (suggestions.length > 0) {
            html += '<div class="desc-suggestions">';
            suggestions.forEach(sug => {
                html += `<div class="desc-suggestion">
                    <span class="desc-suggestion-label">${sug.label}</span>
                    <span class="desc-suggestion-value">${sug.value}</span>
                </div>`;
            });
            html += '</div>';
        }

        // If no structured content found, return plain text
        if (!html) {
            return `<div class="desc-intro">${text}</div>`;
        }

        return html;
    },

    // =========================================================================
    // Step 1: Personal Info / Base Step
    // =========================================================================

    /**
     * Initialize the base step (personal info form)
     */
    initBaseStep() {
        // Add change listeners to form fields
        const fields = ['name', 'age', 'residence', 'birthplace'];
        fields.forEach(field => {
            const input = Utils.qs(`input[name="${field}"]`);
            if (input) {
                input.addEventListener('input', () => this.checkFormCompletion());
                input.addEventListener('change', () => this.handlePersonalInfoChange(input));
            }
        });

        // Initial form validation check
        this.checkFormCompletion();
    },

    /**
     * Handle archetype selection change
     * @param {HTMLSelectElement} selectElement - Archetype select element
     */
    async handleArchetypeSelection(selectElement) {
        const archetypeName = selectElement.value;
        const descriptionElement = Utils.$('archetype-description');
        const occupationContainer = Utils.$('occupation-container');
        const selectedOption = selectElement.options[selectElement.selectedIndex];

        // Update description with formatted HTML
        if (selectedOption && selectedOption.dataset.description) {
            descriptionElement.innerHTML = this.formatDescription(selectedOption.dataset.description);
            descriptionElement.style.display = 'block';
        } else {
            descriptionElement.style.display = 'none';
        }

        // Show/hide occupation container
        occupationContainer.style.display = archetypeName ? 'block' : 'none';

        // Update occupation options
        if (archetypeName) {
            await this.updateOccupationOptions(archetypeName);
        }

        this.checkFormCompletion();
    },

    /**
     * Update occupation dropdown based on selected archetype
     * @param {string} archetypeName - Selected archetype name
     */
    async updateOccupationOptions(archetypeName) {
        const occupationSelect = Utils.$('occupation-select');
        const descriptionElement = Utils.$('occupation-description');

        if (!occupationSelect) return;

        // Reset occupation selection
        occupationSelect.value = '';
        if (descriptionElement) {
            descriptionElement.style.display = 'none';
        }

        try {
            const data = await API.getArchetypeOccupations(archetypeName);

            // Clear existing options (keep first "Select Occupation" option)
            while (occupationSelect.options.length > 1) {
                occupationSelect.removeChild(occupationSelect.lastChild);
            }

            // Add suggested occupations first
            if (data.suggested && data.suggested.length > 0) {
                data.suggested.forEach(occupation => {
                    const option = document.createElement('option');
                    option.value = occupation.name;
                    option.textContent = `⭐ ${occupation.name}`;
                    option.dataset.description = occupation.description;
                    option.className = 'suggested-occupation';
                    occupationSelect.appendChild(option);
                });

                // Add separator
                const separator = document.createElement('option');
                separator.value = '';
                separator.textContent = '────── Other Occupations ──────';
                separator.disabled = true;
                occupationSelect.appendChild(separator);
            }

            // Add other occupations
            if (data.others && data.others.length > 0) {
                data.others.forEach(occupation => {
                    const option = document.createElement('option');
                    option.value = occupation.name;
                    option.textContent = occupation.name;
                    option.dataset.description = occupation.description;
                    occupationSelect.appendChild(option);
                });
            }

            // Refresh custom dropdown to reflect new options
            if (window.CustomDropdown) {
                CustomDropdown.refresh('occupation-select');
            }
        } catch (error) {
            console.error('Error updating occupation options:', error);
        }
    },

    /**
     * Handle occupation selection change
     * @param {HTMLSelectElement} selectElement - Occupation select element
     */
    handleOccupationSelection(selectElement) {
        const descriptionElement = Utils.$('occupation-description');
        const selectedOption = selectElement.options[selectElement.selectedIndex];

        if (selectedOption && selectedOption.dataset.description) {
            descriptionElement.innerHTML = this.formatDescription(selectedOption.dataset.description);
            descriptionElement.style.display = 'block';
        } else {
            descriptionElement.style.display = 'none';
        }

        this.checkFormCompletion();
    },

    /**
     * Handle personal info field changes
     * @param {HTMLInputElement} input - Input element
     */
    async handlePersonalInfoChange(input) {
        const investigatorId = Utils.getCurrentCharacterId();

        if (investigatorId) {
            const field = input.dataset.field || input.name;
            const value = input.value;

            try {
                await API.updateInvestigator(investigatorId, 'personalInfo', field, value);
                Utils.showSuccess(input);
            } catch (error) {
                console.error('Error updating personal info:', error);
                Utils.showError(input);
            }
        }

        this.checkFormCompletion();
    },

    /**
     * Check if the base form is complete
     * @returns {boolean}
     */
    checkFormCompletion() {
        const nameInput = Utils.qs('input[name="name"]');
        const ageInput = Utils.qs('input[name="age"]');
        const residenceInput = Utils.qs('input[name="residence"]');
        const birthplaceInput = Utils.qs('input[name="birthplace"]');
        const archetypeSelect = Utils.$('archetype-select');
        const occupationSelect = Utils.$('occupation-select');
        const nextButton = Utils.$('next-step-button');

        if (!nextButton) return false;

        const isComplete = (
            nameInput?.value.trim() !== '' &&
            ageInput?.value !== '' &&
            residenceInput?.value.trim() !== '' &&
            birthplaceInput?.value.trim() !== '' &&
            archetypeSelect?.value !== '' &&
            occupationSelect?.value !== ''
        );

        Utils.updateButtonState(nextButton, isComplete);
        return isComplete;
    },

    /**
     * Handle form submission for step 1
     * @param {Event} event - Submit event
     */
    async handleFormSubmission(event) {
        event.preventDefault();

        const form = event.target.closest('form');
        const button = Utils.$('next-step-button');
        const existingId = Utils.$('investigatorId');

        Utils.setButtonLoading(button, true, 'Processing...');

        try {
            if (existingId?.value) {
                // Investigator exists, go to attributes
                const html = await API.getWizardStep('attributes', existingId.value);
                Utils.setHTML('character-sheet', html);
            } else {
                // Create new investigator
                const formData = new FormData(form);
                const jsonData = Object.fromEntries(formData.entries());

                const data = await API.createInvestigator(jsonData);
                if (data.Key) {
                    const html = await API.getWizardStep('attributes', data.Key);
                    Utils.setHTML('character-sheet', html);
                }
            }
        } catch (error) {
            console.error('Form submission error:', error);
            Utils.setButtonLoading(button, false);
            Utils.showToast('Error', 'Failed to proceed. Please try again.', '❌');
        }
    },

    // =========================================================================
    // Step 2: Attributes
    // =========================================================================

    /**
     * Initialize the attributes step
     */
    initAttributeForm() {
        // Clean URL query parameters
        if (window.location.search) {
            const cleanUrl = window.location.protocol + '//' + window.location.host + window.location.pathname;
            window.history.replaceState({}, document.title, cleanUrl);
        }

        // Set up input change listeners
        Utils.qsa('.attribute-input').forEach(input => {
            input.addEventListener('input', () => this.updateDerivedValues(input));
            input.addEventListener('change', () => this.updateAttributeValue(input));

            // Initial calculation for pre-filled values
            this.updateDerivedValues(input);
        });

        this.checkAttributesComplete();
    },

    /**
     * Update derived values for an attribute input
     * @param {HTMLInputElement} input - Attribute input element
     */
    updateDerivedValues(input) {
        const container = input.closest('.attribute-container');
        const value = Utils.parseInt(input.value);
        Utils.updateDerivedValues(container, value);
    },

    /**
     * Roll all attributes at once
     */
    rollAllAttributes() {
        const rollButton = Utils.qs('button[onclick*="rollAllAttributes"]');
        if (rollButton) {
            rollButton.classList.add('dice-rolling');
            setTimeout(() => rollButton.classList.remove('dice-rolling'), 500);
        }

        const attributeInputs = Utils.qsa('.attribute-input');
        const shuffledInputs = Array.from(attributeInputs).sort(() => Math.random() - 0.5);

        shuffledInputs.forEach((input, index) => {
            setTimeout(() => {
                this.rollSingleAttribute(input);

                const container = input.closest('.attribute-container');
                if (container) {
                    container.classList.add('highlight');
                    setTimeout(() => container.classList.remove('highlight'), 500);
                }
            }, index * 150);
        });
    },

    /**
     * Roll a single attribute with animation
     * @param {HTMLInputElement} input - Attribute input element
     */
    rollSingleAttribute(input) {
        const formula = input.dataset.formula;
        this.animateDiceRoll(input, formula);
    },

    /**
     * Animate dice roll
     * @param {HTMLInputElement} input - Attribute input element
     * @param {string} formula - Roll formula
     */
    animateDiceRoll(input, formula) {
        let iterations = 3;

        const animate = setInterval(() => {
            input.value = Utils.rollAttribute(formula);
            iterations--;

            if (iterations <= 0) {
                clearInterval(animate);
                input.value = Utils.rollAttribute(formula);
                this.updateDerivedValues(input);
                this.updateAttributeValue(input);
            }
        }, 100);
    },

    /**
     * Update attribute value on server
     * @param {HTMLInputElement} input - Attribute input element
     */
    async updateAttributeValue(input) {
        const attrAbbrev = input.name;
        const value = Utils.parseInt(input.value);
        const attrName = Utils.getAttributeName(attrAbbrev);

        this.updateDerivedValues(input);

        try {
            await API.updateInvestigator(
                Utils.getCurrentCharacterId(),
                'attributes',
                attrName,
                value
            );
            Utils.showSuccess(input);
            this.checkAttributesComplete();
        } catch (error) {
            console.error('Error updating attribute:', error);
            Utils.showError(input);
        }
    },

    /**
     * Check if all attributes are complete
     */
    checkAttributesComplete() {
        const attributeInputs = Utils.qsa('.attribute-input');
        const proceedButton = Utils.qs('.gradient-button');

        if (!proceedButton) return;

        const allFilled = Array.from(attributeInputs).every(input => {
            return Utils.parseInt(input.value) > 0;
        });

        Utils.updateButtonState(proceedButton, allFilled);
    },

    /**
     * Navigate to skills step
     * @param {string} investigatorId - Investigator ID
     */
    async proceedToSkills(investigatorId) {
        try {
            const html = await API.getWizardStep('skills', investigatorId);
            Utils.setHTML('character-sheet', html);
        } catch (error) {
            console.error('Error loading skills step:', error);
            Utils.showToast('Error', 'Failed to load skills step.', '❌');
        }
    },

    // =========================================================================
    // Step 3: Skills
    // =========================================================================

    /**
     * Initialize the skills form
     */
    initSkillForm() {
        // Ensure tab continue buttons are visible
        Utils.qsa('.transition-opacity').forEach(container => {
            container.style.opacity = '1';
            container.style.pointerEvents = 'auto';
        });

        // Hover effects are now handled via CSS :hover selectors
        // No JavaScript color overrides needed
    },

    /**
     * Recalculate skill values when changed
     * @param {HTMLInputElement} input - Skill input element
     */
    async recalculateSkillValues(input) {
        const skillName = input.dataset.skill;
        const value = Utils.parseInt(input.value);
        const prevValue = Utils.parseInt(input.dataset.skillvalue);
        const defaultValue = Utils.parseInt(input.dataset.skilldefault);
        const skillType = input.dataset.skilltype || 'archetype';

        // Apply limits
        if (value > 90) {
            input.value = 90;
            return;
        }
        if (value < defaultValue) {
            input.value = defaultValue;
            return;
        }

        const difference = value - prevValue;
        if (difference === 0) return;

        // Get proper points element based on skill type
        const pointsMap = {
            'archetype': { points: 'archetype-points', confirm: 'confirm-archetype-container' },
            'occupation': { points: 'occupation-points', confirm: 'confirm-occupation-container' },
            'general': { points: 'general-points', confirm: 'confirm-general-container' },
        };

        const activeTab = Utils.qs('.tab-pane.active');
        let type = skillType;
        if (activeTab) {
            if (activeTab.id === 'archetype-skills') type = 'archetype';
            else if (activeTab.id === 'occupation-skills') type = 'occupation';
            else if (activeTab.id === 'general-skills') type = 'general';
        }

        const { points: pointsId, confirm: confirmId } = pointsMap[type] || pointsMap.archetype;
        const pointsElement = Utils.$(pointsId);

        if (pointsElement) {
            const currentPoints = Utils.parseInt(pointsElement.textContent);
            const newPoints = currentPoints - difference;

            // Don't allow negative points
            if (newPoints < 0) {
                input.value = prevValue;
                Utils.showInvalid(input);
                return;
            }

            // Update points display
            pointsElement.textContent = newPoints;
            pointsElement.style.color = newPoints < 10 ? '#e84a5f' : '#63c74d';

            // Flash highlight the skill box
            const skillBox = input.closest('.skill-box');
            if (skillBox) Utils.flashHighlight(skillBox);

            // Update tracking
            input.dataset.skillvalue = value;

            // Update derived values
            const container = input.closest('.skill-values');
            if (container) Utils.updateDerivedValues(container, value);

            // Show continue button if all points used
            const confirmContainer = Utils.$(confirmId);
            if (confirmContainer && newPoints === 0) {
                confirmContainer.style.opacity = '1';
                confirmContainer.style.pointerEvents = 'auto';
            }
        }

        // Update server
        try {
            await API.updateInvestigator(
                Utils.getCurrentCharacterId(),
                'skills',
                skillName,
                value
            );
        } catch (error) {
            console.error('Error updating skill:', error);
        }
    },

    /**
     * Adjust skill value using +/- buttons
     * @param {HTMLButtonElement} btn - Button element
     * @param {boolean} increment - True to increment, false to decrement
     */
    adjustSkillValue(btn, increment) {
        const inputGroup = btn.closest('.skill-input-group');
        const input = inputGroup.querySelector('.skill-input');
        const currentValue = Utils.parseInt(input.value);

        input.value = currentValue + (increment ? 1 : -1);
        input.dispatchEvent(new Event('change'));
    },

    /**
     * Navigate between tabs
     * @param {string} tabName - Tab name (archetype, occupation, general)
     */
    navigateToTab(tabName) {
        if (tabName !== 'archetype') {
            const tab = Utils.$(tabName + '-tab');
            if (tab) tab.disabled = false;
        }

        setTimeout(() => {
            const tab = new bootstrap.Tab(Utils.$(tabName + '-tab'));
            tab.show();

            const content = Utils.$(tabName + '-skills');
            if (content) {
                content.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
        }, 100);
    },

    /**
     * Go back to attributes step
     * @param {string} investigatorId - Investigator ID
     */
    async goBackToAttributes(investigatorId) {
        try {
            const html = await API.getWizardStep('attributes', investigatorId);
            Utils.setHTML('character-sheet', html);
        } catch (error) {
            console.error('Error loading attributes step:', error);
        }
    },

    /**
     * Complete character creation and show character sheet
     * @param {string} investigatorId - Investigator ID
     */
    async completeCharacter(investigatorId) {
        try {
            const html = await API.getInvestigator(investigatorId);
            Utils.setHTML('character-sheet', html);
        } catch (error) {
            console.error('Error completing character:', error);
            Utils.showToast('Error', 'Failed to load character sheet.', '\u274C');
        }
    },

    /**
     * Load personal info step (go back from attributes)
     * @param {string} investigatorId - Investigator ID
     */
    async loadPersonalInfo(investigatorId) {
        try {
            const html = await API.getWizardStep('base', investigatorId);
            Utils.setHTML('character-sheet', html);
        } catch (error) {
            console.error('Error loading personal info:', error);
            Utils.showToast('Error', 'Failed to load personal info step.', '\u274C');
        }
    },
};

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = Wizard;
}

// Make available globally
window.Wizard = Wizard;