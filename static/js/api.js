/**
 * API Module - Handles all server communication
 * @module api
 */

const API = {
    /**
     * Base fetch wrapper with error handling
     * @param {string} url - API endpoint
     * @param {object} options - Fetch options
     * @returns {Promise<Response>}
     */
    async request(url, options = {}) {
        const defaultOptions = {
            headers: {
                'Content-Type': 'application/json',
            },
        };

        const mergedOptions = {
            ...defaultOptions,
            ...options,
            headers: {
                ...defaultOptions.headers,
                ...options.headers,
            },
        };

        try {
            const response = await fetch(url, mergedOptions);

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            return response;
        } catch (error) {
            console.error(`API request failed: ${url}`, error);
            throw error;
        }
    },

    /**
     * GET request returning JSON
     * @param {string} url - API endpoint
     * @returns {Promise<object>}
     */
    async getJSON(url) {
        const response = await this.request(url, { method: 'GET' });
        return response.json();
    },

    /**
     * GET request returning HTML
     * @param {string} url - API endpoint
     * @returns {Promise<string>}
     */
    async getHTML(url) {
        const response = await this.request(url, { method: 'GET' });
        return response.text();
    },

    /**
     * POST request with JSON body
     * @param {string} url - API endpoint
     * @param {object} data - Request body
     * @returns {Promise<object>}
     */
    async postJSON(url, data) {
        const response = await this.request(url, {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return response.json();
    },

    /**
     * PUT request with JSON body
     * @param {string} url - API endpoint
     * @param {object} data - Request body
     * @returns {Promise<Response>}
     */
    async put(url, data) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    },

    // =========================================================================
    // Investigator API
    // =========================================================================

    /**
     * Create a new investigator
     * @param {object} data - Investigator data
     * @returns {Promise<{Key: string}>}
     */
    async createInvestigator(data) {
        return this.postJSON('/api/investigator/', data);
    },

    /**
     * Get investigator by ID (returns HTML)
     * @param {string} id - Investigator ID
     * @returns {Promise<string>}
     */
    async getInvestigator(id) {
        return this.getHTML(`/api/investigator/${id}`);
    },

    /**
     * Update investigator field
     * @param {string} id - Investigator ID
     * @param {string} section - Section to update (attributes, skills, stats, etc.)
     * @param {string} field - Field name
     * @param {*} value - New value
     * @returns {Promise<Response>}
     */
    async updateInvestigator(id, section, field, value) {
        return this.put(`/api/investigator/${id}`, {
            section,
            field,
            value,
        });
    },

    /**
     * Export investigator as PDF
     * @param {string} id - Investigator ID
     * @returns {Promise<Blob>}
     */
    async exportPDF(id) {
        const response = await this.request(`/api/investigator/PDF/${id}`, {
            method: 'POST',
            body: JSON.stringify({}),
        });
        return response.blob();
    },

    /**
     * Get export code for all investigators
     * @returns {Promise<string>}
     */
    async getExportCode() {
        const response = await this.request('/api/investigator/list/export', {
            method: 'GET',
        });
        return response.text();
    },

    /**
     * Import investigators from code
     * @param {string} importCode - Import code
     * @returns {Promise<Response>}
     */
    async importInvestigators(importCode) {
        return this.request('/api/investigator/list/import/', {
            method: 'POST',
            body: JSON.stringify({ ImportCode: importCode }),
        });
    },

    // =========================================================================
    // Archetype API
    // =========================================================================

    /**
     * Get occupations for an archetype
     * @param {string} archetypeName - Archetype name
     * @returns {Promise<{suggested: Array, others: Array}>}
     */
    async getArchetypeOccupations(archetypeName) {
        return this.getJSON(`/api/archetype/${encodeURIComponent(archetypeName)}/occupations`);
    },

    // =========================================================================
    // Wizard API
    // =========================================================================

    /**
     * Get wizard step HTML
     * @param {string} step - Step name (base, attributes, skills)
     * @param {string} id - Investigator ID
     * @returns {Promise<string>}
     */
    async getWizardStep(step, id) {
        return this.getHTML(`/wizard/${step}/${id}`);
    },

    // =========================================================================
    // Bug Report API
    // =========================================================================

    /**
     * Submit bug report
     * @param {object} report - Report data
     * @returns {Promise<object>}
     */
    async submitBugReport(report) {
        return this.postJSON('/api/report-issue', report);
    },
};

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = API;
}

// Make available globally
window.API = API;