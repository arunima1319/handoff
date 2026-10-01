document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await login();
});

document.getElementById('signup-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await signup();
})


let accessToken;

async function signup() {

    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;
    const display_name = document.getElementById('display-name').value;

    try {
        const res = await fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password, display_name })
        });

        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to sign up: ${data.error}`)
        } else {
            accessToken = data.access_token;
            let userName = data.display_name;
            document.getElementById('auth-section').style.display = 'none';
            document.getElementById('user-page').style.display = 'block';
            document.getElementById('welcome-message').textContent = `Welcome ${userName}`;
        }
    } catch (error) {
        alert(`Error: ${error.message}`)
    }


}

async function login() {

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });

        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to Login: ${data.error}`)
        }
        console.log(data);
        if (data.access_token) {
            console.log("something happened")
            accessToken = data.access_token;
            let userName = data.user_details.display_name;
            document.getElementById('auth-section').style.display = 'none';
            document.getElementById('user-page').style.display = 'block';
            document.getElementById('welcome-message').textContent = `Welcome ${userName}`;

        }
    } catch (error) {
        alert(`Error : ${error.message}`)
    }

}