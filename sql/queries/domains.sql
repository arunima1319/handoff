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

-- name: DeleteAllDomains :exec

DELETE FROM domains; 

-- name: GetDomainByID :one

SELECT * FROM domains 
WHERE id = $1; 

-- name: GetUserDomains :many

WITH domain_ids AS (
    SELECT domain_id FROM domains_users
    WHERE user_id = $1
)
SELECT domains.* FROM domains JOIN domain_ids
ON domains.id = domain_ids.domain_id ; 


