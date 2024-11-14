// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: ad.sql

package sqlc

import (
	"context"
	"time"
)

const createAd = `-- name: CreateAd :one
INSERT INTO ad (publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author,
                url, title, description, city, neighborhood, house_type, meterage, rooms_count, year,
                floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng)
VALUES ($1, $2, NOW(), NOW(), $3,
        $4, $5,
        $6, $7, $8, $9, $10,
        $11, $12, $13, $14,
        $15, $16, $17, $18,
        $19, $20, $21)
RETURNING id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
`

type CreateAdParams struct {
	PublisherAdKey string    `json:"publisher_ad_key"`
	PublisherID    *int32    `json:"publisher_id"`
	PublishedAt    time.Time `json:"published_at"`
	Category       string    `json:"category"`
	Author         *string   `json:"author"`
	Url            *string   `json:"url"`
	Title          *string   `json:"title"`
	Description    *string   `json:"description"`
	City           *string   `json:"city"`
	Neighborhood   *string   `json:"neighborhood"`
	HouseType      string    `json:"house_type"`
	Meterage       *int32    `json:"meterage"`
	RoomsCount     *int32    `json:"rooms_count"`
	Year           *int32    `json:"year"`
	Floor          *int32    `json:"floor"`
	TotalFloors    *int32    `json:"total_floors"`
	HasWarehouse   *bool     `json:"has_warehouse"`
	HasElevator    *bool     `json:"has_elevator"`
	HasParking     *bool     `json:"has_parking"`
	Lat            *float64  `json:"lat"`
	Lng            *float64  `json:"lng"`
}

// Insert a new ad
func (q *Queries) CreateAd(ctx context.Context, arg CreateAdParams) (Ad, error) {
	row := q.db.QueryRow(ctx, createAd,
		arg.PublisherAdKey,
		arg.PublisherID,
		arg.PublishedAt,
		arg.Category,
		arg.Author,
		arg.Url,
		arg.Title,
		arg.Description,
		arg.City,
		arg.Neighborhood,
		arg.HouseType,
		arg.Meterage,
		arg.RoomsCount,
		arg.Year,
		arg.Floor,
		arg.TotalFloors,
		arg.HasWarehouse,
		arg.HasElevator,
		arg.HasParking,
		arg.Lat,
		arg.Lng,
	)
	var i Ad
	err := row.Scan(
		&i.ID,
		&i.PublisherAdKey,
		&i.PublisherID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublishedAt,
		&i.Category,
		&i.Author,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.City,
		&i.Neighborhood,
		&i.HouseType,
		&i.Meterage,
		&i.RoomsCount,
		&i.Year,
		&i.Floor,
		&i.TotalFloors,
		&i.HasWarehouse,
		&i.HasElevator,
		&i.HasParking,
		&i.Lat,
		&i.Lng,
	)
	return i, err
}

const deleteAd = `-- name: DeleteAd :exec
DELETE
FROM ad
WHERE id = $1
`

// Delete an ad by publisher_ad_key
func (q *Queries) DeleteAd(ctx context.Context, id *int64) error {
	_, err := q.db.Exec(ctx, deleteAd, id)
	return err
}

const filterAds = `-- name: FilterAds :many
SELECT id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
FROM ad
WHERE (publisher_id = coalesce($1, publisher_id))
  AND (updated_at BETWEEN coalesce($2, updated_at) AND coalesce($3, updated_at))
  AND (published_at BETWEEN coalesce($4, published_at) AND coalesce($5, published_at))
  AND (category = coalesce($6, category))
  AND (author = coalesce($7, author))
  AND (city = coalesce($8, city))
  AND (neighborhood = coalesce($9, neighborhood))
  AND (house_type = coalesce($10, house_type))
  AND (meterage BETWEEN coalesce($11, meterage) AND coalesce($12, meterage))
  AND (rooms_count BETWEEN coalesce($13, rooms_count) AND coalesce($14, rooms_count))
  AND (year BETWEEN coalesce($15, year) AND coalesce($16, year))
  AND (floor BETWEEN coalesce($17, floor) AND coalesce($18, floor))
  AND (total_floors BETWEEN coalesce($19, total_floors) AND coalesce($20, total_floors))
  AND (has_warehouse = coalesce($21, has_warehouse))
  AND (has_elevator = coalesce($22, has_elevator))
  AND (has_parking = coalesce($23, has_parking))
  AND (lat BETWEEN coalesce($24, lat) AND coalesce($25, lat))
  AND (lng BETWEEN coalesce($26, lng) AND coalesce($27, lng))
ORDER BY created_at DESC
LIMIT $29 OFFSET $28
`

