
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
        const taskList = document.createElement('ul');
        document.getElementById('actionable-tasks').appendChild(taskList);
        for (let i = 0; i < data.length; i++) {



            const task = document.createElement('li');
            task.textContent = `${data[i].description}   ---> in ${data[i].domain_name}`;
            taskList.appendChild(task);
            task.addEventListener('click', async () => {

                try {
                    const res = await fetch(`/api/tasks/${data[i].id}`, {
                        method: 'PATCH',
                        headers: {
                            'Authorization': `Bearer ${user.access_token}`,
                        },
                        body: JSON.stringify({ completed: true }),
                    });

                    const checkedTask = await res.json();

                    if (!res.ok) {
                        throw new Error(`Could not check off task: ${checkedTask.error}`);
                    } else {
                        task.style.textDecoration = 'line-through 3px';
                    }
                } catch (error) {
                    alert(`Error: ${error.message}`)
                }


            });
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
        domainsSection = document.getElementById('domains')
        domainsSection.appendChild(domainList);

        const createDomainButton = document.createElement('button');
        createDomainButton.id = 'create-new-domain'
        createDomainButton.textContent = 'Create new Domain';
        domainsSection.appendChild(createDomainButton);

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
                };

            });

        };

        createDomainButton.addEventListener('click', async () => {

            try {
                createDomainButton.style.display = 'none';
                const domainName = document.createElement('input')
                const create = document.createElement('button');
                create.id = 'create-domain';
                create.textContent = "Create";
                domainsSection.appendChild(domainName)
                domainsSection.appendChild(create)

                create.addEventListener('click', async () => {

                    try {
                        createDomainButton.style.display = 'block';
                        domainName.style.display = 'none';
                        create.style.display = 'none';
                        const name = domainName.value;
                        const res = await fetch('/api/domains', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                                'Authorization': `Bearer ${user.access_token}`,
                            },
                            body: JSON.stringify({ name }),
                        });


                        const data = await res.json();
                        if (!res.ok) {
                            throw new Error(`Failed to create new domain: ${data.error}`);
                        }

                        const domainBullet = document.createElement('li');
                        domainBullet.textContent = name;
                        domainList.appendChild(domainBullet)

                        const button = document.createElement('button')
                        domainBullet.appendChild(button);
                        button.textContent = 'Enter domain';
                        button.addEventListener('click', async () => {
                            try {
                                await renderDomainPage(data, user);
                            } catch (error) {
                                alert(`Error: ${error.message}`)
                            };

                        });


                    } catch (error) {
                        alert(`Error: ${error.message}`)
                    }

                });

            } catch (error) {
                alert(`Error: ${error.message}`)
            };

        });

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
        if (domain.owner === user.id) {
            await renderTasksWithDependency(data, user);
        } else {
            renderTasks(data, user);
        }
    };

    if (domain.owner === user.id) {

        const domainPage = document.getElementById('domain-page');

        const buttonCreateTask = document.createElement('button');
        buttonCreateTask.id = 'create-task-button';
        buttonCreateTask.textContent = 'Create New Task'

        domainPage.appendChild(buttonCreateTask);

        buttonCreateTask.addEventListener('click', async () => {

            try {
                await showTaskForm(domain, user, data);
            } catch (error) {
                alert(`Error: ${error.message}`)
            }

        });

        const buttonAddUser = document.createElement('button');
        buttonAddUser.id = 'add-user-button';
        buttonAddUser.textContent = 'Add User to Domain'

        domainPage.appendChild(buttonAddUser);

        buttonAddUser.addEventListener('click', async () => {


            try {
                await showUserForm(domain, user)
            } catch (error) {
                alert(`Error: ${error.message}`)
            }

        });

    };
}

