-- name: CreateTaskDependency :exec

INSERT INTO task_dependencies(task_id, dependency_id)
VALUES(
    $1, 
    $2
); 

-- name: GetTaskDependenciesByTask :many 

SELECT * FROM task_dependencies
WHERE task_id = $1; 