type FilterAdsParams struct {
	PublisherID    *int32    `json:"publisher_id"`
	MinUpdatedAt   time.Time `json:"min_updated_at"`
	MaxUpdatedAt   time.Time `json:"max_updated_at"`
	MinPublishedAt time.Time `json:"min_published_at"`
	MaxPublishedAt time.Time `json:"max_published_at"`
	Category       string    `json:"category"`
	Author         *string   `json:"author"`
	City           *string   `json:"city"`
	Neighborhood   *string   `json:"neighborhood"`
	HouseType      string    `json:"house_type"`
	MinMeterage    *int32    `json:"min_meterage"`
	MaxMeterage    *int32    `json:"max_meterage"`
	MinRooms       *int32    `json:"min_rooms"`
	MaxRooms       *int32    `json:"max_rooms"`
	MinYear        *int32    `json:"min_year"`
	MaxYear        *int32    `json:"max_year"`
	MinFloor       *int32    `json:"min_floor"`
	MaxFloor       *int32    `json:"max_floor"`
	MinTotalFloors *int32    `json:"min_total_floors"`
	MaxTotalFloors *int32    `json:"max_total_floors"`
	HasWarehouse   *bool     `json:"has_warehouse"`
	HasElevator    *bool     `json:"has_elevator"`
	HasParking     *bool     `json:"has_parking"`
	MinLat         *float64  `json:"min_lat"`
	MaxLat         *float64  `json:"max_lat"`
	MinLng         *float64  `json:"min_lng"`
	MaxLng         *float64  `json:"max_lng"`
	Offset         *int32    `json:"offset"`
	Limit          *int32    `json:"limit"`
}

