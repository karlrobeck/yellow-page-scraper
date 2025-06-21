-- name: CreateCategory :one
insert into categories (name,url,size) values (?,?,?) returning *;

-- name: GetAllCategories :many
select * from categories;

-- name: GetCategory :one
select * from categories where name = ?;

-- name: MarkCategoryAsComplete :one
update categories set is_completed = 1 where id = ? returning *;

-- name: MarkCategoryAsIncomplete :one
update categories set is_completed = 0 where id = ? returning *;

-- name: RemoveCategory :exec
delete from categories where id = ?;