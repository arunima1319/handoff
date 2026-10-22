-- name: CreateTaskDependency :exec

INSERT INTO task_dependencies(task_id, dependency_id)
VALUES(
    $1, 
    $2
); 

-- name: GetTaskDependenciesByTaskID :many 

SELECT * FROM task_dependencies
WHERE task_id = $1; 

-- name: GetDependencyTasksByTaskID :many

SELECT tasks.* FROM tasks JOIN task_dependencies 
ON tasks.id = task_dependencies.dependency_id
WHERE task_dependencies.task_id = $1; 