// Comprehensive ad search with all attribute filters, including ranges and additional fields
func (q *Queries) FilterAds(ctx context.Context, arg FilterAdsParams) ([]Ad, error) {
	rows, err := q.db.Query(ctx, filterAds,
		arg.PublisherID,
		arg.MinUpdatedAt,
		arg.MaxUpdatedAt,
		arg.MinPublishedAt,
		arg.MaxPublishedAt,
		arg.Category,
		arg.Author,
		arg.City,
		arg.Neighborhood,
		arg.HouseType,
		arg.MinMeterage,
		arg.MaxMeterage,
		arg.MinRooms,
		arg.MaxRooms,
		arg.MinYear,
		arg.MaxYear,
		arg.MinFloor,
		arg.MaxFloor,
		arg.MinTotalFloors,
		arg.MaxTotalFloors,
		arg.HasWarehouse,
		arg.HasElevator,
		arg.HasParking,
		arg.MinLat,
		arg.MaxLat,
		arg.MinLng,
		arg.MaxLng,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ad
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.ID,
			&i.PublisherAdKey,
			&i.PublisherID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublishedAt,
			&i.Category,
			&i.Author,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.City,
			&i.Neighborhood,
			&i.HouseType,
			&i.Meterage,
			&i.RoomsCount,
			&i.Year,
			&i.Floor,
			&i.TotalFloors,
			&i.HasWarehouse,
			&i.HasElevator,
			&i.HasParking,
			&i.Lat,
			&i.Lng,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAdByID = `-- name: GetAdByID :one
SELECT ad.id, ad.publisher_ad_key, ad.publisher_id, ad.created_at, ad.updated_at, ad.published_at, ad.category, ad.author, ad.url, ad.title, ad.description, ad.city, ad.neighborhood, ad.house_type, ad.meterage, ad.rooms_count, ad.year, ad.floor, ad.total_floors, ad.has_warehouse, ad.has_elevator, ad.has_parking, ad.lat, ad.lng
FROM ad
WHERE ad.id = $1
`

// Get Ad by its ID
func (q *Queries) GetAdByID(ctx context.Context, id int64) (Ad, error) {
	row := q.db.QueryRow(ctx, getAdByID, id)
	var i Ad
	err := row.Scan(
		&i.ID,
		&i.PublisherAdKey,
		&i.PublisherID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublishedAt,
		&i.Category,
		&i.Author,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.City,
		&i.Neighborhood,
		&i.HouseType,
		&i.Meterage,
		&i.RoomsCount,
		&i.Year,
		&i.Floor,
		&i.TotalFloors,
		&i.HasWarehouse,
		&i.HasElevator,
		&i.HasParking,
		&i.Lat,
		&i.Lng,
	)
	return i, err
}

const getAdsByIds = `-- name: GetAdsByIds :many
SELECT id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
FROM ad
WHERE id = ANY ($1::bigint[])
ORDER BY created_at DESC
`

// Get ads based on list of IDs
func (q *Queries) GetAdsByIds(ctx context.Context, adIds []int64) ([]Ad, error) {
	rows, err := q.db.Query(ctx, getAdsByIds, adIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ad
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.ID,
			&i.PublisherAdKey,
			&i.PublisherID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublishedAt,
			&i.Category,
			&i.Author,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.City,
			&i.Neighborhood,
			&i.HouseType,
			&i.Meterage,
			&i.RoomsCount,
			&i.Year,
			&i.Floor,
			&i.TotalFloors,
			&i.HasWarehouse,
			&i.HasElevator,
			&i.HasParking,
			&i.Lat,
			&i.Lng,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAdsByPublisher = `-- name: GetAdsByPublisher :many
SELECT id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
FROM ad
WHERE publisher_id = $1
ORDER BY created_at DESC
`

// Get ads by publisher ID
func (q *Queries) GetAdsByPublisher(ctx context.Context, publisherID *int32) ([]Ad, error) {
	rows, err := q.db.Query(ctx, getAdsByPublisher, publisherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ad
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.ID,
			&i.PublisherAdKey,
			&i.PublisherID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublishedAt,
			&i.Category,
			&i.Author,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.City,
			&i.Neighborhood,
			&i.HouseType,
			&i.Meterage,
			&i.RoomsCount,
			&i.Year,
			&i.Floor,
			&i.TotalFloors,
			&i.HasWarehouse,
			&i.HasElevator,
			&i.HasParking,
			&i.Lat,
			&i.Lng,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAdsPublisherByAdKey = `-- name: GetAdsPublisherByAdKey :one
SELECT ad.publisher_ad_key
FROM ad
WHERE ad.id = $1
`

// Get PublisherAdKey for one specific ad
func (q *Queries) GetAdsPublisherByAdKey(ctx context.Context, adKey *int64) (string, error) {
	row := q.db.QueryRow(ctx, getAdsPublisherByAdKey, adKey)
	var publisher_ad_key string
	err := row.Scan(&publisher_ad_key)
	return publisher_ad_key, err
}

const getAllAds = `-- name: GetAllAds :many
SELECT id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
FROM ad
ORDER BY id DESC
LIMIT $2 OFFSET $1
`

type GetAllAdsParams struct {
	Offset *int32 `json:"offset"`
	Limit  *int32 `json:"limit"`
}

// Get all ads with dynamic ordering, limit, and offset
func (q *Queries) GetAllAds(ctx context.Context, arg GetAllAdsParams) ([]Ad, error) {
	rows, err := q.db.Query(ctx, getAllAds, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ad
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.ID,
			&i.PublisherAdKey,
			&i.PublisherID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublishedAt,
			&i.Category,
			&i.Author,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.City,
			&i.Neighborhood,
			&i.HouseType,
			&i.Meterage,
			&i.RoomsCount,
			&i.Year,
			&i.Floor,
			&i.TotalFloors,
			&i.HasWarehouse,
			&i.HasElevator,
			&i.HasParking,
			&i.Lat,
			&i.Lng,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAd = `-- name: UpdateAd :one
UPDATE ad
SET publisher_ad_key = COALESCE($1, publisher_ad_key),
    publisher_id     = COALESCE($2, publisher_id),
    updated_at       = NOW(),
    published_at     = COALESCE($3, published_at),
    category         = COALESCE($4, category),
    author           = COALESCE($5, author),
    url              = COALESCE($6, url),
    title            = COALESCE($7, title),
    description      = COALESCE($8, description),
    city             = COALESCE($9, city),
    neighborhood     = COALESCE($10, neighborhood),
    house_type       = COALESCE($11, house_type),
    meterage         = COALESCE($12, meterage),
    rooms_count      = COALESCE($13, rooms_count),
    year             = COALESCE($14, year),
    floor            = COALESCE($15, floor),
    total_floors     = COALESCE($16, total_floors),
    has_warehouse    = COALESCE($17, has_warehouse),
    has_elevator     = COALESCE($18, has_elevator),
    has_parking      = COALESCE($19, has_parking),
    lat              = COALESCE($20, lat),
    lng              = COALESCE($21, lng)
WHERE publisher_ad_key = $1
RETURNING id, publisher_ad_key, publisher_id, created_at, updated_at, published_at, category, author, url, title, description, city, neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, has_parking, lat, lng
`

type UpdateAdParams struct {
	PublisherAdKey *string   `json:"publisher_ad_key"`
	PublisherID    *int32    `json:"publisher_id"`
	PublishedAt    time.Time `json:"published_at"`
	Category       string    `json:"category"`
	Author         *string   `json:"author"`
	Url            *string   `json:"url"`
	Title          *string   `json:"title"`
	Description    *string   `json:"description"`
	City           *string   `json:"city"`
	Neighborhood   *string   `json:"neighborhood"`
	HouseType      string    `json:"house_type"`
	Meterage       *int32    `json:"meterage"`
	RoomsCount     *int32    `json:"rooms_count"`
	Year           *int32    `json:"year"`
	Floor          *int32    `json:"floor"`
	TotalFloors    *int32    `json:"total_floors"`
	HasWarehouse   *bool     `json:"has_warehouse"`
	HasElevator    *bool     `json:"has_elevator"`
	HasParking     *bool     `json:"has_parking"`
	Lat            *float64  `json:"lat"`
	Lng            *float64  `json:"lng"`
}

// Update an existing ad's details with optional fields
func (q *Queries) UpdateAd(ctx context.Context, arg UpdateAdParams) (Ad, error) {
	row := q.db.QueryRow(ctx, updateAd,
		arg.PublisherAdKey,
		arg.PublisherID,
		arg.PublishedAt,
		arg.Category,
		arg.Author,
		arg.Url,
		arg.Title,
		arg.Description,
		arg.City,
		arg.Neighborhood,
		arg.HouseType,
		arg.Meterage,
		arg.RoomsCount,
		arg.Year,
		arg.Floor,
		arg.TotalFloors,
		arg.HasWarehouse,
		arg.HasElevator,
		arg.HasParking,
		arg.Lat,
		arg.Lng,
	)
	var i Ad
	err := row.Scan(
		&i.ID,
		&i.PublisherAdKey,
		&i.PublisherID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublishedAt,
		&i.Category,
		&i.Author,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.City,
		&i.Neighborhood,
		&i.HouseType,
		&i.Meterage,
		&i.RoomsCount,
		&i.Year,
		&i.Floor,
		&i.TotalFloors,
		&i.HasWarehouse,
		&i.HasElevator,
		&i.HasParking,
		&i.Lat,
		&i.Lng,
	)
	return i, err
}
