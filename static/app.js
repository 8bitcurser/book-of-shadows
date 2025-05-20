// app.js
document.body.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'character-sheet') {
        document.getElementById('main-content').style.display = 'none';
        document.getElementById('character-sheet').style.display = 'block';
    }
});

// Document ready functions
document.addEventListener('DOMContentLoaded', function() {
    // Create CAPTCHA when the modal is shown
    const reportModal = document.getElementById('reportIssueModal');
    if (reportModal) {
        reportModal.addEventListener('show.bs.modal', function() {
            generateCaptcha();
        });
    }
    
    // Handle issue report submission
    const submitBtn = document.getElementById('submitIssueBtn');
    if (submitBtn) {
        submitBtn.addEventListener('click', function() {
            submitIssueReport();
        });
    }
});

// Generate a random math CAPTCHA
function generateCaptcha() {
    const operations = ['+', '-', '*'];
    const operation = operations[Math.floor(Math.random() * 2)]; // Only use + and - for simplicity
    
    let num1, num2, answer;
    
    if (operation === '+') {
        num1 = Math.floor(Math.random() * 10) + 1;
        num2 = Math.floor(Math.random() * 10) + 1;
        answer = num1 + num2;
    } else if (operation === '-') {
        num1 = Math.floor(Math.random() * 10) + 10; // Ensure first number is larger
        num2 = Math.floor(Math.random() * num1);
        answer = num1 - num2;
    } else {
        num1 = Math.floor(Math.random() * 5) + 1;
        num2 = Math.floor(Math.random() * 5) + 1;
        answer = num1 * num2;
    }
    
    document.getElementById('captchaQuestion').textContent = `What is ${num1} ${operation} ${num2}?`;
    document.getElementById('expectedAnswer').value = answer;
    document.getElementById('captchaAnswer').value = '';
    
    // Reset form message
    const messageEl = document.getElementById('formMessage');
    messageEl.classList.add('d-none');
    messageEl.classList.remove('alert-danger', 'alert-success');
}

function submitIssueReport() {
    const form = document.getElementById('issueReportForm');
    const messageEl = document.getElementById('formMessage');
    
    // Check form validity
    if (!form.checkValidity()) {
        messageEl.textContent = 'Please fill in all required fields.';
        messageEl.classList.remove('d-none', 'alert-success');
        messageEl.classList.add('alert-danger');
        return;
    }
    
    // Verify CAPTCHA
    const userAnswer = document.getElementById('captchaAnswer').value;
    const expectedAnswer = document.getElementById('expectedAnswer').value;
    
    if (userAnswer !== expectedAnswer) {
        messageEl.textContent = 'Incorrect CAPTCHA answer. Please try again.';
        messageEl.classList.remove('d-none', 'alert-success');
        messageEl.classList.add('alert-danger');
        generateCaptcha(); // Generate a new CAPTCHA
        return;
    }
    
    // Collect form data
    const issueType = document.getElementById('issueType').value;
    const description = document.getElementById('issueDescription').value;
    const email = document.getElementById('contactEmail').value || 'No email provided';
    
    // Show loading state
    const submitBtn = document.getElementById('submitIssueBtn');
    const originalBtnText = submitBtn.textContent;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Sending...';
    
    // Prepare data for sending to backend
    const reportData = {
        issueType: issueType,
        description: description,
        email: email,
        timestamp: new Date().toISOString()
    };
    
    // Send to your Go backend endpoint
    fetch('/api/report-issue', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(reportData)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to send report');
        }
        return response.json();
    })
    .then(data => {
        // Show success message
        messageEl.textContent = 'Thank you for your report! We\'ll look into this issue.';
        messageEl.classList.remove('d-none', 'alert-danger');
        messageEl.classList.add('alert-success');
        
        // Reset form fields
        document.getElementById('issueType').value = '';
        document.getElementById('issueDescription').value = '';
        document.getElementById('contactEmail').value = '';
        
        // Auto-close the modal after 2 seconds
        setTimeout(() => {
            const modal = bootstrap.Modal.getInstance(document.getElementById('reportIssueModal'));
            modal.hide();
        }, 2000);
    })
    .catch(error => {
        console.error('Error sending report:', error);
        messageEl.textContent = 'Unable to send your report. Please try again later.';
        messageEl.classList.remove('d-none', 'alert-success');
        messageEl.classList.add('alert-danger');
    })
    .finally(() => {
        // Reset button state
        submitBtn.disabled = false;
        submitBtn.textContent = originalBtnText;
    });
}