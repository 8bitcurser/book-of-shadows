/**
 * Combat Tracker Module - Manages combat encounters
 * @module combat-tracker
 */

const CombatTracker = {
    combatants: [],
    actions: [],
    currentRound: 0,
    currentTurnIndex: -1,
    status: 'setup', // setup, active, ended
    log: [],

    /**
     * Initialize the combat tracker
     */
    init() {
        this.loadState();
        this.renderInitiativeList();
        this.updateUI();
    },

    /**
     * Generate a unique ID
     */
    generateId() {
        return 'id_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    },

    /**
     * Add a combatant
     */
    addCombatant() {
        const nameInput = document.getElementById('combatant-name');
        const typeSelect = document.getElementById('combatant-type');
        const hpInput = document.getElementById('combatant-hp');
        const dexInput = document.getElementById('combatant-dex');

        const name = nameInput.value.trim();
        if (!name) {
            this.showToast('Please enter a name', 'warning');
            return;
        }

        const combatant = {
            id: this.generateId(),
            name: name,
            type: typeSelect.value,
            hp: parseInt(hpInput.value) || 12,
            maxHp: parseInt(hpInput.value) || 12,
            dex: parseInt(dexInput.value) || 50,
            initiative: 0,
            status: 'active', // active, unconscious, dead, fled
            conditions: []
        };

        this.combatants.push(combatant);
        this.saveState();
        this.renderInitiativeList();
        this.updateTargetDropdown();
        this.addLog(`${name} joined combat (HP: ${combatant.hp}, DEX: ${combatant.dex})`);

        // Clear name input
        nameInput.value = '';
        nameInput.focus();
    },

    /**
     * Remove a combatant
     */
    removeCombatant(id) {
        const combatant = this.combatants.find(c => c.id === id);
        if (combatant) {
            this.combatants = this.combatants.filter(c => c.id !== id);
            this.saveState();
            this.renderInitiativeList();
            this.updateActiveCombatant();
            this.updateTargetDropdown();
            this.addLog(`${combatant.name} removed from combat`);
        }
    },

    /**
     * Roll initiative for a single combatant
     */
    rollInitiative(id) {
        const combatant = this.combatants.find(c => c.id === id);
        if (combatant) {
            // Roll d100 against DEX
            const roll = Math.floor(Math.random() * 100) + 1;
            combatant.initiative = combatant.dex + (roll <= combatant.dex ? 50 : 0);
            this.sortByInitiative();
            this.saveState();
            this.renderInitiativeList();
            this.addLog(`${combatant.name} rolled initiative: ${combatant.initiative}`);
        }
    },

    /**
     * Roll initiative for all combatants
     */
    rollAllInitiative() {
        this.combatants.forEach(c => {
            const roll = Math.floor(Math.random() * 100) + 1;
            c.initiative = c.dex + (roll <= c.dex ? 50 : 0);
        });
        this.sortByInitiative();
        this.saveState();
        this.renderInitiativeList();
        this.addLog('Initiative rolled for all combatants', 'round');
    },

    /**
     * Sort combatants by initiative
     */
    sortByInitiative() {
        this.combatants.sort((a, b) => b.initiative - a.initiative);
    },

    /**
     * Start combat
     */
    startCombat() {
        if (this.combatants.length < 1) {
            this.showToast('Need at least 1 combatant', 'warning');
            return;
        }

        // Auto-roll initiative if not set
        const noInitiative = this.combatants.every(c => c.initiative === 0);
        if (noInitiative) {
            this.rollAllInitiative();
        }

        this.status = 'active';
        this.currentRound = 1;
        this.currentTurnIndex = 0;
        this.saveState();
        this.updateUI();
        this.updateActiveCombatant();
        this.addLog('--- Combat Started! ---', 'important');
        this.addLog(`--- Round ${this.currentRound} ---`, 'round');
    },

    /**
     * End combat
     */
    endCombat() {
        this.status = 'ended';
        this.currentTurnIndex = -1;
        this.saveState();
        this.updateUI();
        this.updateActiveCombatant();
        this.addLog('--- Combat Ended ---', 'important');
    },

    /**
     * Move to next turn
     */
    nextTurn() {
        if (this.status !== 'active') return;

        // Find next active combatant
        let nextIndex = this.currentTurnIndex + 1;
        let looped = false;

        while (true) {
            if (nextIndex >= this.combatants.length) {
                nextIndex = 0;
                looped = true;
            }

            const combatant = this.combatants[nextIndex];
            if (combatant && combatant.status === 'active') {
                break;
            }

            nextIndex++;

            // Prevent infinite loop if no active combatants
            if (nextIndex === this.currentTurnIndex) {
                this.addLog('No active combatants remaining', 'important');
                return;
            }
        }

        // If we looped, it's a new round
        if (looped && nextIndex <= this.currentTurnIndex) {
            this.nextRound();
        }

        this.currentTurnIndex = nextIndex;
        this.saveState();
        this.renderInitiativeList();
        this.updateActiveCombatant();

        const current = this.combatants[this.currentTurnIndex];
        if (current) {
            this.addLog(`${current.name}'s turn`);
        }
    },

    /**
     * Move to next round
     */
    nextRound() {
        this.currentRound++;
        this.currentTurnIndex = 0;

        // Find first active combatant
        while (this.currentTurnIndex < this.combatants.length &&
               this.combatants[this.currentTurnIndex].status !== 'active') {
            this.currentTurnIndex++;
        }

        this.saveState();
        this.updateUI();
        this.addLog(`--- Round ${this.currentRound} ---`, 'round');
    },

    /**
     * Record an action
     */
    recordAction(actionType) {
        if (this.status !== 'active' || this.currentTurnIndex < 0) return;

        const combatant = this.combatants[this.currentTurnIndex];
        if (!combatant) return;

        // Get target for targeted actions
        const targetedActions = ['attack', 'spell', 'item'];
        let targetId = null;
        let targetName = null;
        let target = null;

        if (targetedActions.includes(actionType)) {
            const targetSelect = document.getElementById('action-target');
            if (targetSelect && targetSelect.value) {
                targetId = targetSelect.value;
                target = this.combatants.find(c => c.id === targetId);
                targetName = target ? target.name : null;
            }
        }

        // Get damage values
        const damageInput = document.getElementById('action-damage');
        const fightbackInput = document.getElementById('action-fightback');
        const damageToTarget = parseInt(damageInput?.value) || 0;
        const fightbackDamage = parseInt(fightbackInput?.value) || 0;

        const action = {
            round: this.currentRound,
            combatantId: combatant.id,
            combatantName: combatant.name,
            actionType: actionType,
            targetId: targetId,
            targetName: targetName,
            damageDealt: damageToTarget,
            damageReceived: fightbackDamage,
            timestamp: Date.now()
        };

        this.actions.push(action);

        const actionNames = {
            attack: '⚔️ Attack',
            defend: '🛡️ Defend',
            dodge: '↪️ Dodge',
            flee: '🏃 Flee',
            spell: '✨ Cast Spell',
            item: '🎒 Use Item'
        };

        let logMessage = `${combatant.name}: ${actionNames[actionType] || actionType}`;
        if (targetName) {
            logMessage += ` → ${targetName}`;
        }
        this.addLog(logMessage);

        // Apply damage to target
        if (target && damageToTarget > 0) {
            this.adjustHP(targetId, -damageToTarget);
        }

        // Apply fight back damage to attacker
        if (fightbackDamage > 0) {
            this.adjustHP(combatant.id, -fightbackDamage);
        }

        // Clear inputs after action
        const targetSelect = document.getElementById('action-target');
        if (targetSelect) targetSelect.value = '';
        if (damageInput) damageInput.value = '0';
        if (fightbackInput) fightbackInput.value = '0';

        // Auto-advance to next turn
        this.nextTurn();
    },

    /**
     * Adjust HP for a combatant
     */
    adjustHP(id, delta) {
        const combatant = this.combatants.find(c => c.id === id);
        if (!combatant) return;

        const oldHP = combatant.hp;
        combatant.hp = Math.max(0, Math.min(combatant.maxHp, combatant.hp + delta));

        // Check for major wound (half HP or more in one hit)
        if (delta < 0 && Math.abs(delta) >= Math.floor(combatant.maxHp / 2)) {
            if (!combatant.conditions.includes('major wound')) {
                combatant.conditions.push('major wound');
                this.addLog(`${combatant.name} suffers a MAJOR WOUND!`, 'failure');
            }
        }

        // Check for unconscious/dead
        if (combatant.hp <= 0) {
            combatant.status = 'unconscious';
            this.addLog(`${combatant.name} is UNCONSCIOUS!`, 'failure');
        } else if (oldHP <= 0 && combatant.hp > 0) {
            combatant.status = 'active';
            this.addLog(`${combatant.name} regains consciousness`);
        }

        if (delta < 0) {
            this.addLog(`${combatant.name} takes ${Math.abs(delta)} damage (HP: ${combatant.hp}/${combatant.maxHp})`);
        } else if (delta > 0) {
            this.addLog(`${combatant.name} heals ${delta} (HP: ${combatant.hp}/${combatant.maxHp})`, 'success');
        }

        this.saveState();
        this.renderInitiativeList();
        this.updateActiveCombatant();
    },

    /**
     * Adjust HP for the active combatant
     */
    adjustActiveHP(delta) {
        if (this.currentTurnIndex >= 0 && this.combatants[this.currentTurnIndex]) {
            this.adjustHP(this.combatants[this.currentTurnIndex].id, delta);
        }
    },

    /**
     * Heal active combatant using amount input
     */
    healActive() {
        const amountInput = document.getElementById('active-hp-amount');
        const amount = parseInt(amountInput.value) || 1;
        this.adjustActiveHP(amount);
    },

    /**
     * Damage active combatant using amount input
     */
    damageActive() {
        const amountInput = document.getElementById('active-hp-amount');
        const amount = parseInt(amountInput.value) || 1;
        this.adjustActiveHP(-amount);
    },

    /**
     * Toggle major wound condition for active combatant
     */
    toggleActiveMajorWound() {
        if (this.currentTurnIndex < 0) return;
        const combatant = this.combatants[this.currentTurnIndex];
        if (!combatant) return;

        const checkbox = document.getElementById('active-major-wound');
        const hasMajorWound = combatant.conditions.includes('major wound');

        if (checkbox.checked && !hasMajorWound) {
            combatant.conditions.push('major wound');
            this.addLog(`${combatant.name} suffers a MAJOR WOUND!`, 'failure');
        } else if (!checkbox.checked && hasMajorWound) {
            combatant.conditions = combatant.conditions.filter(c => c !== 'major wound');
            this.addLog(`${combatant.name}'s major wound treated`);
        }

        this.saveState();
        this.renderInitiativeList();
        this.updateActiveCombatant();
    },

    /**
     * Set combatant status
     */
    setStatus(id, status) {
        const combatant = this.combatants.find(c => c.id === id);
        if (combatant) {
            combatant.status = status;
            this.saveState();
            this.renderInitiativeList();
            this.updateActiveCombatant();
            this.addLog(`${combatant.name} is now ${status}`);
        }
    },

    /**
     * Render initiative list
     */
    renderInitiativeList() {
        const container = document.getElementById('initiative-list');
        if (!container) return;

        if (this.combatants.length === 0) {
            container.innerHTML = `<div class="text-center text-muted p-4">
                <i class="bi bi-hourglass display-6"></i>
                <p class="mt-2 mb-0">No combatants yet</p>
            </div>`;
            return;
        }

        let html = '';
        this.combatants.forEach((c, index) => {
            const isActive = this.status === 'active' && index === this.currentTurnIndex;
            const typeClass = c.type === 'enemy' ? 'bg-danger' : c.type === 'npc' ? 'bg-info' : 'bg-primary';
            const statusClass = c.status === 'unconscious' ? 'combatant-unconscious' :
                               c.status === 'dead' ? 'combatant-dead' :
                               c.status === 'fled' ? 'combatant-fled' : '';

            const hpPercent = Math.max(0, (c.hp / c.maxHp) * 100);
            const hpColor = hpPercent > 50 ? '#63c74d' : hpPercent > 25 ? '#f7b731' : '#e84a5f';

            html += `<div class="initiative-item ${isActive ? 'active-turn' : ''} ${statusClass}">
                <div class="d-flex justify-content-between align-items-center">
                    <div class="d-flex align-items-center">
                        <span class="initiative-number me-2">${c.initiative || '-'}</span>
                        <span class="badge ${typeClass} me-2">${c.type.charAt(0).toUpperCase()}</span>
                        <strong>${c.name}</strong>
                        ${c.conditions.length > 0 ? `<span class="badge bg-warning ms-2">${c.conditions.join(', ')}</span>` : ''}
                    </div>
                    <div class="d-flex align-items-center gap-2">
                        <div class="hp-mini">
                            <div class="hp-mini-bar" style="width: ${hpPercent}%; background: ${hpColor}"></div>
                            <span class="hp-mini-text">${c.hp}/${c.maxHp}</span>
                        </div>
                        <div class="btn-group btn-group-sm">
                            <button class="btn btn-outline-success" onclick="CombatTracker.adjustHP('${c.id}', 1)" title="Heal">
                                <i class="bi bi-plus"></i>
                            </button>
                            <button class="btn btn-outline-danger" onclick="CombatTracker.adjustHP('${c.id}', -1)" title="Damage">
                                <i class="bi bi-dash"></i>
                            </button>
                            <button class="btn btn-outline-secondary" onclick="CombatTracker.removeCombatant('${c.id}')" title="Remove">
                                <i class="bi bi-x"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>`;
        });

        container.innerHTML = html;
    },

    /**
     * Update active combatant display
     */
    updateActiveCombatant() {
        const card = document.getElementById('active-combatant-card');
        const nameEl = document.getElementById('active-combatant-name');
        const typeEl = document.getElementById('active-combatant-type');
        const hpEl = document.getElementById('active-combatant-hp');
        const maxHpEl = document.getElementById('active-combatant-maxhp');
        const hpBar = document.getElementById('active-combatant-hp-bar');
        const majorWoundBadge = document.getElementById('active-combatant-major-wound');
        const majorWoundCheckbox = document.getElementById('active-major-wound');

        if (this.status !== 'active' || this.currentTurnIndex < 0) {
            if (card) card.style.display = 'none';
            return;
        }

        const combatant = this.combatants[this.currentTurnIndex];
        if (!combatant) {
            if (card) card.style.display = 'none';
            return;
        }

        if (card) card.style.display = 'block';
        if (nameEl) nameEl.textContent = combatant.name;
        if (typeEl) {
            typeEl.textContent = combatant.type;
            typeEl.className = 'badge ' + (combatant.type === 'enemy' ? 'bg-danger' :
                                          combatant.type === 'npc' ? 'bg-info' : 'bg-primary');
        }
        if (hpEl) hpEl.textContent = combatant.hp;
        if (maxHpEl) maxHpEl.textContent = combatant.maxHp;
        if (hpBar) {
            const percent = Math.max(0, (combatant.hp / combatant.maxHp) * 100);
            hpBar.style.width = percent + '%';
            hpBar.style.background = percent > 50 ? '#63c74d' : percent > 25 ? '#f7b731' : '#e84a5f';
        }

        // Update major wound display
        const hasMajorWound = combatant.conditions.includes('major wound');
        if (majorWoundBadge) {
            majorWoundBadge.style.display = hasMajorWound ? 'inline' : 'none';
        }
        if (majorWoundCheckbox) {
            majorWoundCheckbox.checked = hasMajorWound;
        }

        // Update target dropdown
        this.updateTargetDropdown();
    },

    /**
     * Update target dropdown for actions
     */
    updateTargetDropdown() {
        const select = document.getElementById('action-target');
        if (!select) return;

        const currentValue = select.value;
        select.innerHTML = '<option value="">Select target...</option>';

        // Exclude the active combatant from targets
        const activeId = this.currentTurnIndex >= 0 ? this.combatants[this.currentTurnIndex]?.id : null;

        this.combatants
            .filter(c => c.id !== activeId && c.status === 'active')
            .forEach(c => {
                const typeLabel = c.type === 'enemy' ? '👹' : c.type === 'npc' ? '👤' : '🔍';
                const option = document.createElement('option');
                option.value = c.id;
                option.textContent = `${typeLabel} ${c.name} (HP: ${c.hp}/${c.maxHp})`;
                select.appendChild(option);
            });

        // Restore selection if still valid
        if (currentValue && this.combatants.find(c => c.id === currentValue && c.status === 'active')) {
            select.value = currentValue;
        }
    },

    /**
     * Update UI state
     */
    updateUI() {
        const statusEl = document.getElementById('combat-status');
        const roundEl = document.getElementById('combat-round');
        const startBtn = document.getElementById('btn-start-combat');
        const nextRoundBtn = document.getElementById('btn-next-combat-round');
        const nextTurnBtn = document.getElementById('btn-next-turn');

        if (statusEl) {
            statusEl.textContent = this.status.charAt(0).toUpperCase() + this.status.slice(1);
            statusEl.className = 'badge ' + (this.status === 'active' ? 'bg-danger' : 'bg-secondary');
        }

        if (roundEl) {
            roundEl.textContent = 'Round ' + this.currentRound;
        }

        if (startBtn) startBtn.disabled = this.status !== 'setup';
        if (nextRoundBtn) nextRoundBtn.disabled = this.status !== 'active';
        if (nextTurnBtn) nextTurnBtn.disabled = this.status !== 'active';
    },

    /**
     * Add entry to the log
     */
    addLog(message, type = 'normal') {
        const container = document.getElementById('combat-log');
        if (!container) return;

        const entry = document.createElement('div');
        entry.className = `log-entry log-${type}`;

        const time = new Date().toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
        entry.innerHTML = `<span class="log-time">${time}</span> ${message}`;

        // Add to beginning
        if (container.firstChild) {
            container.insertBefore(entry, container.firstChild);
        } else {
            container.appendChild(entry);
        }

        this.log.unshift({ message, type, time: Date.now() });
        if (this.log.length > 100) this.log.pop();
    },

    /**
     * Clear the log
     */
    clearLog() {
        const container = document.getElementById('combat-log');
        if (container) {
            container.innerHTML = '<div class="log-entry text-muted"><i class="bi bi-info-circle me-1"></i>Log cleared.</div>';
        }
        this.log = [];
    },

    /**
     * Reset the tracker
     */
    reset() {
        this.showConfirmModal(
            'Reset Combat',
            'Are you sure you want to reset combat? This will clear all combatants and history.',
            () => {
                this.combatants = [];
                this.actions = [];
                this.currentRound = 0;
                this.currentTurnIndex = -1;
                this.status = 'setup';
                this.log = [];

                this.saveState();
                this.renderInitiativeList();
                this.updateActiveCombatant();
                this.updateTargetDropdown();
                this.updateUI();
                this.clearLog();
                this.addLog('Combat tracker reset');
            }
        );
    },

    /**
     * Show confirmation modal
     */
    showConfirmModal(title, message, onConfirm) {
        const modal = document.getElementById('confirm-modal');
        const titleEl = document.getElementById('confirm-modal-title');
        const bodyEl = document.getElementById('confirm-modal-body');
        const btnEl = document.getElementById('confirm-modal-btn');

        if (!modal) {
            // Fallback to confirm if modal not found
            if (confirm(message)) onConfirm();
            return;
        }

        titleEl.innerHTML = `<i class="bi bi-exclamation-triangle text-warning me-2"></i>${title}`;
        bodyEl.innerHTML = `<p>${message}</p>`;

        // Remove old event listener and add new one
        const newBtn = btnEl.cloneNode(true);
        btnEl.parentNode.replaceChild(newBtn, btnEl);
        newBtn.addEventListener('click', onConfirm);

        const bsModal = new bootstrap.Modal(modal);
        bsModal.show();
    },

    /**
     * Save state to sessionStorage
     */
    saveState() {
        const state = {
            combatants: this.combatants,
            actions: this.actions,
            currentRound: this.currentRound,
            currentTurnIndex: this.currentTurnIndex,
            status: this.status
        };
        sessionStorage.setItem('combatTracker', JSON.stringify(state));
    },

    /**
     * Load state from sessionStorage
     */
    loadState() {
        const saved = sessionStorage.getItem('combatTracker');
        if (saved) {
            try {
                const state = JSON.parse(saved);
                this.combatants = state.combatants || [];
                this.actions = state.actions || [];
                this.currentRound = state.currentRound || 0;
                this.currentTurnIndex = state.currentTurnIndex ?? -1;
                this.status = state.status || 'setup';
            } catch (e) {
                console.error('Failed to load combat state:', e);
            }
        }
    },

    /**
     * Show toast notification
     */
    showToast(message, type = 'info') {
        const icons = { warning: '⚠️', success: '✅', info: 'ℹ️', error: '❌' };
        console.log(`${icons[type] || ''} ${message}`);
    }
};

// Export globally
window.CombatTracker = CombatTracker;
