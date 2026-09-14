-- name: CreateTask :one

INSERT INTO tasks(id, created_at, updated_at, description, domain_id, assignee_id)
VALUES(
    GEN_RANDOM_UUID(),
    NOW(), 
    NOW(), 
    $1, 
    $2,
    $3
)
RETURNING *; 

-- name: GetTaskByID :one

SELECT * FROM tasks 
WHERE id = $1; 

-- name: GetTasksOfDomain :many 

SELECT * FROM tasks 
WHERE domain_id = $1; 

-- name: GetUnblockedUserTasks :many 

WITH incomplete_dependencies AS (
    SELECT task_dependencies.task_id FROM tasks JOIN task_dependencies
    ON tasks.id = task_dependencies.dependency_id
    WHERE tasks.completed_at IS NULL
)
SELECT tasks.* FROM tasks
WHERE tasks.id NOT IN (SELECT task_id FROM incomplete_dependencies)
AND tasks.assignee_id = $1; 


