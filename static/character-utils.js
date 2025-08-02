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

            const response = await fetch('/api/investigator/PDF/' + key, {
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
            const cookieId = this.getCurrentCharacterId();
            const response = await fetch(`/api/investigator/${cookieId}`, {
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
                input.max = parseInt(input.value) || 0;
            } else {
                // Remove max restriction when points are available
                input.removeAttribute('max');
            }
        });
    },

    // Function to handle archetype selection
    handleArchetypeSelection(selectElement) {
        // Show/hide description
        const descriptionElement = document.getElementById('archetype-description');
        const selectedOption = Array.from(selectElement.options).find(option => option.value === selectElement.value);
        
        if (selectedOption && selectedOption.dataset.description) {
            descriptionElement.textContent = selectedOption.dataset.description;
            descriptionElement.style.display = 'block';
        } else {
            descriptionElement.style.display = 'none';
        }

        // Show/hide occupation container
        const occupationContainer = document.getElementById('occupation-container');
        occupationContainer.style.display = selectElement.value ? 'block' : 'none';

        // Update occupation options based on selected archetype
        if (selectElement.value) {
            this.updateOccupationOptions(selectElement.value);
        }

        // Check if the form is complete
        this.checkFormCompletion();
    },

    // Function to update occupation options based on selected archetype
    async updateOccupationOptions(archetypeName) {
        try {
            // Get the occupation select element
            const occupationSelect = document.getElementById('occupation-select');
            if (!occupationSelect) return;

            // Reset occupation selection
            occupationSelect.value = '';
            const descriptionElement = document.getElementById('occupation-description');
            if (descriptionElement) {
                descriptionElement.style.display = 'none';
            }

            // Fetch updated occupation options from server
            const response = await fetch(`/api/archetype/${encodeURIComponent(archetypeName)}/occupations`);
            if (!response.ok) {
                throw new Error('Failed to fetch occupation options');
            }

            const data = await response.json();
            
            // Clear existing options (keep the first "Select Occupation" option)
            while (occupationSelect.options.length > 1) {
                occupationSelect.removeChild(occupationSelect.lastChild);
            }

            // Add suggested occupations first with star indicator
            if (data.suggested && data.suggested.length > 0) {
                data.suggested.forEach(occupation => {
                    const option = document.createElement('option');
                    option.value = occupation.name;
                    option.textContent = `⭐ ${occupation.name}`;
                    option.setAttribute('data-description', occupation.description);
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
                    option.setAttribute('data-description', occupation.description);
                    occupationSelect.appendChild(option);
                });
            }

        } catch (error) {
            console.error('Error updating occupation options:', error);
            // Fall back to showing all occupations if the API call fails
            this.resetOccupationOptions();
        }
    },

    // Function to reset occupation options to show all occupations
    resetOccupationOptions() {
        // This would be called if the archetype-specific API fails
        // For now, we'll leave the current options as they are
        console.log('Falling back to current occupation options');
    },

    // Function to handle occupation selection
    handleOccupationSelection(selectElement) {
        const descriptionElement = document.getElementById('occupation-description');
        const selectedOption = Array.from(selectElement.options).find(option => option.value === selectElement.value);
        
        if (selectedOption && selectedOption.dataset.description) {
            descriptionElement.textContent = selectedOption.dataset.description;
            descriptionElement.style.display = 'block';
        } else {
            descriptionElement.style.display = 'none';
        }

        this.checkFormCompletion();
    },

    // Function to check if the form is complete
    checkFormCompletion() {
        const nameInput = document.querySelector('input[name="name"]');
        const ageInput = document.querySelector('input[name="age"]');
        const residenceInput = document.querySelector('input[name="residence"]');
        const birthplaceInput = document.querySelector('input[name="birthplace"]');
        const archetypeSelect = document.getElementById('archetype-select');
        const occupationSelect = document.getElementById('occupation-select');
        const nextButton = document.getElementById('next-step-button');

        const isFormComplete = (
            nameInput && nameInput.value.trim() !== '' &&
            ageInput && ageInput.value !== '' &&
            residenceInput && residenceInput.value.trim() !== '' &&
            birthplaceInput && birthplaceInput.value.trim() !== '' &&
            archetypeSelect && archetypeSelect.value !== '' &&
            occupationSelect && occupationSelect.value !== ''
        );
        
        nextButton.disabled = !isFormComplete;
        
        if (isFormComplete) {
            nextButton.style.background = 'linear-gradient(135deg, #6d6875 0%, #b5838d 100%)';
            nextButton.classList.add('pulse-button');
        } else {
            nextButton.style.background = '#e5e5e5';
            nextButton.classList.remove('pulse-button');
        }
        
        return isFormComplete;
    },

    // Add these functions to character-utils.js

    // Roll all attributes at once
    rollAllAttributes() {
        // Get the dice roll button and add animation
        const rollButton = document.querySelector('button[onclick="characterUtils.rollAllAttributes()"]');
        rollButton.classList.add('dice-rolling');
        setTimeout(() => {
            rollButton.classList.remove('dice-rolling');
        }, 500);
        
        // Get all attribute inputs
        const attributeInputs = document.querySelectorAll('.attribute-input');
        
        // Randomize the order to make it more visually interesting
        const shuffledInputs = Array.from(attributeInputs).sort(() => Math.random() - 0.5);
        
        // Sequentially roll each attribute with a small delay
        shuffledInputs.forEach((input, index) => {
            setTimeout(() => {
                this.rollAttribute(input);
                
                // Highlight the container
                const container = input.closest('.attribute-container');
                container.classList.add('highlight');
                setTimeout(() => {
                    container.classList.remove('highlight');
                }, 500);
            }, index * 150); // Stagger the rolls
        });
    },

    // Roll a single attribute
    rollAttribute(input) {
        const formula = input.dataset.formula;

        if (formula === '3d6x5') {
            // Roll 3d6 * 5 with animation
            this.animateDiceRoll(3, 6, 5, input);
        } else if (formula === '2d6p6x5') {
            // Roll (2d6 + 6) * 5 with animation
            this.animateDiceRoll(2, 6, 5, input, 6);
        }
        
        // updateAttributeValue will be called by animateDiceRoll when animation completes
    },

    // Add a new function for dice roll animation
    animateDiceRoll(numDice, sides, multiplier, input, bonus = 0) {
        // Start with a random value
        let currentValue = Math.floor(Math.random() * sides * numDice) * multiplier;
        if (bonus > 0) currentValue += (bonus * multiplier);
        input.value = currentValue;
        
        // Animate through several values
        let iterations = 3;
        const animateRoll = setInterval(() => {
            currentValue = Math.floor(Math.random() * sides * numDice) * multiplier;
            if (bonus > 0) currentValue += (bonus * multiplier);
            input.value = currentValue;
            iterations--;
            
            if (iterations <= 0) {
                clearInterval(animateRoll);
                // Final actual roll
                let result = 0;
                for (let i = 0; i < numDice; i++) {
                    result += Math.floor(Math.random() * sides) + 1;
                }
                if (bonus > 0) result += bonus;
                result *= multiplier;
                
                input.value = result;
                this.updateAttributeValue(input);
            }
        }, 100);
        
        return currentValue;
    },

    // Add function to update half and fifth values
    updateDerivedValues(input) {
        const container = input.closest('.attribute-container');
        const halfSpan = container.querySelector('.attr-half');
        const fifthSpan = container.querySelector('.attr-fifth');
        
        const value = parseInt(input.value) || 0;
        
        if (halfSpan) halfSpan.textContent = Math.floor(value / 2);
        if (fifthSpan) fifthSpan.textContent = Math.floor(value / 5);
    },

    // Initialize attribute form
    initAttributeForm() {
        // Set up input change listeners for all attribute inputs
        document.querySelectorAll('.attribute-input').forEach(input => {
            input.addEventListener('input', function() {
                characterUtils.updateDerivedValues(this);
            });
            
            // Initial calculation for any pre-filled values
            characterUtils.updateDerivedValues(input);
        });

        // Check if attributes are complete on initialization
        this.checkAttributesComplete();
    },

    // Add these to character-utils.js

    // Function to go back to attributes page
    goBackToAttributes(investigatorId) {
        htmx.ajax('GET', '/wizard/attributes/' + investigatorId, {
            target: '#character-sheet',
            swap: 'innerHTML'
        });
    },

    // Function to navigate between tabs
    navigateToTab(tabName) {
        // Enable the tab if it's not the archetype tab (which is always enabled)
        if (tabName !== 'archetype') {
            document.getElementById(tabName + '-tab').disabled = false;
        }

        // Add a slight delay before switching tabs for better animation
        setTimeout(() => {
            // Switch to the tab
            const tab = new bootstrap.Tab(document.getElementById(tabName + '-tab'));
            tab.show();
            
            // Scroll to top of the tab content
            document.getElementById(tabName + '-skills').scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }, 100);
    },

    // Function to increment or decrement skill value using buttons
    adjustSkillValue(btn, increment) {
        const inputGroup = btn.closest('.skill-input-group');
        const input = inputGroup.querySelector('.skill-input');
        const currentValue = parseInt(input.value) || 0;
        
        // Increment or decrement the value
        input.value = currentValue + (increment ? 1 : -1);
        
        // Trigger the change event to update calculations
        input.dispatchEvent(new Event('change'));
    },


    recalculateSkillValues(input) {
        const skillName = input.dataset.skill;
        const value = parseInt(input.value) || 0;
        const prevValue = parseInt(input.dataset.skillvalue) || 0;
        const defaultValue = parseInt(input.dataset.skilldefault) || 0;
        const skillType = input.dataset.skilltype || 'archetype';

        // Apply max limit
        if (value > 90) {
            input.value = 90;
            return;
        }
        
        // Apply min limit (default value)
        if (value < defaultValue) {
            input.value = defaultValue;
            return;
        }

        // Calculate difference
        const difference = value - prevValue;
        
        // Skip if there's no change
        if (difference === 0) return;

        // Get proper points element
        let pointsId = "";
        let confirmId = "";

        if (skillType === "archetype" || document.querySelector('#archetype-skills.active')) {
            pointsId = "archetype-points";
            confirmId = "confirm-archetype-container";
        } else if (skillType === "occupation" || document.querySelector('#occupation-skills.active')) {
            pointsId = "occupation-points";
            confirmId = "confirm-occupation-container";
        } else if (skillType === "general" || document.querySelector('#general-skills.active')) {
            pointsId = "general-points";
            confirmId = "confirm-general-container";
        }

        // Get points element and current remaining points
        const pointsElement = document.getElementById(pointsId);

        if (pointsElement) {
            const currentPoints = parseInt(pointsElement.textContent) || 0;
            const newPoints = currentPoints - difference;

            // Don't allow negative points
            if (newPoints < 0) {
                input.value = prevValue;
                
                // Visual feedback for error
                input.classList.add('is-invalid');
                
                // Remove invalid class after a short delay
                setTimeout(() => {
                    input.classList.remove('is-invalid');
                }, 800);
                
                return;
            }

            // Update points display with visual feedback
            pointsElement.textContent = newPoints;
            
            // Add color coding based on remaining points
            if (newPoints < 10) {
                pointsElement.style.color = '#e74c3c';
            } else {
                pointsElement.style.color = '#b5838d';
            }
            
            // Highlight the skill box that was changed
            const skillBox = input.closest('.skill-box');
            skillBox.classList.add('flash-highlight');
            setTimeout(() => {
                skillBox.classList.remove('flash-highlight');
            }, 800);

            // Update skill value tracking
            input.dataset.skillvalue = value;

            // Update half and fifth values
            const container = input.closest('.skill-values');
            const halfSpan = container.querySelector('[data-half]');
            const fifthSpan = container.querySelector('[data-fifth]');

            if (halfSpan) halfSpan.textContent = Math.floor(value / 2);
            if (fifthSpan) fifthSpan.textContent = Math.floor(value / 5);

            // Show continue button if all points are used
            const confirmContainer = document.getElementById(confirmId);
            if (confirmContainer && newPoints === 0) {
                confirmContainer.style.opacity = "1";
                confirmContainer.style.pointerEvents = "auto";
            }
        }

        // Call server update
        this.updateInvestigator("skills", skillName, value);
    },

    // Initialize skill form
    initSkillForm() {
        // Ensure all tab continue buttons are visible
        document.querySelectorAll('.transition-opacity').forEach(container => {
            container.style.opacity = "1";
            container.style.pointerEvents = "auto";
        });
        
        // Add hover effects to skill boxes
        document.querySelectorAll('.skill-box').forEach(box => {
            box.addEventListener('mouseenter', function() {
                this.style.backgroundColor = '#f0f0f0';
            });
            box.addEventListener('mouseleave', function() {
                this.style.backgroundColor = '#f8f9fa';
            });
        });
        
        // Initialize any skill inputs
        document.querySelectorAll('.skill-input').forEach(input => {
            // Add focus/blur effects
            input.addEventListener('focus', function() {
                this.closest('.skill-box').style.borderColor = '#b5838d';
            });
            
            input.addEventListener('blur', function() {
                this.closest('.skill-box').style.borderColor = 'rgba(0,0,0,0.05)';
            });
        });
    },


    toggleLock(checkbox) {
        const isLocked = checkbox.checked;
        
        // Get all editable elements
        const editables = document.querySelectorAll('.editable');
        
        editables.forEach(element => {
            element.disabled = isLocked;
        });
        
        // Show a message 
        const lockMessage = isLocked ? 'Character sheet is now locked. Unlock to make changes.' : 'Character sheet is now editable.';
        const toast = document.createElement('div');
        toast.className = 'position-fixed bottom-0 end-0 p-3';
        toast.style.zIndex = 1050;
        toast.innerHTML = `
            <div class="toast show" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">${isLocked ? '🔒 Locked' : '🔓 Unlocked'}</strong>
                    <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close" onclick="this.parentElement.parentElement.remove()"></button>
                </div>
                <div class="toast-body">
                    ${lockMessage}
                </div>
            </div>
        `;
        document.body.appendChild(toast);
        
        // Remove toast after 1.5 seconds
        setTimeout(() => {
            toast.remove();
        }, 1500);
},

