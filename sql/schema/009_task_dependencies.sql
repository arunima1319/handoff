-- +goose Up

CREATE TABLE task_dependencies(
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, 
    dependency_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY(task_id, dependency_id)
); 

-- +goose Down 

DROP TABLE task_dependencies; 