function showTaskDetails(taskBullet, task, user) {

    const taskCard = document.getElementById('task-card');
    const taskAssignee = document.getElementById('task-assignee');
    const dependencyList = document.getElementById('task-dependencies');

    taskBullet.addEventListener('mousedown', async (event) => {

        const mouseX = event.pageX;
        const mouseY = event.pageY;

        taskCard.style.display = '';

        taskCard.style.left = `${mouseX + 10}px`;
        taskCard.style.top = `${mouseY + 10}px`;

        const userRes = await fetch(`/api/users/${task.assignee_id}`, {
            method: 'GET',
        });

        const userData = await userRes.json();
        if (!(userRes.ok)) {
            throw new Error(`Could not get user: ${userData.error}`);
        }

        taskAssignee.textContent = `The task is assigned to ${userData.display_name}`;


        const res = await fetch(`/api/tasks/${task.id}/dependencies`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer: ${user.access_token}`,
            },
        });

        const dependencyTasks = await res.json();

        if (!(res.ok)) {
            throw new Error(`Could not get dependency tasks: ${dependencyTasks.error}`);
        }
        for (let k = 0; k < dependencyTasks.length; k++) {
            dependencyList.textContent = 'First complete these tasks:';
            const bullet = document.createElement('li');
            dependencyList.appendChild(bullet);
            bullet.textContent = dependencyTasks[k].description;

        }

    });

    taskBullet.addEventListener('mouseup', () => {
        dependencyList.innerHTML = '';
        taskAssignee.textContent = '';
        taskCard.style.display = 'none';
    });



}

function renderTasks(tasks, user) {
    const domainTaskList = document.createElement('ul');
    domainTaskList.id = 'domain-task-list'

    document.getElementById('domain-tasks').appendChild(domainTaskList);
    for (let i = 0; i < tasks.length; i++) {
        const task = document.createElement('li');
        task.textContent = tasks[i].description;

        domainTaskList.appendChild(task);
        showTaskDetails(task, tasks[i], user);
    }
}

async function renderTasksWithDependency(tasks, user) {

    const domainTaskList = document.createElement('ul');
    domainTaskList.id = 'domain-task-list'

    document.getElementById('domain-tasks').appendChild(domainTaskList);
    for (let i = 0; i < tasks.length; i++) {
        const task = document.createElement('li');
        task.textContent = tasks[i].description;

        domainTaskList.appendChild(task);
        showTaskDetails(task, tasks[i], user);

        if (!tasks[i].completed_at.Valid) {

            await showAddDependencyButton(task, tasks[i], tasks, user);

        };
    };

};

async function showAddDependencyButton(taskBullet, presentTask, allTasks, user) {

    const dependency = document.createElement('button');
    dependency.textContent = 'Add Dependency';
    taskBullet.appendChild(dependency);

    dependency.addEventListener('click', async () => {

        try {
            dependency.style.display = 'none';
            const dependencyForm = document.createElement('form');
            taskBullet.appendChild(dependencyForm);
            const select = document.createElement('select');
            select.id = 'dependency';
            const label = document.createElement('label');
            label.for = 'dependency';
            label.textContent = 'Dependency';
            dependencyForm.appendChild(label);
            dependencyForm.appendChild(select);

            for (let j = 0; j < allTasks.length; j++) {
                if (!allTasks[j].completed_at.Valid && presentTask.id != allTasks[j].id) {
                    const option = document.createElement('option');
                    option.value = allTasks[j].id;
                    option.textContent = allTasks[j].description;
                    select.appendChild(option)
                }

            };

            const submitDependencyForm = document.createElement('button');
            submitDependencyForm.textContent = 'Add';
            dependencyForm.appendChild(submitDependencyForm);

            dependencyForm.addEventListener('submit', async (event) => {
                try {
                    event.preventDefault();
                    const dependencyID = select.value;
                    dependency.style.display = '';
                    dependencyForm.style.display = 'none';

                    await addDependencyRequest(presentTask.id, dependencyID, user);
                } catch (error) {
                    alert(`Error: ${error.message}`)
                }


            });
        } catch (error) {
            alert(`Error: ${error.message}`)
        }

    });

}

async function addDependencyRequest(taskID, dependencyID, user) {

    const res = await fetch(`/api/tasks/${taskID}/dependencies`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${user.access_token}`,
        },
        body: JSON.stringify({ dependency_id: dependencyID }),
    });

    const data = await res.json();

    if (!res.ok) {
        throw new Error(`Could not add dependency: ${data.error}`)
    } else {
        alert(`${data.message}`)
    }

}

async function showUserForm(domain, user) {

    const userForm = document.createElement('form');
    const memberLabel = document.createElement('label');
    memberLabel.for = 'member';
    memberLabel.textContent = 'User Email';
    memberLabel.placeholder = 'xyz@google.com';

    const memberInput = document.createElement('input');
    memberInput.id = 'member';
    memberInput.required = true;

    const submitAddUser = document.createElement('button')
    submitAddUser.textContent = 'Add';

    userForm.appendChild(memberLabel);
    userForm.appendChild(memberInput);
    userForm.appendChild(submitAddUser);

    const buttonAddUser = document.getElementById('add-user-button');
    buttonAddUser.style.display = 'none';

    const domainPage = document.getElementById('domain-page');
    domainPage.appendChild(userForm);



    userForm.addEventListener('submit', async (event) => {

        try {
            event.preventDefault();
            const email = document.getElementById('member').value;
            buttonAddUser.style.display = '';
            userForm.remove();



            const res = await fetch(`api/domains/${domain.id}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${user.access_token}`,
                },
                body: JSON.stringify({ email }),
            });

            const data = await res.json();
            if (!res.ok) {
                throw new Error(`Could not add user to domain: ${data.error}`)
            } else {
                alert(`${data.message}`)
            }
        } catch (error) {
            alert(`${error.message}`)
        }

    });
}

async function showTaskForm(domain, user, tasks) {


    const taskForm = document.createElement('form');
    const descriptionLabel = document.createElement('label');
    descriptionLabel.for = 'description';
    descriptionLabel.textContent = 'Description';
    descriptionLabel.placeholder = 'Call Sangita';

    const descriptionInput = document.createElement('input');
    descriptionInput.id = 'description';
    descriptionInput.required = true;

    const assigneeLabel = document.createElement('assignee');
    assigneeLabel.for = 'assignee';
    assigneeLabel.textContent = 'Assignee';

    const selectAssignee = document.createElement('select');
    selectAssignee.id = 'assignee';

    const submitCreatedTask = document.createElement('button')
    submitCreatedTask.textContent = 'Create';

    taskForm.appendChild(descriptionLabel);
    taskForm.appendChild(descriptionInput);
    taskForm.appendChild(assigneeLabel);
    taskForm.appendChild(selectAssignee);
    taskForm.appendChild(submitCreatedTask)


    document.getElementById('create-task-button').style.display = 'none';

    const domainPage = document.getElementById('domain-page');
    domainPage.appendChild(taskForm)

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
        document.getElementById('create-task-button').style.display = '';
        taskForm.remove();

        try {
            await createTask(domain, user, tasks, assignee_id, description);
        } catch (error) {
            alert(`Error: ${error.message}`)
        }

    });

};


async function createTask(domain, user, tasks, assignee_id, description) {

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
        showTaskDetails(newTask, task, user);
        showAddDependencyButton(newTask, task, tasks, user);

        document.getElementById('domain-task-list').appendChild(newTask);

    };
};

async function getDomainUsers(domain, user) {


    const res = await fetch(`/api/domains/${domain.id}/users`, {
        method: 'GET',
    });

    const data = await res.json();

    if (!res.ok) {
        throw new Error(`Could not get users of domain: ${data.error}`);
    } else {
        return data;
    };

}