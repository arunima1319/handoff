document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await login();
});

let accessToken;


async function login() {

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const res = await fetch('api/login', {
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