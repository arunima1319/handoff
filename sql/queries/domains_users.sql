-- name: AddUserToDomain :exec 

INSERT INTO domains_users(domain_id, user_id)
VALUES(
    $1, 
    $2
); 

-- name: DeleteAllUsersFromDomains :exec

DELETE FROM domains_users; 