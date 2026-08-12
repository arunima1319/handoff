-- +goose Up

ALTER TABLE tasks
ADD CONSTRAINT tasks_domain_assignee_fkey
FOREIGN KEY (domain_id, assignee_id) 
REFERENCES domains_users(domain_id, user_id)
ON DELETE RESTRICT; 

-- +goose Down

ALTER TABLE tasks 
DROP CONSTRAINT tasks_domain_assignee_fkey; 