-- +goose Up

ALTER TABLE domains_users
DROP CONSTRAINT domains_users_domain_id_fkey; 

ALTER TABLE domains_users
ADD CONSTRAINT domains_users_domain_id_fkey
FOREIGN KEY (domain_id)
REFERENCES domains(id)
ON DELETE RESTRICT; 

ALTER TABLE domains_users 
DROP CONSTRAINT domains_users_user_id_fkey; 

ALTER TABLE domains_users
ADD CONSTRAINT domains_users_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE; 

-- +goose Down 

ALTER TABLE domains_users
DROP CONSTRAINT domains_users_domain_id_fkey; 

ALTER TABLE domains_users
ADD CONSTRAINT domains_users_domain_id_fkey
FOREIGN KEY (domain_id)
REFERENCES domains(id);

ALTER TABLE domains_users 
DROP CONSTRAINT domains_users_user_id_fkey; 

ALTER TABLE domains_users
ADD CONSTRAINT domains_users_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id);
