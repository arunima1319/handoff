
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
    let username = user.display_name;
    document.getElementById('welcome-message').textContent = `Welcome ${username}`;


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
            const domain = document.createElement('li');
            domain.textContent = dataDomain[i].name;
            domainList.appendChild(domain);

            const button = document.createElement('button')
            domain.appendChild(button);
            button.textContent = 'Enter domain';
            button.addEventListener('click', async () => {
                try {
                    await renderDomainPage(dataDomain[i], user);
                } catch (error) {
                    alert(`Error: ${error.message}`)
                }

            });


        };
    };
}

async function renderDomainPage(domain, user) {

    document.getElementById('user-page').style.display = 'none';
    document.getElementById('domain-page').style.display = 'block';
    document.getElementById('domain-heading').textContent = domain.name;


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
        const domainTaskList = document.createElement('ul');
        domainTaskList.id = 'domain-task-list'

        document.getElementById('domain-tasks').appendChild(domainTaskList);
        for (let i = 0; i < data.length; i++) {
            const task = document.createElement('li');
            task.textContent = data[i].description;

            domainTaskList.appendChild(task);
        };
    };

    if (domain.owner === user.id) {
        const buttonCreateTask = document.createElement('button');
        buttonCreateTask.id = 'create-task-button'
        const domainPage = document.getElementById('domain-page');

        domainPage.appendChild(buttonCreateTask);
        buttonCreateTask.textContent = 'Create New Task'
        buttonCreateTask.addEventListener('click', async () => {

            try {
                await showTaskForm(domain, user);
            } catch (error) {
                alert(`Error: ${error.message}`)
            }

        });

    };
}

async function showTaskForm(domain, user) {

    document.getElementById('create-task-button').style.display = 'none';

    const domainPage = document.getElementById('domain-Page');

    const taskForm = document.getElementById('task-form');
    taskForm.style.display = 'block';

    const selectAssignee = document.getElementById('assignee');


    assignees = await getDomainUsers(domain, user);

    for (let i = 0; i < assignees.length; i++) {
        const option = document.createElement('option');
        selectAssignee.appendChild(option);
        option.value = `${assignees[i].id}`;
        option.textContent = `${assignees[i].display_name}`;
    }


    taskForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const assignee_id = selectAssignee.value;
        const description = document.getElementById('description').value;
        document.getElementById('create-task-button').style.display = "block";
        document.getElementById('task-form').style.display = 'none';
        try {
            await createTask(domain, user, assignee_id, description);
        } catch (error) {
            alert(`Error: ${error.message}`)
        }

    });

};


async function createTask(domain, user, assignee_id, description) {

    const res = await fetch(`/api/domains/${domain.id}/tasks`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${user.access_token}`,
        },
        body: JSON.stringify({ description, assignee_id }),
    });

    const task = await res.json();

    if (!res.ok) {
        if (res.status === 401) {
            alert("Your session has expired, please reload and login again.");
            return;
        }
        throw new Error(`Failed to create task: ${task.error}`)
    } else {
        const newTask = document.createElement('li');
        newTask.textContent = `${task.description}`;
        newTask.id = `${task.id}`;

        document.getElementById('domain-task-list').appendChild(newTask);

    };
};

async function getDomainUsers(domain, user) {


    const res = await fetch(`/api/domains/${domain.id}/users`, {
        method: 'GET',
    });

    const data = await res.json()

    if (!res.ok) {
        throw new Error(`Could not get users of domain: ${data.error}`);
    } else {
        return data;
    }

}