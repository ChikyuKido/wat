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
            console.log('Success:', data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
}