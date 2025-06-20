-- name: CreateBusinessInCategory :one
insert into business_in_category (category_id,name,url,page) values (?,?,?,?) returning *;

-- name: GetBusinessInCategory :many
select * from business_in_category where category_id = (select id from categories where categories.name = ?) and page = ?;