// Enhanced recalculateValues function for character sheet
    recalculateSheetValues(input, type) {
        let value = parseInt(input.value) || 0;
        
        // Update half and fifth values if applicable
        const container = input.closest('.characteristic-box') || input.parentElement;
        const halfSpan = container.querySelector('[data-half]');
        const fifthSpan = container.querySelector('[data-fifth]');
        
        if (halfSpan) halfSpan.textContent = Math.floor(value / 2);
        if (fifthSpan) fifthSpan.textContent = Math.floor(value / 5);
        
        // Handle attribute updates
        if (type === 'attribute') {
            const attrName = input.dataset.attr;
            
            // Get character data if it exists
            const charData = this.getCurrentCharacter();
            if (charData && charData.attributes && charData.attributes[attrName]) {
                charData.attributes[attrName].value = value;
            }
            
            // Update the server with the new value
            this.updateInvestigator("combat", attrName, value)
                .then(() => {
                    // Add visual feedback
                    input.classList.add('bg-success', 'bg-opacity-10');
                    setTimeout(() => {
                        input.classList.remove('bg-success', 'bg-opacity-10');
                    }, 300);
                })
                .catch(error => {
                    console.error('Error updating attribute:', error);
                    input.classList.add('bg-danger', 'bg-opacity-10');
                    setTimeout(() => {
                        input.classList.remove('bg-danger', 'bg-opacity-10');
                    }, 300);
                });
        } 
        // Handle skill updates
        else if (type === 'skill') {
            // Update skill derived values (half/fifth)
            const skillItem = input.closest('.skill-item');
            if (skillItem) {
                const derivedValues = skillItem.querySelector('.derived-values');
                if (derivedValues) {
                    const spans = derivedValues.querySelectorAll('span');
                    if (spans.length >= 3) {
                        spans[0].textContent = Math.floor(value / 2);
                        spans[2].textContent = Math.floor(value / 5);
                    }
                }
            }
            
            const skillName = input.dataset.skill;
            const prevValue = parseInt(input.dataset.skillvalue) || value;
            
            // Update the server with the new value
            this.updateInvestigator("skills", skillName, value)
                .then(() => {
                    // Add visual feedback
                    input.classList.add('bg-success', 'bg-opacity-10');
                    setTimeout(() => {
                        input.classList.remove('bg-success', 'bg-opacity-10');
                    }, 300);
                })
                .catch(error => {
                    console.error('Error updating skill:', error);
                    input.classList.add('bg-danger', 'bg-opacity-10');
                    setTimeout(() => {
                        input.classList.remove('bg-danger', 'bg-opacity-10');
                    }, 300);
                });
            
            // Update skill value tracking
            input.dataset.skillvalue = value.toString();
        } else {
            const statusName = input.dataset.stat;
            const charData = this.getCurrentCharacter();
            charData[statusName] = value;
            this.updateInvestigator("stats", statusName, value);
        
        }
    },

