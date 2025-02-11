const characterUtils = {
    getCurrentCharacter() {
        const hiddenInput = document.getElementById('currentCharacter');
        if (!hiddenInput || !hiddenInput.value) {
            return null;
        }
        return JSON.parse(hiddenInput.value);
    },

    async exportPDF(evt, key) {
        try {

            const response = await fetch('/api/investigator/export/' + key, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({})
            });

            if (!response.ok) {
                throw new Error('PDF export failed');
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = key + ".pdf";
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error exporting PDF:', error);
            alert('Failed to export PDF. Please try again.');
        }
    },

    async importInvestigators() {
        try {
            const key = document.getElementById('importCode').value
            const response = await fetch('/api/investigator/list/import/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "ImportCode": key
                })
            });

            if (!response.ok) {
                throw new Error('PDF export failed');
            }
            // Close modal and trigger HTMX refresh
            document.getElementById('importModal').classList.add('hidden');
            htmx.trigger('body', 'import');

        } catch (error) {
            console.error('Importing Investigators:', error);
            alert('Failed to export PDF. Please try again.');
        }
    },

    async handleSkillToggleCheck(input){
        const skillName = input.dataset.skill;
        await this.updateInvestigator(
            "skill_check",
            skillName,
            true
        )
    },

    async handleSkillNameChange(input){
        const skillName = input.dataset.skill;
        await this.updateInvestigator(
            "skill_name",
            skillName,
            input.value
        )
    },

    async recalculateValues(input, type) {
        let value = 0
        value = parseInt(input.value) || 0;

        const container = input.parentElement;
        const halfSpan = container.querySelector('[data-half]');
        const fifthSpan = container.querySelector('[data-fifth]');

        if (halfSpan) halfSpan.textContent = Math.floor(value / 2);
        if (fifthSpan) fifthSpan.textContent = Math.floor(value / 5);

        if (type === 'attribute') {
            const attrName = input.dataset.attr;
            if (characterData.attributes?.[attrName]) {
                characterData.attributes[attrName].value = value;
            }
            await this.updateInvestigator(
                "combat",
                attrName,
                value
            )
        } else if (type === 'skill') {
            const skillName = input.dataset.skill;
            const prevValue = parseInt(input.dataset.skillvalue);

            const difference = value - prevValue;

            let currentPoints = document.getElementById("archetype-points");
            if (!currentPoints) {
                currentPoints = document.getElementById("occupation-points");
            }
            if (!currentPoints) {
                currentPoints = document.getElementById("general-points");
            }
            if (currentPoints) {
                let val= parseInt(currentPoints.innerText) || 0;
                let newPoints = val - difference;

                // Don't allow negative points
                if (newPoints < 0) {
                    input.value = prevValue;
                    return;
                }
                currentPoints.innerText = newPoints.toString();
                input.dataset.skillvalue = value.toString();
                this.updateArchetypeConfirmButton(newPoints);

                // Disable/enable inputs based on points
                this.updateArchetypeInputs(newPoints);
            }
            input.dataset.skillvalue = value.toString();
                await this.updateInvestigator(
                "skills",
                skillName,
                value
            )

        }
    },

    async updatePersonalInfo(input) {
        const field = input.dataset.field;
        const value = field === 'age' ? parseInt(input.value) || 0 : input.value;
        const characterData = this.getCurrentCharacter();
        if (!characterData) return;

        switch(field) {
            case 'Name':
                characterData.Investigators_Name = value;
                break;
            case 'Age':
                characterData.Age = value;
                break;
            case 'Residence':
                characterData.Residence = value;
                break;
            case 'Birthplace':
                characterData.Birthplace = value;
                break;
        }
        await this.updateInvestigator(
            "personalInfo",
            field,
            value
        )
        document.getElementById('currentCharacter').value = JSON.stringify(characterData);
    },

    initializeEventListeners() {
        const buttonHandlers = {
            'exportPdf': () => this.exportPDF()
        };

        Object.entries(buttonHandlers).forEach(([id, handler]) => {
            const button = document.getElementById(id);
            if (button) {
                button.addEventListener('click', (e) => {
                    e.preventDefault();
                    handler();
                });
            }
        });

        this.initializeInputHandlers();
    },

    initializeInputHandlers() {
        document.querySelectorAll('input[type="number"]').forEach(input => {
            input.addEventListener('input', () => {
                this.recalculateValues(input, input.dataset.skill ? 'skill' : 'attribute');
            });
        });

        document.querySelectorAll('input[data-field]').forEach(input => {
            input.addEventListener('change', () => {
                this.updatePersonalInfo(input);
            });
        });
    },

    async updateInvestigator(section, field, value) {
        try {
            const cookieId = document.querySelector('input[data-field="Name"]').id
            const response = await fetch(`/api/investigator/update/${cookieId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    section: section,
                    field: field,
                    value: value
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return null
        } catch (error) {
            console.error('Error updating investigator:', error);
            throw error;
        }
    },

    showDescription(value, targetId) {
        const select = event.target;
        const selectedOption = select.options[select.selectedIndex];
        const description = selectedOption.getAttribute('data-description');
        const target = document.getElementById(targetId);

        if (description) {
            target.textContent = description;
            target.style.opacity = "1";
        } else {
            target.textContent = "";
            target.style.opacity = "0";
        }
    },

    rollAttribute(button, formula) {
        const input = button.previousElementSibling;
        let result = 0;

        if (formula === '3d6x5') {
            // Roll 3d6 * 5
            result = (rollDice(3, 6)) * 5;
        } else if (formula === '2d6p6x5') {
            // Roll (2d6 + 6) * 5
            result = (rollDice(2, 6) + 6) * 5;
        }

        input.value = result;
    },

    updateArchetypeConfirmButton(points) {
        const confirmContainer = document.getElementById('confirm-archetype-container');
        if (confirmContainer) {
            confirmContainer.style.opacity = points === 0 ? "1" : "0";
            confirmContainer.style.pointerEvents = points === 0 ? "auto" : "none";
        }
    },

    updateArchetypeInputs(points) {
        const inputs = document.querySelectorAll('input[data-skill]');
        inputs.forEach(input => {
            if (points === 0) {
                // When points are 0, don't allow increasing values
                const currentValue = parseInt(input.value) || 0;
                input.max = currentValue;
            } else {
                // Remove max restriction when points are available
                input.removeAttribute('max');
            }
        });
    }


};


function rollDice(numDice, sides) {
    let total = 0;
    for (let i = 0; i < numDice; i++) {
        total += Math.floor(Math.random() * sides) + 1;
    }
    return total;
}



// Initialize only once when the page loads
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        characterUtils.initializeEventListeners();
    });

} else {
    characterUtils.initializeEventListeners();
}


// Make utils available globally
window.characterUtils = characterUtils;