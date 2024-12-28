const characterUtils = {
    getCurrentCharacter() {
        const hiddenInput = document.getElementById('currentCharacter');
        if (!hiddenInput || !hiddenInput.value) {
            return null;
        }
        return JSON.parse(hiddenInput.value);
    },

    getCurrentUIState() {
        const baseCharacter = JSON.parse(document.getElementById('currentCharacter').value);
        const flatState = {
            // Basic info
            Investigators_Name: document.querySelector('input[data-field="name"]').value,
            Age: String(parseInt(document.querySelector('input[data-field="age"]').value) || 0),
            Residence: document.querySelector('input[data-field="residence"]').value,
            Birthplace: document.querySelector('input[data-field="birthplace"]').value,
            Occupation: baseCharacter.Occupation?.name || '',
            Archetype: baseCharacter.Archetype?.name || '',
            MOV: String(baseCharacter.MOV || 0),
            Build: baseCharacter.Build || '',
            DamageBonus: baseCharacter.DamageBonus || '',
        };

        // Add attributes with their half and fifth values
        const attributes = {
            STR: parseInt(document.querySelector('span[data-attr="STR"]')?.textContent) || 0,
            DEX: parseInt(document.querySelector('span[data-attr="DEX"]')?.textContent) || 0,
            INT: parseInt(document.querySelector('span[data-attr="INT"]')?.textContent) || 0,
            CON: parseInt(document.querySelector('span[data-attr="CON"]')?.textContent) || 0,
            APP: parseInt(document.querySelector('span[data-attr="APP"]')?.textContent) || 0,
            POW: parseInt(document.querySelector('span[data-attr="POW"]')?.textContent) || 0,
            SIZ: parseInt(document.querySelector('span[data-attr="SIZ"]')?.textContent) || 0,
            EDU: parseInt(document.querySelector('span[data-attr="EDU"]')?.textContent) || 0,
            CurrentMagic: parseInt(document.querySelector('input[data-attr="CurrentMagic"]')?.value) || 0,
            CurrentHP: parseInt(document.querySelector('input[data-attr="CurrentHP"]')?.value) || 0,
            CurrentSanity: parseInt(document.querySelector('input[data-attr="CurrentSanity"]')?.value) || 0,
            CurrentLuck: parseInt(document.querySelector('input[data-attr="CurrentLuck"]')?.value) || 0
        };
        // Add each attribute with its half and fifth values
        Object.entries(attributes).forEach(([key, value]) => {
            flatState[key] = String(value);
            flatState[`${key}_half`] = String(Math.floor(value / 2));
            flatState[`${key}_fifth`] = String(Math.floor(value / 5));
        });

        // Add skills with their half and fifth values
        Object.entries(baseCharacter.Skill || {}).forEach(([key, _]) => {
            const input = document.querySelector(`input[data-skill="${key}"]`);
            const value = parseInt(input?.value || 0);

            flatState[`Skill_${key}`] = String(value);
            flatState[`Skill_${key}_half`] = String(Math.floor(value / 2));
            flatState[`Skill_${key}_fifth`] = String(Math.floor(value / 5));
        });

        // Special case for Dodge_Copy
        if (baseCharacter.Skill?.Dodge) {
            const dodgeValue = parseInt(baseCharacter.Skill.Dodge) || 0;
            flatState['Dodge_Copy'] = String(dodgeValue);
            flatState['Dodge_Copy_half'] = String(Math.floor(dodgeValue / 2));
            flatState['Dodge_Copy_fifth'] = String(Math.floor(dodgeValue / 5));
        }

        // Add Pulp Talents as comma-separated string
        if (baseCharacter["Pulp-Talents"]?.length > 0) {
            flatState["Pulp Talents"] = baseCharacter["Pulp-Talents"]
                .map(talent => talent.name)
                .join(", ") + ", ";

            flatState["Pulp Talents Descriptions"] = baseCharacter["Pulp-Talents"]
                .map(talent => talent.description)
                .join("~ ") + "~ ";
        }

        return flatState;
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

    recalculateValues(input, type) {
        let value = 0
        if (type === 'attribute') {
            value = parseInt(input.textContent) || 0
        } else {
            value = parseInt(input.value) || 0;
        }
        const container = input.parentElement;
        const halfSpan = container.querySelector('[data-half]');
        const fifthSpan = container.querySelector('[data-fifth]');

        if (halfSpan) halfSpan.textContent = Math.floor(value / 2);
        if (fifthSpan) fifthSpan.textContent = Math.floor(value / 5);

        const characterData = this.getCurrentCharacter();
        if (!characterData) return;

        if (type === 'attribute') {
            const attrName = input.dataset.attribute;
            if (characterData.attributes?.[attrName]) {
                characterData.attributes[attrName].value = value;
            }
        } else if (type === 'skill') {
            const skillName = input.dataset.skill;
            if (characterData.Skill?.[skillName]) {
                characterData.Skill[skillName].value = value;
            }
        }

        document.getElementById('currentCharacter').value = JSON.stringify(characterData);
    },

    updatePersonalInfo(input) {
        const field = input.dataset.field;
        const value = field === 'age' ? parseInt(input.value) || 0 : input.value;
        const characterData = this.getCurrentCharacter();
        if (!characterData) return;

        switch(field) {
            case 'name':
                characterData.Investigators_Name = value;
                break;
            case 'age':
                characterData.Age = value;
                break;
            case 'residence':
                characterData.Residence = value;
                break;
            case 'birthplace':
                characterData.Birthplace = value;
                break;
        }

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


};

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