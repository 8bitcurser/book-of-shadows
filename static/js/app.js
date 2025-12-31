/**
 * Main Application Module - Entry point and event coordination
 * @module app
 */

const App = {
    // =========================================================================
    // Initialization
    // =========================================================================

    /**
     * Initialize the application
     */
    init() {
        this.initHTMXHandlers();
        this.initBugReportModal();
    },

    /**
     * Initialize HTMX event handlers
     */
    initHTMXHandlers() {
        // Handle content swap for character sheet display
        document.body.addEventListener('htmx:afterSwap', (evt) => {
            if (evt.detail.target.id === 'character-sheet') {
                const mainContent = Utils.$('main-content');
                const characterSheet = Utils.$('character-sheet');

                if (mainContent) mainContent.style.display = 'none';
                if (characterSheet) characterSheet.style.display = 'block';
            }
        });
    },

    // =========================================================================
    // Bug Report Modal
    // =========================================================================

    /**
     * Initialize bug report modal functionality
     */
    initBugReportModal() {
        const reportModal = Utils.$('reportIssueModal');
        if (reportModal) {
            reportModal.addEventListener('show.bs.modal', () => {
                this.generateCaptcha();
            });
        }

        const submitBtn = Utils.$('submitIssueBtn');
        if (submitBtn) {
            submitBtn.addEventListener('click', () => {
                this.submitIssueReport();
            });
        }
    },

    /**
     * Generate a random math CAPTCHA
     */
    generateCaptcha() {
        const operations = ['+', '-'];
        const operation = operations[Math.floor(Math.random() * operations.length)];

        let num1, num2, answer;

        if (operation === '+') {
            num1 = Math.floor(Math.random() * 10) + 1;
            num2 = Math.floor(Math.random() * 10) + 1;
            answer = num1 + num2;
        } else {
            num1 = Math.floor(Math.random() * 10) + 10;
            num2 = Math.floor(Math.random() * num1);
            answer = num1 - num2;
        }

        Utils.setText('captchaQuestion', `What is ${num1} ${operation} ${num2}?`);

        const expectedAnswer = Utils.$('expectedAnswer');
        if (expectedAnswer) expectedAnswer.value = answer;

        const captchaAnswer = Utils.$('captchaAnswer');
        if (captchaAnswer) captchaAnswer.value = '';

        // Reset form message
        const messageEl = Utils.$('formMessage');
        if (messageEl) {
            messageEl.classList.add('d-none');
            messageEl.classList.remove('alert-danger', 'alert-success');
        }
    },

    /**
     * Submit bug/issue report
     */
    async submitIssueReport() {
        const form = Utils.$('issueReportForm');
        const messageEl = Utils.$('formMessage');
        const submitBtn = Utils.$('submitIssueBtn');

        if (!form || !messageEl || !submitBtn) return;

        // Check form validity
        if (!form.checkValidity()) {
            this.showFormMessage(messageEl, 'Please fill in all required fields.', false);
            return;
        }

        // Verify CAPTCHA
        const userAnswer = Utils.getValue('captchaAnswer');
        const expectedAnswer = Utils.getValue('expectedAnswer');

        if (userAnswer !== expectedAnswer) {
            this.showFormMessage(messageEl, 'Incorrect CAPTCHA answer. Please try again.', false);
            this.generateCaptcha();
            return;
        }

        // Show loading state
        const originalBtnText = submitBtn.textContent;
        Utils.setButtonLoading(submitBtn, true, 'Sending...');

        // Collect form data
        const reportData = {
            issueType: Utils.getValue('issueType'),
            description: Utils.getValue('issueDescription'),
            email: Utils.getValue('contactEmail') || 'No email provided',
            timestamp: new Date().toISOString()
        };

        try {
            await API.submitBugReport(reportData);

            this.showFormMessage(messageEl, "Thank you for your report! We'll look into this issue.", true);

            // Reset form fields
            const issueType = Utils.$('issueType');
            const issueDescription = Utils.$('issueDescription');
            const contactEmail = Utils.$('contactEmail');

            if (issueType) issueType.value = '';
            if (issueDescription) issueDescription.value = '';
            if (contactEmail) contactEmail.value = '';

            // Auto-close modal after success
            setTimeout(() => {
                const modal = bootstrap.Modal.getInstance(Utils.$('reportIssueModal'));
                if (modal) modal.hide();
            }, 2000);
        } catch (error) {
            console.error('Error sending report:', error);
            this.showFormMessage(messageEl, 'Unable to send your report. Please try again later.', false);
        } finally {
            Utils.setButtonLoading(submitBtn, false);
            submitBtn.textContent = originalBtnText;
        }
    },

    /**
     * Show form message
     * @param {HTMLElement} messageEl - Message element
     * @param {string} message - Message text
     * @param {boolean} isSuccess - Success or error
     */
    showFormMessage(messageEl, message, isSuccess) {
        messageEl.textContent = message;
        messageEl.classList.remove('d-none', 'alert-danger', 'alert-success');
        messageEl.classList.add(isSuccess ? 'alert-success' : 'alert-danger');
    },
};

