-- name: CreateDomain :one

INSERT INTO domains(id, created_at, updated_at, owner, name)
VALUES(
    GEN_RANDOM_UUID(), 
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;