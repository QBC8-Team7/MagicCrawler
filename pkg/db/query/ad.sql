-- Insert a new ad
-- name: CreateAd :one
INSERT INTO ad (publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author,
                url, title, description, city, neighborhood, house_type, meterage, rooms_count, year,
                floor, total_floors, has_warehouse, has_elevator, lat, lng)
VALUES (sqlc.arg('publisher_ad_key'), sqlc.arg('publisher_id'), NOW(), NOW(), sqlc.narg('published_at'),
        sqlc.narg('category'), sqlc.narg('author'),
        sqlc.narg('url'), sqlc.narg('title'), sqlc.narg('description'), sqlc.narg('city'), sqlc.narg('neighborhood'),
        sqlc.narg('house_type'), sqlc.narg('meterage'), sqlc.narg('rooms_count'), sqlc.narg('year'),
        sqlc.narg('floor'), sqlc.narg('total_floors'), sqlc.narg('has_warehouse'), sqlc.narg('has_elevator'),
        sqlc.narg('lat'), sqlc.narg('lng'))
RETURNING *;

-- Update an existing ad's details with optional fields
-- name: UpdateAd :one
UPDATE ad
SET publisher_ad_key = COALESCE(sqlc.narg('publisher_ad_key'), publisher_ad_key),
    publisher_id     = COALESCE(sqlc.narg('publisher_id'), publisher_id),
    updated_at       = NOW(),
    published_at     = COALESCE(sqlc.narg('published_at'), published_at),
    category         = COALESCE(sqlc.narg('category'), category),
    author           = COALESCE(sqlc.narg('author'), author),
    url              = COALESCE(sqlc.narg('url'), url),
    title            = COALESCE(sqlc.narg('title'), title),
    description      = COALESCE(sqlc.narg('description'), description),
    city             = COALESCE(sqlc.narg('city'), city),
    neighborhood     = COALESCE(sqlc.narg('neighborhood'), neighborhood),
    house_type       = COALESCE(sqlc.narg('house_type'), house_type),
    meterage         = COALESCE(sqlc.narg('meterage'), meterage),
    rooms_count      = COALESCE(sqlc.narg('rooms_count'), rooms_count),
    year             = COALESCE(sqlc.narg('year'), year),
    floor            = COALESCE(sqlc.narg('floor'), floor),
    total_floors     = COALESCE(sqlc.narg('total_floors'), total_floors),
    has_warehouse    = COALESCE(sqlc.narg('has_warehouse'), has_warehouse),
    has_elevator     = COALESCE(sqlc.narg('has_elevator'), has_elevator),
    lat              = COALESCE(sqlc.narg('lat'), lat),
    lng              = COALESCE(sqlc.narg('lng'), lng)
WHERE publisher_ad_key = sqlc.narg('publisher_ad_key')
RETURNING *;

-- Delete an ad by publisher_ad_key
-- name: DeleteAd :exec
DELETE
FROM ad
WHERE publisher_ad_key = sqlc.narg('publisher_ad_key');

-- Get all ads with dynamic ordering, limit, and offset
-- name: GetAllAds :many
SELECT *
FROM ad
ORDER BY CASE
             WHEN sqlc.narg('order_by') = 'published_at' THEN published_at
             WHEN sqlc.narg('order_by') = 'updated_at' THEN updated_at
             WHEN sqlc.narg('order_by') = 'created_at' THEN created_at
             WHEN sqlc.narg('order_by') = 'year' THEN year
             ELSE id -- Default to ordering by id if no valid order_by is provided
             END
        DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Comprehensive ad search with all attribute filters, including ranges and additional fields
-- name: FilterAds :many
SELECT *
FROM ad
WHERE (publisher_id = coalesce(sqlc.narg('publisher_id'), publisher_id))
  AND (updated_at BETWEEN coalesce(sqlc.narg('min_updated_at'), updated_at) AND coalesce(sqlc.narg('max_updated_at'), updated_at))
  AND (published_at BETWEEN coalesce(sqlc.narg('min_published_at'), published_at) AND coalesce(sqlc.narg('max_published_at'), published_at))
  AND (category = coalesce(sqlc.narg('category'), category))
  AND (author = coalesce(sqlc.narg('author'), author))
  AND (city = coalesce(sqlc.narg('city'), city))
  AND (neighborhood = coalesce(sqlc.narg('neighborhood'), neighborhood))
  AND (house_type = coalesce(sqlc.narg('house_type'), house_type))
  AND (meterage BETWEEN coalesce(sqlc.narg('min_meterage'), meterage) AND coalesce(sqlc.narg('max_meterage'), meterage))
  AND (rooms_count BETWEEN coalesce(sqlc.narg('min_rooms'), rooms_count) AND coalesce(sqlc.narg('max_rooms'), rooms_count))
  AND (year BETWEEN coalesce(sqlc.narg('min_year'), year) AND coalesce(sqlc.narg('max_year'), year))
  AND (floor BETWEEN coalesce(sqlc.narg('min_floor'), floor) AND coalesce(sqlc.narg('max_floor'), floor))
  AND (total_floors BETWEEN coalesce(sqlc.narg('min_total_floors'), total_floors) AND coalesce(sqlc.narg('max_total_floors'), total_floors))
  AND (has_warehouse = coalesce(sqlc.narg('has_warehouse'), has_warehouse))
  AND (has_elevator = coalesce(sqlc.narg('has_elevator'), has_elevator))
  AND (lat BETWEEN coalesce(sqlc.narg('min_lat'), lat) AND coalesce(sqlc.narg('max_lat'), lat))
  AND (lng BETWEEN coalesce(sqlc.narg('min_lng'), lng) AND coalesce(sqlc.narg('max_lng'), lng))
ORDER BY created_at DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Update ad updated_at timestamp by publisher_ad_key
-- name: UpdateAdUpdatedAt :exec
UPDATE ad
SET updated_at = NOW()
WHERE publisher_ad_key = sqlc.narg('publisher_ad_key');

-- Get ads by publisher ID
-- name: GetAdsByPublisher :many
SELECT *
FROM ad
WHERE publisher_id = sqlc.narg('publisher_id')
ORDER BY created_at DESC;

-- Get ads with their latest price within the specified range.
-- This only includes ads where the latest price falls within the specified range.
-- For each ad, it selects the most recent price (based on fetched_at) in the range.
-- name: FilterAdsByPriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, total_price
               FROM price
               WHERE (total_price BETWEEN sqlc.narg('min_price') AND sqlc.narg('max_price'))
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Filter ads based on list of IDs and price range
-- name: FilterAdsByIdsAndPriceRange :many
SELECT ad.*
FROM ad
         JOIN (SELECT DISTINCT ON (ad_id) ad_id, total_price
               FROM price
               WHERE ad_id = ANY (sqlc.narg('ad_ids')::int[])
                 AND total_price BETWEEN sqlc.narg('min_price') AND sqlc.narg('max_price')
               ORDER BY ad_id, fetched_at DESC) latest_price ON latest_price.ad_id = ad.id
ORDER BY ad.created_at DESC;

-- Get ads without associated price
-- name: GetAdsWithoutPrice :many
SELECT ad.*
FROM ad
         LEFT JOIN price ON price.ad_id = ad.id
WHERE price.id IS NULL;

-- Get PublisherAdKey for one specific ad
-- name: GetAdsPublisherByAdKey :one
SELECT ad.publisher_ad_key
FROM ad
WHERE ad.id = sqlc.narg('ad_key');
