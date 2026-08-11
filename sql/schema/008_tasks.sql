-- +goose Up

ALTER TABLE tasks
ADD completed_at TIMESTAMP DEFAULT NULL; 

-- +goose Down 

ALTER TABLE tasks
DROP completed_at; 