document.addEventListener("DOMContentLoaded", async ()  => {
    const params = new URLSearchParams(window.location.search);
    const uuid = params.get('uuid');

    if (!uuid) {
        displayMessage("UUID parameter is missing in the URL.");
        return;
    }

    try {
        const response = await fetch(`/api/v1/auth/verify?uuid=${uuid}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            window.location.href = '/auth/login';
        } else {
            displayMessage("Verification failed. Please try again.");
        }
    } catch (error) {
        console.error("Error during verification:", error);
        displayMessage("An error occurred. Please try again later.");
    }
});

function displayMessage(message) {
    const messageContainer = document.getElementById('message');
    if (messageContainer) {
        messageContainer.textContent = message;
        messageContainer.style.color = 'red';
    } else {
        alert(message);
    }
}