// Initialize character sheet
    initCharacterSheet() {
        // Add hover effects to stat pills and characteristic boxes
        const hoverElements = document.querySelectorAll('.stat-pill, .characteristic-box, .skill-item');
        hoverElements.forEach(element => {
            element.addEventListener('mouseenter', function() {
                if (!this.closest('.card').classList.contains('locked')) {
                    this.style.backgroundColor = '#f0f0f0';
                }
            });
            element.addEventListener('mouseleave', function() {
                if (!this.closest('.card').classList.contains('locked')) {
                    this.style.backgroundColor = '#f8f9fa';
                }
            });
        });
    
        // Adjust skill names dynamically based on available space
        adjustSkillNames = function() {
            const skillItems = document.querySelectorAll('.skill-item');
            skillItems.forEach(item => {
                const container = item.querySelector('.skill-name-container');
                const nameWrapper = item.querySelector('.skill-name-wrapper');
                const name = item.querySelector('.skill-name');
                
                if (container && nameWrapper && name) {
                    // Calculate available width
                    const containerWidth = container.offsetWidth;
                    
                    // Adjust max-width dynamically if needed
                    if (containerWidth < 120) {
                        name.style.maxWidth = (containerWidth - 30) + 'px';
                    }
                }
            });
        };
    
        // Run initially and on window resize
        adjustSkillNames();
        window.addEventListener('resize', adjustSkillNames);
    },

    // Toggle skill pin status
    togglePinSkill(button) {
        const skillName = button.dataset.skill;
        const isPinned = button.dataset.pinned === 'true';
        const newPinnedStatus = !isPinned;
        
        // Update button appearance
        button.dataset.pinned = newPinnedStatus.toString();
        
        // Update icon based on pinned status
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
        
        // Update skill item parent to show priority styling (optional)
        const skillItem = button.closest('.skill-item');
        if (skillItem) {
            skillItem.dataset.priority = newPinnedStatus.toString();
        }
        
        // Send update to server
        this.updateInvestigator("skill_prio", skillName, newPinnedStatus)
            .then(() => {
                // Add visual feedback
                button.classList.add('scale-up');
                setTimeout(() => {
                    button.classList.remove('scale-up');
                }, 300);
                
                // Reload the skills section to reorder the skills
                htmx.ajax('GET', `/api/investigator/${this.getCurrentCharacterId()}`, {
                    target: '#character-sheet',
                    swap: 'innerHTML'
                });
            })
            .catch(error => {
                console.error('Error updating skill priority status:', error);
                // Revert to previous state on error
                button.dataset.pinned = isPinned.toString();
                if (icon) {
                    if (isPinned) {
                        icon.classList.remove('bi-pin');
                        icon.classList.add('bi-pin-fill');
                        icon.style.color = '#C97700'; // Occult Amber
                    } else {
                        icon.classList.remove('bi-pin-fill');
                        icon.classList.add('bi-pin');
                        icon.style.color = '#B0B0B0'; // Phantom Gray
                    }
                }
                // Revert priority styling
                if (skillItem) {
                    skillItem.dataset.priority = isPinned.toString();
                }
            });
    },

    // Helper to get current character ID
    getCurrentCharacterId() {
        // First try to get from name field (character sheet context)
        const nameInput = document.querySelector('input[data-field="Name"]');
        if (nameInput && nameInput.id) {
            return nameInput.id;
        }
        
        // Then try to get from hidden investigatorId field (wizard context)
        const hiddenInput = document.getElementById('investigatorId');
        if (hiddenInput && hiddenInput.value) {
            return hiddenInput.value;
        }
        
        return '';
    },

    // Function to load personal info step
    loadPersonalInfo(investigatorId) {
        fetch(`/wizard/base/${investigatorId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.text();
            })
            .then(html => {
                document.getElementById('character-sheet').innerHTML = html;
            })
            .catch(error => {
                console.error('Error loading personal info:', error);
            });
    },

    // Function to proceed to skills step
    proceedToSkills(investigatorId) {
        fetch(`/wizard/skills/${investigatorId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.text();
            })
            .then(html => {
                document.getElementById('character-sheet').innerHTML = html;
            })
            .catch(error => {
                console.error('Error loading skills step:', error);
            });
    },

    // Function to complete character and navigate to character sheet
    completeCharacter(investigatorId) {
        fetch(`/api/investigator/${investigatorId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.text();
            })
            .then(html => {
                document.getElementById('character-sheet').innerHTML = html;
            })
            .catch(error => {
                console.error('Error completing character:', error);
            });
    },

    // Function to handle personal info changes when investigator exists
    async handlePersonalInfoChange(input) {
        // Check if investigator exists (has hidden ID field)
        const investigatorId = this.getCurrentCharacterId();
        
        if (investigatorId) {
            // Update existing investigator
            const field = input.dataset.field;
            const value = field === 'Age' ? parseInt(input.value) || 0 : input.value;
            
            try {
                await this.updateInvestigator("personalInfo", field, value);
                
                // Add visual feedback
                input.classList.add('bg-success', 'bg-opacity-10');
                setTimeout(() => {
                    input.classList.remove('bg-success', 'bg-opacity-10');
                }, 300);
                
            } catch (error) {
                console.error('Error updating personal info:', error);
                input.classList.add('bg-danger', 'bg-opacity-10');
                setTimeout(() => {
                    input.classList.remove('bg-danger', 'bg-opacity-10');
                }, 300);
            }
        }
        
        // Always check form completion
        this.checkFormCompletion();
    },

    // Function to update attribute value and check if all attributes are filled
    async updateAttributeValue(input) {
        const attrAbbrev = input.name;
        const value = parseInt(input.value) || 0;

        // Mapping from abbreviation to full name
        const attributeNames = {
            "POW": "Power",
            "STR": "Strength", 
            "LCK": "Luck",
            "APP": "Appearance",
            "DEX": "Dexterity",
            "INT": "Intelligence",
            "EDU": "Education",
            "SIZ": "Size",
            "CON": "Constitution"
        };

        const attrName = attributeNames[attrAbbrev] || attrAbbrev;

        // Update derived values (half/fifth)
        this.updateDerivedValues(input);

        // Update the server with the new value
        try {
            await this.updateInvestigator("attributes", attrName, value);
            
            // Add visual feedback
            input.classList.add('bg-success', 'bg-opacity-10');
            setTimeout(() => {
                input.classList.remove('bg-success', 'bg-opacity-10');
            }, 300);

            // Check if all attributes have values and enable/disable proceed button
            this.checkAttributesComplete();
            
        } catch (error) {
            console.error('Error updating attribute:', error);
            input.classList.add('bg-danger', 'bg-opacity-10');
            setTimeout(() => {
                input.classList.remove('bg-danger', 'bg-opacity-10');
            }, 300);
        }
    },

    // Function to check if all attributes are complete and enable proceed button
    checkAttributesComplete() {
        const attributeInputs = document.querySelectorAll('.attribute-input');
        const proceedButton = document.querySelector('.gradient-button');
        
        if (!proceedButton) return;

        const allFilled = Array.from(attributeInputs).every(input => {
            const value = parseInt(input.value) || 0;
            return value > 0;
        });

        if (allFilled) {
            proceedButton.disabled = false;
            proceedButton.style.background = 'linear-gradient(135deg, #6d6875 0%, #b5838d 100%)';
            proceedButton.classList.add('pulse-button');
        } else {
            proceedButton.disabled = true;
            proceedButton.style.background = '#e5e5e5';
            proceedButton.classList.remove('pulse-button');
        }
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