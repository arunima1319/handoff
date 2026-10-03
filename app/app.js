
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
            document.getElementById('auth-section').style.display = 'none';
            document.getElementById('user-page').style.display = 'block';
            await renderUserPage(data)
        }
    } catch (error) {
        alert(`Error : ${error.message}`)
    }

}

async function renderUserPage(user) {
    let username = user.user_details.display_name;
    document.getElementById('welcome-message').textContent = `Welcome ${username}`;

    try {
        const res = await fetch('/api/users/me/tasks', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${user.access_token}`,
            },
        });

        const data = await res.json();

        if (!res.ok) {
            throw new Error(`Failed to get tasks: ${data.error}`)
        } else {
            taskList = document.createElement('ul');
            document.getElementById('actionable-tasks').appendChild(taskList);
            for (let i = 0; i < data.length; i++) {
                task = document.createElement('li');
                task.textContent = data[i].description;
                taskList.appendChild(task);
            };
        };

        const resDomain = await fetch('/api/users/me/domains', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${user.access_token}`,
            },
        });

        const dataDomain = await resDomain.json();

        if (!resDomain.ok) {
            throw new Error(`Failed to get domains: ${dataDomain.error}`)
        } else {
            domainList = document.createElement('ul');
            document.getElementById('domains').appendChild(domainList);
            for (let i = 0; i < dataDomain.length; i++) {
                domain = document.createElement('li');
                domain.textContent = dataDomain[i].name;
                domainList.appendChild(domain);
                button = document.createElement('button')
                domain.appendChild(button);
                button.textContent = 'Enter domain';
                button.addEventListener('click', async () => {
                    await renderDomainPage(dataDomain[i], user);
                });


            };
        };


    } catch (error) {
        alert(`Error: ${error.message}`)
    }


}

async function renderDomainPage(domain, user) {

    document.getElementById('user-page').style.display = 'none';
    document.getElementById('domain-page').style.display = 'block';
    document.getElementById('domain-heading').textContent = domain.name;

    try {
        const res = await fetch(`/api/domains/${domain.id}/tasks`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${user.access_token}`
            },
        });

        const data = await res.json();

        if (!res.ok) {
            throw new Error(`Failed to get domain tasks: ${data.error}`);
        } else {
            taskList = document.createElement('ul');
            document.getElementById('domain-tasks').appendChild(taskList);
            for (let i = 0; i < data.length; i++) {
                task = document.createElement('li');
                task.textContent = data[i].description;

                taskList.appendChild(task);
            };
        };
    } catch (error) {
        alert(`Error: ${error.message}`);
    };

} 
