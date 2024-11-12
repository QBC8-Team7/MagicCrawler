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


-- Get the latest price for a specific ad by its id
-- name: GetLatestPriceByAdID :one
SELECT price.*
FROM price
WHERE ad_id = sqlc.arg('id')
ORDER BY fetched_at DESC
LIMIT 1;

-- Get all prices for a specific ad by its id
-- name: GetAllPricesByAdID :many
SELECT price.*
FROM price
WHERE ad_id = sqlc.arg('id')
ORDER BY fetched_at;

-- Get ads with their latest total price within the specified range.
-- Handles cases with min_price and max_price individually or together.
-- name: FilterAdsByTotalPriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, total_price
               FROM price
               WHERE total_price >= COALESCE(sqlc.narg('min_price'), total_price)
                 AND total_price <= COALESCE(sqlc.narg('max_price'), total_price)
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Get ads with their latest mortgage price within the specified range.
-- Handles cases with min_price and max_price individually or together.
-- name: FilterAdsByMortgagePriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, mortgage
               FROM price
               WHERE mortgage >= COALESCE(sqlc.narg('min_price'), mortgage)
                 AND mortgage <= COALESCE(sqlc.narg('max_price'), mortgage)
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Filter ads based on list of IDs and price range
-- name: FilterAdsByIdsAndTotalPriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, total_price
               FROM price
               WHERE ad_id = ANY (sqlc.narg('ad_ids')::int[])
                 AND total_price BETWEEN sqlc.narg('min_price') AND sqlc.narg('max_price')
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC;

-- Filter ads based on list of IDs and price range
-- name: FilterAdsByIdsAndMortgagePriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, total_price
               FROM price
               WHERE ad_id = ANY (sqlc.narg('ad_ids')::int[])
                 AND mortgage BETWEEN sqlc.narg('min_price') AND sqlc.narg('max_price')
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC;

-- Get ads without associated price
-- name: GetAdsWithoutPrice :many
SELECT ad.*
FROM ad
         LEFT JOIN price ON price.ad_id = ad.id
WHERE price.id IS NULL;

