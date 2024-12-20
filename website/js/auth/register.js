function handleRegister() {
    let username = document.getElementById('username').value;
    let email = document.getElementById('email').value;
    let password = document.querySelector('input[name="password"]').value;

    fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
            email: email
        })
    })
        .then(response => response.json())
        .then(data => {
            if(data.verification && data.emailSent) {
                showMessage("Sent message to your email. Please verify it")
                return
            }
            if(data.verification && !data.emailSent) {
                showMessage("Failed to send email.")
                return
            }
            if(!data.verification) {
                window.location = "/auth/login"
            }
        })
        .catch(error => {
            showMessage("Failed to send request")
            console.log(error)
        });
}

function showMessage(content) {
    const message = document.getElementById('message')
    message.style.display = 'block';
    message.innerHTML = content;
}