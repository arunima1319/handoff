-- +goose Up

ALTER TABLE domains
ADD name TEXT NOT NULL;

-- +goose Down 

ALTER TABLE domains
DROP name; 