-- name: CreateBusinessInfo :one
insert into business_info (
  trade_name,
  business_name,
  address,
  phone_number,
  email,
  website,
  social_media,
  canonical_link,
  rating,
  description
) values (
  ?,?,?,?,?,?,?,?,?,?
) returning *;

-- name: GetBusinessInfo :one
select * from business_info where canonical_link = ?;
