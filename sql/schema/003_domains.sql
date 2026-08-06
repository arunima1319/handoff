-- +goose Up

ALTER TABLE domains 
DROP owner;

ALTER TABLE domains
ADD owner UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT; 


-- +goose Down

ALTER TABLE domains
DROP owner; 

ALTER TABLE domains
ADD owner UUID REFERENCES users(id) ON DELETE RESTRICT; 



