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
            STR: parseInt(document.querySelector('span[data-attr="Strength"]')?.textContent) || 0,
            DEX: parseInt(document.querySelector('span[data-attr="Dexterity"]')?.textContent) || 0,
            INT: parseInt(document.querySelector('span[data-attr="Intelligence"]')?.textContent) || 0,
            CON: parseInt(document.querySelector('span[data-attr="Constitution"]')?.textContent) || 0,
            APP: parseInt(document.querySelector('span[data-attr="Appearance"]')?.textContent) || 0,
            POW: parseInt(document.querySelector('span[data-attr="Power"]')?.textContent) || 0,
            SIZ: parseInt(document.querySelector('span[data-attr="Size"]')?.textContent) || 0,
            EDU: parseInt(document.querySelector('span[data-attr="Education"]')?.textContent) || 0,
            CurrentMagic: parseInt(document.querySelector('span[data-attr="MagicPoints"]')?.textContent) || 0,
            CurrentHP: parseInt(document.querySelector('span[data-attr="HitPoints"]')?.textContent) || 0,
            CurrentSanity: parseInt(document.querySelector('span[data-attr="Sanity"]')?.textContent) || 0,
            CurrentLuck: parseInt(document.querySelector('span[data-attr="Luck"]')?.textContent) || 0
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
        }

        return flatState;
    },

    exportJSON() {
        try {
            const currentState = this.getCurrentUIState();
            const blob = new Blob([JSON.stringify(currentState, null, 2)], {
                type: 'application/json'
            });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = currentState["Investigators_Name"] + '.json';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error exporting JSON:', error);
            alert('Failed to export character data.');
        }
    },

    async exportPDF() {
        try {
            const element = document.getElementById('character-sheet');
            if (!element) {
                throw new Error('Character sheet not found');
            }

            const clone = element.cloneNode(true);
            this.preparePDFElement(clone);

            const opt = {
                margin: 0,
                filename: 'character-sheet.pdf',
                image: { type: 'jpeg', quality: 1 },
                html2canvas: {
                    scale: 1,
                    useCORS: true,
                    letterRendering: true,
                    backgroundColor: '#ffffff'
                },
                jsPDF: {
                    unit: 'mm',
                    format: 'a4',
                    orientation: 'portrait'
                }
            };

            await html2pdf().from(clone).set(opt).save();
        } catch (error) {
            console.error('Error generating PDF:', error);
            alert('Failed to generate PDF. Please try again.');
        }
    },

    async exportPDFPrime() {
        try {
            const currentState = this.getCurrentUIState();
            const response = await fetch('/api/export-pdf', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(currentState)
            });

            if (!response.ok) {
                throw new Error('PDF export failed');
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${currentState.Investigators_Name || 'character'}.pdf`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error exporting PDF:', error);
            alert('Failed to export PDF. Please try again.');
        }
    },

    preparePDFElement(element) {
        element.style.width = '800px';
        element.style.backgroundColor = '#ffffff';

        element.querySelectorAll('input').forEach(input => {
            const span = document.createElement('span');
            span.textContent = input.value;
            input.parentNode.replaceChild(span, input);
        });

        return element;
    },

    recalculateValues(input, type) {
        const value = parseInt(input.value) || 0;
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

    async loadJSON(input) {
        if (!input.files?.length) return;

        try {
            const file = input.files[0];
            const text = await file.text();
            const character = JSON.parse(text);

            const hiddenInput = document.getElementById('currentCharacter');
            if (hiddenInput) {
                hiddenInput.value = JSON.stringify(character);
            }

            // Trigger HTMX reload of the character sheet
            htmx.trigger('#character-sheet', 'load');
        } catch (error) {
            console.error('Error loading JSON:', error);
            alert('Failed to load character data.');
        }
    },

    initializeEventListeners() {
        const buttonHandlers = {
            'exportPdfBtn': () => this.exportPDF(),
            'exportPdfPrimeBtn': () => this.exportPDFPrime(),
            'exportJsonBtn': () => this.exportJSON()
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

        const fileInput = document.getElementById('loadJsonInput');
        if (fileInput) {
            fileInput.addEventListener('change', () => this.loadJSON(fileInput));
        }

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
    }
};

// Initialize only once when the page loads
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        characterUtils.initializeEventListeners();
    });
} else {
    characterUtils.initializeEventListeners();
}

// Handle HTMX content updates
document.body.addEventListener('htmx:afterSwap', () => {
    characterUtils.initializeInputHandlers();
});

// Make utils available globally
window.characterUtils = characterUtils;