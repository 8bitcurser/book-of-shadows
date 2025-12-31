/**
 * Helper Panel Module - Dice roller, recovery calculator, and quick rules
 * @module helper-panel
 */

const HelperPanel = {
    /** @type {boolean} */
    isOpen: false,

    /** @type {Array} */
    rollHistory: [],

    // =========================================================================
    // Panel Controls
    // =========================================================================

    /**
     * Toggle the helper panel
     */
    toggle() {
        this.isOpen ? this.close() : this.open();
    },

    /**
     * Open the helper panel
     */
    open() {
        const panel = document.getElementById('helper-panel');
        const overlay = document.getElementById('helper-panel-overlay');
        const toggleBtn = document.querySelector('.helper-panel-toggle');

        if (panel) {
            panel.classList.add('open');
            this.isOpen = true;
        }
        if (overlay) {
            overlay.classList.add('visible');
        }
        if (toggleBtn) {
            toggleBtn.classList.add('active');
        }
    },

    /**
     * Close the helper panel
     */
    close() {
        const panel = document.getElementById('helper-panel');
        const overlay = document.getElementById('helper-panel-overlay');
        const toggleBtn = document.querySelector('.helper-panel-toggle');

        if (panel) {
            panel.classList.remove('open');
            this.isOpen = false;
        }
        if (overlay) {
            overlay.classList.remove('visible');
        }
        if (toggleBtn) {
            toggleBtn.classList.remove('active');
        }
    },

    /**
     * Switch between tabs
     * @param {string} tabId - The tab to switch to
     */
    switchTab(tabId) {
        // Update tab buttons
        document.querySelectorAll('.helper-tab').forEach(tab => {
            tab.classList.toggle('active', tab.dataset.tab === tabId);
        });

        // Update tab content
        document.querySelectorAll('.helper-tab-content').forEach(content => {
            content.classList.toggle('active', content.id === `tab-${tabId}`);
        });
    },

    // =========================================================================
    // Dice Rolling
    // =========================================================================

    /**
     * Quick roll a skill check (called when clicking on a skill)
     * @param {string} skillName - Name of the skill
     * @param {number} skillValue - Current skill value
     */
    quickRollSkill(skillName, skillValue) {
        const roll = Math.floor(Math.random() * 100) + 1;
        const half = Math.floor(skillValue / 2);
        const fifth = Math.floor(skillValue / 5);

        let outcome = '';
        let outcomeClass = '';
        let icon = '';

        if (roll === 1) {
            outcome = 'CRITICAL!';
            outcomeClass = 'critical-success';
            icon = '\u2728'; // sparkles
        } else if (roll === 100 || (skillValue < 50 && roll >= 96)) {
            outcome = 'FUMBLE!';
            outcomeClass = 'fumble';
            icon = '\uD83D\uDCA5'; // explosion
        } else if (roll <= fifth) {
            outcome = 'Extreme Success';
            outcomeClass = 'extreme-success';
            icon = '\uD83C\uDF1F'; // star
        } else if (roll <= half) {
            outcome = 'Hard Success';
            outcomeClass = 'hard-success';
            icon = '\u2705'; // check
        } else if (roll <= skillValue) {
            outcome = 'Success';
            outcomeClass = 'success';
            icon = '\u2714\uFE0F'; // check mark
        } else {
            outcome = 'Failure';
            outcomeClass = 'failure';
            icon = '\u274C'; // X
        }

        // Show toast with result
        this.showQuickRollToast(skillName, roll, skillValue, outcome, outcomeClass, icon);

        // Add to history
        this.addToHistory(`${skillName} (${skillValue}%)`, roll, outcome);

        // Update helper panel display if open
        this.displayResult(roll, `d100 vs ${skillValue}`, outcome, outcomeClass);
    },

    /**
     * Show a toast notification for quick roll result
     */
    showQuickRollToast(skillName, roll, target, outcome, outcomeClass, icon) {
        // Create toast element
        const toastId = 'quick-roll-toast';
        let toast = document.getElementById(toastId);

        if (!toast) {
            toast = document.createElement('div');
            toast.id = toastId;
            toast.className = 'quick-roll-toast';
            document.body.appendChild(toast);
        }

        // Determine result color class
        let colorClass = '';
        switch (outcomeClass) {
            case 'critical-success':
                colorClass = 'toast-critical';
                break;
            case 'extreme-success':
                colorClass = 'toast-extreme';
                break;
            case 'hard-success':
                colorClass = 'toast-hard';
                break;
            case 'success':
                colorClass = 'toast-success';
                break;
            case 'failure':
                colorClass = 'toast-failure';
                break;
            case 'fumble':
                colorClass = 'toast-fumble';
                break;
        }

        toast.className = `quick-roll-toast ${colorClass}`;
        toast.innerHTML = `
            <div class="toast-skill-name">${skillName}</div>
            <div class="toast-roll-result">
                <span class="toast-roll">${roll}</span>
                <span class="toast-vs">vs</span>
                <span class="toast-target">${target}%</span>
            </div>
            <div class="toast-outcome">${icon} ${outcome}</div>
        `;

        // Trigger animation
        toast.classList.remove('show');
        void toast.offsetWidth; // Force reflow
        toast.classList.add('show');

        // Auto-hide after 3 seconds
        clearTimeout(this._toastTimeout);
        this._toastTimeout = setTimeout(() => {
            toast.classList.remove('show');
        }, 3000);
    },

    /** Toast timeout reference */
    _toastTimeout: null,

    /**
     * Roll a die with specified sides
     * @param {number} sides - Number of sides on the die
     * @returns {number} The roll result
     */
    roll(sides) {
        const result = Math.floor(Math.random() * sides) + 1;
        this.displayResult(result, `d${sides}`);
        this.addToHistory(`d${sides}`, result);
        return result;
    },

    /**
     * Roll a skill check against a target number
     */
    skillCheck() {
        const targetInput = document.getElementById('skill-target');
        const target = parseInt(targetInput?.value || 50);
        const roll = Math.floor(Math.random() * 100) + 1;

        let outcome = '';
        let outcomeClass = '';

        if (roll === 1) {
            outcome = 'CRITICAL SUCCESS!';
            outcomeClass = 'critical-success';
        } else if (roll === 100) {
            outcome = 'FUMBLE!';
            outcomeClass = 'fumble';
        } else if (roll <= Math.floor(target / 5)) {
            outcome = 'Extreme Success';
            outcomeClass = 'extreme-success';
        } else if (roll <= Math.floor(target / 2)) {
            outcome = 'Hard Success';
            outcomeClass = 'hard-success';
        } else if (roll <= target) {
            outcome = 'Regular Success';
            outcomeClass = 'success';
        } else {
            outcome = 'Failure';
            outcomeClass = 'failure';
        }

        this.displayResult(roll, `d100 vs ${target}`, outcome, outcomeClass);
        this.addToHistory(`Skill (${target}%)`, roll, outcome);
    },

    /**
     * Roll with a bonus die
     */
    rollWithBonus() {
        const units = Math.floor(Math.random() * 10);
        const tens1 = Math.floor(Math.random() * 10) * 10;
        const tens2 = Math.floor(Math.random() * 10) * 10;
        const result = Math.min(tens1, tens2) + units || 100;

        const resultEl = document.getElementById('bonus-penalty-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="result-breakdown">
                    <span class="text-muted">Tens: ${tens1}, ${tens2} (keep lower)</span>
                    <span class="text-muted">Units: ${units}</span>
                </div>
                <div class="result-final text-success fw-bold">Result: ${result === 0 ? 100 : result}</div>
            `;
        }
        this.addToHistory('Bonus Die', result === 0 ? 100 : result);
    },

    /**
     * Roll with a penalty die
     */
    rollWithPenalty() {
        const units = Math.floor(Math.random() * 10);
        const tens1 = Math.floor(Math.random() * 10) * 10;
        const tens2 = Math.floor(Math.random() * 10) * 10;
        const result = Math.max(tens1, tens2) + units || 100;

        const resultEl = document.getElementById('bonus-penalty-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="result-breakdown">
                    <span class="text-muted">Tens: ${tens1}, ${tens2} (keep higher)</span>
                    <span class="text-muted">Units: ${units}</span>
                </div>
                <div class="result-final text-danger fw-bold">Result: ${result === 0 ? 100 : result}</div>
            `;
        }
        this.addToHistory('Penalty Die', result === 0 ? 100 : result);
    },

    /**
     * Display roll result
     */
    displayResult(value, dieType, outcome = '', outcomeClass = '') {
        const resultValue = document.getElementById('result-value');
        const outcomeEl = document.getElementById('dice-outcome');

        if (resultValue) {
            resultValue.textContent = value;
            resultValue.className = 'result-value';
            if (outcomeClass) {
                resultValue.classList.add(outcomeClass);
            }
        }

        if (outcomeEl) {
            outcomeEl.textContent = outcome;
            outcomeEl.className = 'dice-outcome';
            if (outcomeClass) {
                outcomeEl.classList.add(outcomeClass);
            }
        }
    },

    /**
     * Add roll to history
     */
    addToHistory(type, result, outcome = '') {
        this.rollHistory.unshift({ type, result, outcome, time: new Date() });
        if (this.rollHistory.length > 10) {
            this.rollHistory.pop();
        }
        this.renderHistory();
    },

    /**
     * Render roll history
     */
    renderHistory() {
        const historyEl = document.getElementById('roll-history');
        if (!historyEl) return;

        if (this.rollHistory.length === 0) {
            historyEl.innerHTML = '<p class="text-muted small">No rolls yet</p>';
            return;
        }

        historyEl.innerHTML = this.rollHistory.map(roll => `
            <div class="history-item">
                <span class="history-type">${roll.type}</span>
                <span class="history-result">${roll.result}</span>
                ${roll.outcome ? `<span class="history-outcome">${roll.outcome}</span>` : ''}
            </div>
        `).join('');
    },

    // =========================================================================
    // Recovery Calculations
    // =========================================================================

    /**
     * Calculate HP recovery based on rest
     */
    calculateHPRecovery() {
        const amount = parseInt(document.getElementById('hp-rest-amount')?.value || 1);
        const unit = document.getElementById('hp-rest-unit')?.value || 'days';
        const isPulp = document.querySelector('.alert-info')?.textContent?.includes('Pulp');

        let hpRecovered = 0;
        let explanation = '';

        if (isPulp) {
            // Pulp recovery rates
            switch (unit) {
                case 'hours':
                    // With successful First Aid, can heal every hour
                    hpRecovered = amount; // 1 HP per hour with medical attention
                    explanation = `${amount} hour(s) of rest with medical attention`;
                    break;
                case 'days':
                    hpRecovered = amount * 2; // 2 HP per day in Pulp
                    explanation = `${amount} day(s) of rest (2 HP/day in Pulp)`;
                    break;
                case 'weeks':
                    hpRecovered = amount * 14; // 2 HP per day * 7 days
                    explanation = `${amount} week(s) of rest (14 HP/week in Pulp)`;
                    break;
            }
        } else {
            // Standard CoC recovery rates
            switch (unit) {
                case 'days':
                    hpRecovered = amount; // 1 HP per day
                    explanation = `${amount} day(s) of rest (1 HP/day)`;
                    break;
                case 'weeks':
                    hpRecovered = amount * 7; // 1 HP per day * 7 days
                    explanation = `${amount} week(s) of rest (7 HP/week)`;
                    break;
            }
        }

        const resultEl = document.getElementById('hp-recovery-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="recovery-amount text-danger fw-bold">+${hpRecovered} HP</div>
                <div class="recovery-explanation small text-muted">${explanation}</div>
            `;
        }
    },

    /**
     * Calculate Sanity recovery
     */
    calculateSanityRecovery() {
        const amount = parseInt(document.getElementById('san-rest-amount')?.value || 1);
        const unit = document.getElementById('san-rest-unit')?.value || 'months';

        let sanRecovered = 0;
        let explanation = '';

        switch (unit) {
            case 'sessions':
                // Successful therapy session: 1d6 Sanity
                const sessionRolls = [];
                for (let i = 0; i < amount; i++) {
                    sessionRolls.push(Math.floor(Math.random() * 6) + 1);
                }
                sanRecovered = sessionRolls.reduce((a, b) => a + b, 0);
                explanation = `${amount} therapy session(s): ${sessionRolls.join(' + ')} = ${sanRecovered}`;
                break;
            case 'months':
                // 1d6 per month of private care / 1d3 per month of institutional care
                const monthRolls = [];
                for (let i = 0; i < amount; i++) {
                    monthRolls.push(Math.floor(Math.random() * 6) + 1);
                }
                sanRecovered = monthRolls.reduce((a, b) => a + b, 0);
                explanation = `${amount} month(s) of care: ${monthRolls.join(' + ')} = ${sanRecovered}`;
                break;
        }

        const resultEl = document.getElementById('san-recovery-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="recovery-amount text-info fw-bold">+${sanRecovered} SAN</div>
                <div class="recovery-explanation small text-muted">${explanation}</div>
            `;
        }
    },

    /**
     * Roll First Aid recovery (Pulp)
     */
    rollFirstAid() {
        const roll = Math.floor(Math.random() * 6) + 1;
        const total = roll + 1;

        const resultEl = document.getElementById('first-aid-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="recovery-amount text-success fw-bold">+${total} HP</div>
                <div class="recovery-explanation small text-muted">Rolled ${roll} + 1 = ${total}</div>
            `;
        }
        this.addToHistory('First Aid', total, `1d6+1`);
    },

    /**
     * Roll Luck recovery (Pulp)
     */
    rollLuckRecovery() {
        const d1 = Math.floor(Math.random() * 6) + 1;
        const d2 = Math.floor(Math.random() * 6) + 1;
        const total = d1 + d2 + 10;

        const resultEl = document.getElementById('luck-recovery-result');
        if (resultEl) {
            resultEl.innerHTML = `
                <div class="recovery-amount text-primary fw-bold">+${total} Luck</div>
                <div class="recovery-explanation small text-muted">Rolled ${d1} + ${d2} + 10 = ${total}</div>
            `;
        }
        this.addToHistory('Luck Recovery', total, `2d6+10`);
    },

    // =========================================================================
    // Initialization
    // =========================================================================

    init() {
        this.renderHistory();
    }
};

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    HelperPanel.init();
});

// Export globally
window.HelperPanel = HelperPanel;
