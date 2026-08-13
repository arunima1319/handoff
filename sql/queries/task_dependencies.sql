-- name: CreateTaskDependency :exec

INSERT INTO task_dependencies(task_id, dependency_id)
VALUES(
    $1, 
    $2
); 

