-- Insert a new price entry for a specific ad
-- name: CreatePrice :one
INSERT INTO price (ad_id, fetched_at, has_price, total_price, price_per_meter, mortgage, normal_price, weekend_price)
VALUES (sqlc.arg('ad_id'),
        NOW(),
        sqlc.narg('has_price'),
        sqlc.narg('total_price'),
        sqlc.narg('price_per_meter'),
        sqlc.narg('mortgage'),
        sqlc.narg('normal_price'),
        sqlc.narg('weekend_price'))
RETURNING *;
