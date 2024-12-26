function handleLogin() {
    let usernameOrEmail = document.getElementById('usernameOrEmail').value;
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
            }else {
                console.log('Success:', data);
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