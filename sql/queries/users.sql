-- name: CreateUser :one 

INSERT INTO users(id, created_at, updated_at, email, display_name)
VALUES(
    GEN_RANDOM_UUID(), 
    NOW(), 
    NOW(), 
    $1, 
    $2
)
RETURNING *;

-- name: GetUsersOfDomain :many

SELECT users.id, users.created_at, users.updated_at, users.email, users.display_name 
FROM users JOIN domains_users 
ON users.id = domains_users.user_id
WHERE domain_id = $1; 

