-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id ASC;

-- name: CreateTask :exec
INSERT INTO tasks (
  name
) VALUES (
  ?
);


-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;   

