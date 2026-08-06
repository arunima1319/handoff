-- +goose Up

CREATE TABLE domains_users(
    domain_id UUID NOT NULL REFERENCES domains(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY(domain_id, user_id)
);

-- +goose Down

DROP TABLE domains_users;