// =========================================================================
// Backward Compatibility Layer
// =========================================================================

/**
 * Legacy characterUtils object for backward compatibility with templates
 * Delegates to new modular structure
 */
const characterUtils = {
    // Wizard functions
    handleArchetypeSelection: (el) => Wizard.handleArchetypeSelection(el),
    handleOccupationSelection: (el) => Wizard.handleOccupationSelection(el),
    checkFormCompletion: () => Wizard.checkFormCompletion(),
    handlePersonalInfoChange: (input) => Wizard.handlePersonalInfoChange(input),
    initAttributeForm: () => Wizard.initAttributeForm(),
    rollAllAttributes: () => Wizard.rollAllAttributes(),
    rollAttribute: (input) => Wizard.rollSingleAttribute(input),
    updateAttributeValue: (input) => Wizard.updateAttributeValue(input),
    checkAttributesComplete: () => Wizard.checkAttributesComplete(),
    proceedToSkills: (id) => Wizard.proceedToSkills(id),
    initSkillForm: () => Wizard.initSkillForm(),
    recalculateSkillValues: (input) => Wizard.recalculateSkillValues(input),
    adjustSkillValue: (btn, inc) => Wizard.adjustSkillValue(btn, inc),
    navigateToTab: (name) => Wizard.navigateToTab(name),
    goBackToAttributes: (id) => Wizard.goBackToAttributes(id),
    completeCharacter: (id) => Wizard.completeCharacter(id),
    loadPersonalInfo: (id) => Wizard.loadPersonalInfo(id),

    // Character sheet functions
    initCharacterSheet: () => CharacterSheet.init(),
    recalculateSheetValues: (input, type) => CharacterSheet.recalculateValues(input, type),
    recalculateValues: (input, type) => CharacterSheet.recalculateValues(input, type),
    toggleLock: (checkbox) => CharacterSheet.toggleLock(checkbox),
    togglePinSkill: (btn) => CharacterSheet.togglePinSkill(btn),
    handleSkillToggleCheck: (input) => CharacterSheet.handleSkillToggleCheck(input),
    handleSkillNameChange: (input) => CharacterSheet.handleSkillNameChange(input),
    updatePersonalInfo: (input) => CharacterSheet.updatePersonalInfo(input),
    updateHeaderName: (input) => CharacterSheet.updateHeaderName(input),
    exportPDF: (evt, key) => CharacterSheet.exportPDF(evt, key),
    importInvestigators: () => CharacterSheet.importInvestigators(),

    // Utility functions
    getCurrentCharacter: () => Utils.getCurrentCharacter(),
    getCurrentCharacterId: () => Utils.getCurrentCharacterId(),
    updateDerivedValues: (input) => {
        const container = input.closest('.attribute-container');
        if (container) {
            Utils.updateDerivedValues(container, Utils.parseInt(input.value));
        }
    },

    // API function
    updateInvestigator: (section, field, value) => {
        const id = Utils.getCurrentCharacterId();
        return API.updateInvestigator(id, section, field, value);
    },

    // Legacy initialization (no-op, handled by App.init)
    initializeEventListeners: () => {},
    initializeInputHandlers: () => {},
};

// Legacy rollDice function
function rollDice(numDice, sides) {
    return Utils.rollDice(numDice, sides);
}

// =========================================================================
// Application Bootstrap
// =========================================================================

// Initialize application when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => App.init());
} else {
    App.init();
}

// Export modules globally
window.App = App;
window.characterUtils = characterUtils;
