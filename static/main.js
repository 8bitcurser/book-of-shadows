let currentCharacter = null;

function displayCharacter(character) {
    currentCharacter = character;
    const sheet = document.getElementById('character-sheet');

    sheet.innerHTML = `
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="character-info">
                <h2 class="text-2xl font-bold mb-4">${character.Investigators_Name}</h2>
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <p><strong>Occupation:</strong> ${character.Occupation.Name}</p>
                        <p><strong>Age:</strong> ${character.Age}</p>
                        <p><strong>Residence:</strong> ${character.Residence}</p>
                        <p><strong>Birthplace:</strong> ${character.Birthplace}</p>
                    </div>
                    <div>
                        <p><strong>Movement:</strong> ${character.MOV}</p>
                        <p><strong>Build:</strong> ${character.Build}</p>
                        <p><strong>Damage Bonus:</strong> ${character.DamageBonus}</p>
                    </div>
                </div>
            </div>
            
            <div class="attributes">
                <h3 class="text-xl font-bold mb-2">Attributes</h3>
                <div class="grid grid-cols-2 gap-2">
                    ${Object.entries(character.attributes)
        .map(([key, attr]) => `
                            <div class="flex justify-between">
                                <span>${attr.Name}:</span>
                                <span>${attr.Value}</span>
                            </div>
                        `).join('')}
                </div>
            </div>
            
            <div class="skills col-span-full">
                <h3 class="text-xl font-bold mb-2">Skills</h3>
                <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
                    ${Object.entries(character.Skill)
        .map(([key, skill]) => `
                            <div class="flex justify-between">
                                <span>${skill.Name}:</span>
                                <span>${skill.Value}%</span>
                            </div>
                        `).join('')}
                </div>
            </div>
        </div>
    `;
}

async function exportPDF() {
    if (!currentCharacter) {
        alert('Please generate a character first');
        return;
    }

    try {
        const response = await fetch('/api/export-pdf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentCharacter)
        });

        if (!response.ok) throw new Error('Failed to generate PDF');

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'investigator.pdf';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
    } catch (error) {
        console.error('Error exporting PDF:', error);
        alert('Failed to generate PDF');
    }
}

function exportJSON() {
    if (!currentCharacter) {
        alert('Please generate a character first');
        return;
    }

    const blob = new Blob([JSON.stringify(currentCharacter, null, 2)], {
        type: 'application/json'
    });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'investigator.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
}

// Handle the response from the generate endpoint
htmx.on('htmx:afterRequest', function(evt) {
    if (evt.detail.pathInfo.requestPath === '/api/generate') {
        const response = JSON.parse(evt.detail.xhr.responseText);
        if (response.data) {
            displayCharacter(response.data);
        }
    }
});