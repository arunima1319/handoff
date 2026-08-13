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