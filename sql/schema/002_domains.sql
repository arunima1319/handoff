-- +goose Up

CREATE TABLE domains(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    owner UUID REFERENCES users(id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE domains;