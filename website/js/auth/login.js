let verificationUUID = ""

function handleLogin() {
    let usernameOrEmail = document.getElementById('usernameOrEmail').value;
    let verifyButton = document.getElementById('verifyButton');
    let password = document.querySelector('input[name="password"]').value;

    let data;
    if (usernameOrEmail.includes('@')) {
        data = {
            email: usernameOrEmail,
            password: password
        };
    } else {
        data = {
            username: usernameOrEmail,
            password: password
        };
    }
    fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(data => {
            if(data.error) {
                showMessage("Login failed: "+data.error)
                if (data.verificationUUID) {
                    verifyButton.style.display = "block";
                    verificationUUID = data.verificationUUID;
                }
            }else {
                console.log('Success:', data);
                window.location = "/"
            }
        })
        .catch(error => {
            showMessage("Login failed: "+error)
        });
}
function showMessage(content) {
    const message = document.getElementById('message')
    message.style.display = 'block';
    message.innerHTML = content;
}
function handleVerify() {
    fetch('/api/v1/auth/sendVerification?verificationUUID=' + verificationUUID, {
        method: 'POST',
    })
        .then(response => {
            if(!response.ok) {
                showMessage("Login failed due to: " + response.statusText);
                return ""
            }
            return response.json();
        })
        .then(data => {
            if(data === "") {
                return
            }
            if(data.error) {
                showMessage("Login failed: "+data.error)
            }else {
                showMessage("Sent email. Please verify it now.")
            }
        })
        .catch(error => {
            showMessage("Login failed: "+error)
        });
}
document.getElementById('passwordField').addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
        handleLogin()
    }
});