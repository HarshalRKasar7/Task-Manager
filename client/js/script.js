function showForm(formId) {
    // Hide all forms
    const forms = ['taskForm', 'removeForm', 'completeForm'];
    forms.forEach(function(id) {
        document.getElementById(id).style.display = 'none';
    });

    // Show the selected form
    document.getElementById(formId).style.display = 'block';
}

function hideForm(formId) {
    document.getElementById(formId).style.display = 'none';
}

function redirectTo(url) {
    window.location.href = url;
}