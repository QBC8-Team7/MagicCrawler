-- Insert a new ad
-- name: CreateAd :one
INSERT INTO ad (publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author,
                url, title, description, city, neighborhood, house_type, meterage, rooms_count, year,
                floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng)
VALUES (sqlc.arg('publisher_ad_key'), sqlc.arg('publisher_id'), NOW(), NOW(), sqlc.narg('published_at'),
        sqlc.narg('category'), sqlc.narg('author'),
        sqlc.narg('url'), sqlc.narg('title'), sqlc.narg('description'), sqlc.narg('city'), sqlc.narg('neighborhood'),
        sqlc.narg('house_type'), sqlc.narg('meterage'), sqlc.narg('rooms_count'), sqlc.narg('year'),
        sqlc.narg('floor'), sqlc.narg('total_floors'), sqlc.narg('has_warehouse'), sqlc.narg('has_elevator'),
        sqlc.narg('has_parking'), sqlc.narg('lat'), sqlc.narg('lng'))
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
    has_parking      = COALESCE(sqlc.narg('has_parking'), has_parking),
    lat              = COALESCE(sqlc.narg('lat'), lat),
    lng              = COALESCE(sqlc.narg('lng'), lng)
WHERE publisher_ad_key = sqlc.narg('publisher_ad_key')
RETURNING *;


-- Delete an ad by publisher_ad_key
-- name: DeleteAd :exec
DELETE
FROM ad
WHERE id = sqlc.narg('id');

-- Get all ads with dynamic ordering, limit, and offset
-- name: GetAllAds :many
SELECT *
FROM ad
ORDER BY id DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Get ads based on list of IDs
-- name: GetAdsByIds :many
SELECT *
FROM ad
WHERE id = ANY (sqlc.narg('ad_ids')::bigint[])
ORDER BY created_at DESC;

-- Comprehensive ad search with all attribute filters, including ranges and additional fields
-- name: FilterAds :many
SELECT *
FROM ad
WHERE (publisher_id = coalesce(sqlc.narg('publisher_id'), publisher_id))
  AND (published_at BETWEEN coalesce(sqlc.narg('min_published_at'), published_at) AND coalesce(sqlc.narg('max_published_at'), published_at))
  AND (category::TEXT = coalesce(sqlc.narg('category')::TEXT, category::TEXT))
  AND ((sqlc.narg('author')::varchar IS NULL AND author IS NULL) OR author LIKE coalesce(sqlc.narg('author'), author))
  AND (city = coalesce(sqlc.narg('city'), city))
  AND (neighborhood = coalesce(sqlc.narg('neighborhood'), neighborhood))
  AND (house_type::TEXT = coalesce(sqlc.narg('house_type')::TEXT, house_type::TEXT))
  AND ((meterage IS NULL) OR (meterage BETWEEN coalesce(sqlc.narg('min_meterage'), meterage) AND coalesce(sqlc.narg('max_meterage'), meterage)))
  AND ((rooms_count IS NULL) OR (rooms_count BETWEEN coalesce(sqlc.narg('min_rooms'), rooms_count) AND coalesce(sqlc.narg('max_rooms'), rooms_count)))
  AND ((year IS NULL) OR (year BETWEEN coalesce(sqlc.narg('min_year'), year) AND coalesce(sqlc.narg('max_year'), year)))
  AND ((floor IS NULL) OR (floor BETWEEN coalesce(sqlc.narg('min_floor'), floor) AND coalesce(sqlc.narg('max_floor'), floor)))
  AND ((total_floors IS NULL) OR (total_floors BETWEEN coalesce(sqlc.narg('min_total_floors'), total_floors) AND coalesce(sqlc.narg('max_total_floors'), total_floors)))
  AND (has_warehouse is null or has_warehouse = coalesce(sqlc.narg('has_warehouse'), has_warehouse))
  AND (has_elevator is null or has_elevator = coalesce(sqlc.narg('has_elevator'), has_elevator))
  AND (has_parking is null or has_parking = coalesce(sqlc.narg('has_parking'), has_parking))
  AND ((sqlc.narg('lat')::float IS NOT NULL AND
        sqlc.narg('lng')::float IS NOT NULL AND
        sqlc.narg('radius')::int IS NOT NULL AND
        6371 * ACOS(
                COS(RADIANS(sqlc.narg('lat'))) *
                COS(RADIANS(lat)) *
                COS(RADIANS(lng) - RADIANS(sqlc.narg('lng'))) +
                SIN(RADIANS(sqlc.narg('lat'))) *
                SIN(RADIANS(lat))
               ) <= sqlc.narg('radius'))
    OR (sqlc.narg('lat') IS NULL OR sqlc.narg('lng') IS NULL OR sqlc.narg('radius') IS NULL))
ORDER BY created_at DESC
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Get ads by publisher ID
-- name: GetAdsByPublisher :many
SELECT *
FROM ad
WHERE publisher_id = sqlc.narg('publisher_id')
ORDER BY created_at DESC;


-- Get PublisherAdKey for one specific ad
-- name: GetAdsPublisherByAdKey :one
SELECT ad.publisher_ad_key
FROM ad
WHERE ad.id = sqlc.narg('ad_key');

-- name: GetAdByPublisherAdKey :one
SELECT ad.id
FROM ad
WHERE ad.publisher_ad_key = sqlc.arg('ad_key');

-- Get Ad by its ID
-- name: GetAdByID :one
SELECT
    ad.*,
    CASE
        WHEN fa.ad_id IS NOT NULL THEN true
        ELSE false
        END AS favorite_status
FROM ad
         LEFT JOIN public.favorite_ads fa
                   ON ad.id = fa.ad_id AND fa.user_id = sqlc.arg('user_id')
WHERE ad.id = sqlc.arg('id');


-- name: CountAds :one
SELECT COUNT(*) FROM ad;