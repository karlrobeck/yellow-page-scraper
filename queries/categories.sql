-- name: CreateCategory :one
insert into categories (name,url) values (?,?) returning *;

-- name: GetAllCategories :many
select * from categories;

-- name: MarkCategoryAsComplete :one
update categories set is_completed = 1 where id = ? returning *;

-- name: MarkCategoryAsIncomplete :one
update categories set is_completed = 0 where id = ? returning *;

-- name: RemoveCategory :exec
delete from categories where id = ?;