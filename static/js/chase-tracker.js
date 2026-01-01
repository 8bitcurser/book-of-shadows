/**
 * Chase Tracker Module - Manages chase sequences
 * @module chase-tracker
 */

const ChaseTracker = {
    participants: [],
    hazards: [],
    trackLength: 10,
    currentRound: 0,
    status: 'setup', // setup, active, ended
    log: [],

    /**
     * Initialize the chase tracker
     */
    init() {
        this.loadState();
        this.renderTrack();
        this.renderParticipants();
        this.updateUI();
        this.updateParticipantDropdown();
    },

    /**
     * Generate a unique ID
     */
    generateId() {
        return 'id_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    },

    /**
     * Add a participant to the chase
     */
    addParticipant() {
        const nameInput = document.getElementById('participant-name');
        const typeSelect = document.getElementById('participant-type');
        const speedInput = document.getElementById('participant-speed');

        const name = nameInput.value.trim();
        if (!name) {
            this.showToast('Please enter a name', 'warning');
            return;
        }

        const participant = {
            id: this.generateId(),
            name: name,
            type: typeSelect.value,
            speed: parseInt(speedInput.value) || 8,
            position: 1,
            status: 'active' // active, caught, escaped, incapacitated
        };

        this.participants.push(participant);
        this.saveState();
        this.renderParticipants();
        this.renderTrack();
        this.updateParticipantDropdown();
        this.addLog(`${name} joined the chase (Speed: ${participant.speed})`);

        // Clear inputs
        nameInput.value = '';
        nameInput.focus();
    },

    /**
     * Remove a participant
     */
    removeParticipant(id) {
        const participant = this.participants.find(p => p.id === id);
        if (participant) {
            this.participants = this.participants.filter(p => p.id !== id);
            this.saveState();
            this.renderParticipants();
            this.renderTrack();
            this.updateParticipantDropdown();
            this.addLog(`${participant.name} removed from chase`);
        }
    },

    /**
     * Add a hazard to the track
     */
    addHazard() {
        const positionInput = document.getElementById('hazard-position');
        const nameInput = document.getElementById('hazard-name');
        const skillInput = document.getElementById('hazard-skill');

        const name = nameInput.value.trim();
        if (!name) {
            this.showToast('Please enter a hazard name', 'warning');
            return;
        }

        const hazard = {
            id: this.generateId(),
            position: parseInt(positionInput.value) || 3,
            name: name,
            skill: skillInput.value.trim() || 'DEX',
            passed: [] // IDs of participants who passed this hazard
        };

        this.hazards.push(hazard);
        this.saveState();
        this.renderTrack();
        this.addLog(`Hazard added: ${name} at position ${hazard.position} (${hazard.skill} check)`);

        // Clear inputs
        nameInput.value = '';
        skillInput.value = '';
    },

    /**
     * Remove a hazard
     */
    removeHazard(id) {
        const hazard = this.hazards.find(h => h.id === id);
        if (hazard) {
            this.hazards = this.hazards.filter(h => h.id !== id);
            this.saveState();
            this.renderTrack();
            this.addLog(`Hazard removed: ${hazard.name}`);
        }
    },

    /**
     * Update track length
     */
    updateTrackLength() {
        const input = document.getElementById('track-length');
        this.trackLength = parseInt(input.value) || 10;
        this.saveState();
        this.renderTrack();
    },

    /**
     * Start the chase
     */
    startChase() {
        if (this.participants.length < 2) {
            this.showToast('Need at least 2 participants', 'warning');
            return;
        }

        this.status = 'active';
        this.currentRound = 1;
        this.saveState();
        this.updateUI();
        this.addLog('--- Chase Started! ---', 'important');
    },

    /**
     * End the chase
     */
    endChase() {
        this.status = 'ended';
        this.saveState();
        this.updateUI();
        this.addLog('--- Chase Ended ---', 'important');
    },

    /**
     * Move to next round
     */
    nextRound() {
        this.currentRound++;
        this.saveState();
        this.updateUI();
        this.addLog(`--- Round ${this.currentRound} ---`, 'round');
    },

    /**
     * Roll movement for all active participants
     */
    rollMovement() {
        this.participants.forEach(p => {
            if (p.status === 'active') {
                // Roll 1d6 + speed modifier
                const roll = Math.floor(Math.random() * 6) + 1;
                const speedMod = Math.floor((p.speed - 5) / 2);
                const movement = Math.max(1, roll + speedMod);

                this.moveParticipant(p.id, movement);
            }
        });
    },

    /**
     * Move a participant
     */
    moveParticipant(id, spaces) {
        const participant = this.participants.find(p => p.id === id);
        if (!participant || participant.status !== 'active') return;

        const oldPosition = participant.position;
        // Allow backward movement but not below 1
        participant.position = Math.max(1, Math.min(this.trackLength, participant.position + spaces));

        // Check for hazards in the path (only when moving forward)
        if (spaces > 0) {
            const hazardsInPath = this.hazards.filter(h =>
                h.position > oldPosition &&
                h.position <= participant.position &&
                !h.passed.includes(id)
            );

            if (hazardsInPath.length > 0) {
                const hazard = hazardsInPath[0];
                participant.position = hazard.position;
                this.addLog(`${participant.name} moved ${spaces} → position ${participant.position}, faces ${hazard.name}!`, 'hazard');
            } else {
                this.addLog(`${participant.name} moved ${spaces} → position ${participant.position}`);
            }
        } else {
            this.addLog(`${participant.name} moved ${spaces} → position ${participant.position}`);
        }

        // Check if escaped
        if (participant.position >= this.trackLength) {
            participant.status = 'escaped';
            this.addLog(`${participant.name} ESCAPED!`, 'success');
        }

        this.saveState();
        this.renderParticipants();
        this.renderTrack();
        this.updateParticipantDropdown();
    },

    /**
     * Manual move from input fields
     */
    manualMove() {
        const participantSelect = document.getElementById('manual-move-participant');
        const spacesInput = document.getElementById('manual-move-spaces');

        const participantId = participantSelect.value;
        const spaces = parseInt(spacesInput.value) || 0;

        if (!participantId) {
            this.showToast('Select a participant first', 'warning');
            return;
        }

        if (spaces === 0) {
            this.showToast('Enter number of spaces to move', 'warning');
            return;
        }

        this.moveParticipant(participantId, spaces);
    },

    /**
     * Update participant dropdown for manual move
     */
    updateParticipantDropdown() {
        const select = document.getElementById('manual-move-participant');
        if (!select) return;

        const currentValue = select.value;
        select.innerHTML = '<option value="">Select...</option>';

        this.participants
            .filter(p => p.status === 'active')
            .forEach(p => {
                const option = document.createElement('option');
                option.value = p.id;
                option.textContent = `${p.name} (Pos: ${p.position})`;
                select.appendChild(option);
            });

        // Restore selection if still valid
        if (currentValue && this.participants.find(p => p.id === currentValue && p.status === 'active')) {
            select.value = currentValue;
        }
    },

    /**
     * Resolve hazard check for a participant
     */
    resolveHazard(participantId, hazardId, success) {
        const participant = this.participants.find(p => p.id === participantId);
        const hazard = this.hazards.find(h => h.id === hazardId);

        if (!participant || !hazard) return;

        if (success) {
            hazard.passed.push(participantId);
            this.addLog(`${participant.name} passed the ${hazard.name}!`, 'success');
        } else {
            participant.position = Math.max(1, participant.position - 1);
            this.addLog(`${participant.name} failed the ${hazard.name} and fell back!`, 'failure');
        }

        this.saveState();
        this.renderParticipants();
        this.renderTrack();
    },

    /**
     * Set participant status
     */
    setStatus(id, status) {
        const participant = this.participants.find(p => p.id === id);
        if (participant) {
            participant.status = status;
            this.saveState();
            this.renderParticipants();
            this.renderTrack();
            this.addLog(`${participant.name} is now ${status}`);
        }
    },

    /**
     * Render the chase track
     */
    renderTrack() {
        const container = document.getElementById('chase-track');
        if (!container) return;

        let html = '<div class="track-positions">';

        for (let i = 1; i <= this.trackLength; i++) {
            const hazard = this.hazards.find(h => h.position === i);
            const participantsHere = this.participants.filter(p => p.position === i && p.status === 'active');

            const isStart = i === 1;
            const isEnd = i === this.trackLength;

            html += `<div class="track-position ${hazard ? 'has-hazard' : ''} ${isStart ? 'start' : ''} ${isEnd ? 'end' : ''}">`;
            html += `<div class="position-number">${i}</div>`;

            if (hazard) {
                html += `<div class="hazard-marker" title="${hazard.name} (${hazard.skill})">
                    <i class="bi bi-exclamation-triangle-fill"></i>
                </div>`;
            }

            if (participantsHere.length > 0) {
                html += '<div class="position-participants">';
                participantsHere.forEach(p => {
                    const typeClass = p.type === 'enemy' ? 'enemy' : p.type === 'npc' ? 'npc' : 'investigator';
                    html += `<div class="track-participant ${typeClass}" title="${p.name}">
                        ${p.name.charAt(0).toUpperCase()}
                    </div>`;
                });
                html += '</div>';
            }

            html += '</div>';
        }

        html += '</div>';

        // Add legend
        html += `<div class="track-legend mt-3">
            <span class="legend-item"><span class="track-participant investigator">I</span> Investigator</span>
            <span class="legend-item"><span class="track-participant enemy">E</span> Enemy</span>
            <span class="legend-item"><span class="track-participant npc">N</span> NPC</span>
            <span class="legend-item"><i class="bi bi-exclamation-triangle-fill text-warning"></i> Hazard</span>
        </div>`;

        container.innerHTML = html;
    },

    /**
     * Render participants list
     */
    renderParticipants() {
        const container = document.getElementById('participants-list');
        const countEl = document.getElementById('participant-count');

        if (!container) return;

        if (this.participants.length === 0) {
            container.innerHTML = `<div class="text-center text-muted p-4">
                <i class="bi bi-person-dash display-6"></i>
                <p class="mt-2 mb-0">No participants yet</p>
            </div>`;
            if (countEl) countEl.textContent = '0';
            return;
        }

        if (countEl) countEl.textContent = this.participants.length;

        let html = '';
        this.participants.forEach(p => {
            const typeClass = p.type === 'enemy' ? 'bg-danger' : p.type === 'npc' ? 'bg-info' : 'bg-primary';
            const statusClass = p.status === 'escaped' ? 'text-success' :
                               p.status === 'caught' ? 'text-danger' :
                               p.status === 'incapacitated' ? 'text-muted' : '';

            html += `<div class="participant-item ${statusClass}">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <span class="badge ${typeClass} me-2">${p.type.charAt(0).toUpperCase()}</span>
                        <strong>${p.name}</strong>
                        <small class="text-muted ms-2">SPD: ${p.speed} | POS: ${p.position}</small>
                    </div>
                    <div class="btn-group btn-group-sm">
                        ${this.status === 'active' && p.status === 'active' ? `
                            <button class="btn btn-outline-primary" onclick="ChaseTracker.moveParticipant('${p.id}', 1)" title="Move +1">
                                <i class="bi bi-arrow-right"></i>
                            </button>
                            <button class="btn btn-outline-secondary" onclick="ChaseTracker.moveParticipant('${p.id}', -1)" title="Move -1">
                                <i class="bi bi-arrow-left"></i>
                            </button>
                        ` : ''}
                        <button class="btn btn-outline-danger" onclick="ChaseTracker.removeParticipant('${p.id}')" title="Remove">
                            <i class="bi bi-x"></i>
                        </button>
                    </div>
                </div>
                ${p.status !== 'active' ? `<small class="badge bg-${p.status === 'escaped' ? 'success' : 'secondary'} mt-1">${p.status.toUpperCase()}</small>` : ''}
            </div>`;
        });

        container.innerHTML = html;
    },

    /**
     * Update UI state
     */
    updateUI() {
        const statusEl = document.getElementById('chase-status');
        const roundEl = document.getElementById('chase-round');
        const startBtn = document.getElementById('btn-start');
        const nextRoundBtn = document.getElementById('btn-next-round');
        const rollMovementBtn = document.getElementById('btn-roll-movement');

        if (statusEl) {
            statusEl.textContent = this.status.charAt(0).toUpperCase() + this.status.slice(1);
            statusEl.className = 'badge ' + (this.status === 'active' ? 'bg-success' : 'bg-secondary');
        }

        if (roundEl) {
            roundEl.textContent = 'Round ' + this.currentRound;
        }

        if (startBtn) startBtn.disabled = this.status !== 'setup';
        if (nextRoundBtn) nextRoundBtn.disabled = this.status !== 'active';
        if (rollMovementBtn) rollMovementBtn.disabled = this.status !== 'active';
    },

    /**
     * Add entry to the log
     */
    addLog(message, type = 'normal') {
        const container = document.getElementById('chase-log');
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
        if (this.log.length > 50) this.log.pop();
    },

    /**
     * Clear the log
     */
    clearLog() {
        const container = document.getElementById('chase-log');
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
            'Reset Chase',
            'Are you sure you want to reset the entire chase? This will clear all participants and hazards.',
            () => {
                this.participants = [];
                this.hazards = [];
                this.currentRound = 0;
                this.status = 'setup';
                this.log = [];

                this.saveState();
                this.renderTrack();
                this.renderParticipants();
                this.updateUI();
                this.updateParticipantDropdown();
                this.clearLog();
                this.addLog('Chase tracker reset');
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
            participants: this.participants,
            hazards: this.hazards,
            trackLength: this.trackLength,
            currentRound: this.currentRound,
            status: this.status
        };
        sessionStorage.setItem('chaseTracker', JSON.stringify(state));
    },

    /**
     * Load state from sessionStorage
     */
    loadState() {
        const saved = sessionStorage.getItem('chaseTracker');
        if (saved) {
            try {
                const state = JSON.parse(saved);
                this.participants = state.participants || [];
                this.hazards = state.hazards || [];
                this.trackLength = state.trackLength || 10;
                this.currentRound = state.currentRound || 0;
                this.status = state.status || 'setup';

                // Update track length input
                const trackLengthInput = document.getElementById('track-length');
                if (trackLengthInput) trackLengthInput.value = this.trackLength;
            } catch (e) {
                console.error('Failed to load chase state:', e);
            }
        }
    },

    /**
     * Show toast notification
     */
    showToast(message, type = 'info') {
        // Simple alert for now, can be enhanced
        const icons = { warning: '⚠️', success: '✅', info: 'ℹ️', error: '❌' };
        console.log(`${icons[type] || ''} ${message}`);
    }
};

// Export globally
window.ChaseTracker = ChaseTracker;
