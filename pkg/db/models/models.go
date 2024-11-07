package models

import (
	"time"
)

// UserRole is an Enum type
type UserRole string

// AdCategory is an Enum type
type AdCategory string

// HouseType is an Enum type
type HouseType string

const (
	SuperAdmin UserRole = "super_admin"
	Admin      UserRole = "admin"
	Simple     UserRole = "simple"

	Rent     AdCategory = "rent"
	Buy      AdCategory = "buy"
	Mortgage AdCategory = "mortgage"
	Other    AdCategory = "other"

	Apartment HouseType = "apartment"
	Villa     HouseType = "villa"
	OtherType HouseType = "other"
)

// Publisher represents the publisher table
type Publisher struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required" sql:"type:varchar(31);not null"`
	URL  string `json:"url" binding:"required" sql:"type:varchar(63);not null"`
}

// Ad represents the ad table
type Ad struct {
	ID             int64      `json:"id"`
	PublisherAdKey string     `json:"publisher_ad_key" binding:"required" sql:"type:varchar(255);unique;not null"`
	PublisherID    *int       `json:"publisher_id"`
	CreatedAt      time.Time  `json:"created_at" sql:"default:now()"`
	UpdatedAt      *time.Time `json:"updated_at"`
	PublishedAt    *time.Time `json:"published_at"`
	Category       AdCategory `json:"category"`
	Author         *string    `json:"author" sql:"type:varchar(63)"`
	URL            *string    `json:"url" sql:"type:varchar(255)"`
	Title          *string    `json:"title" sql:"type:varchar(255)"`
	Description    *string    `json:"description" sql:"type:text"`
	City           *string    `json:"city" sql:"type:varchar(63)"`
	Neighborhood   *string    `json:"neighborhood" sql:"type:varchar(63)"`
	HouseType      HouseType  `json:"house_type"`
	Meterage       int        `json:"meterage" sql:"check:meterage >= 0"`
	RoomsCount     int        `json:"rooms_count" sql:"check:rooms_count >= 0"`
	Year           int        `json:"year" sql:"check:year >= 0"`
	Floor          *int       `json:"floor"`
	TotalFloors    *int       `json:"total_floors"`
	HasWarehouse   bool       `json:"has_warehouse"`
	HasElevator    bool       `json:"has_elevator"`
	Lat            *float64   `json:"lat" sql:"check:lat between -90 and 90"`
	Lng            *float64   `json:"lng" sql:"check:lng between -180 and 180"`
}

// User represents the user table
type User struct {
	TgID            string   `json:"tg_id" binding:"required" sql:"type:varchar(31);primary_key"`
	Role            UserRole `json:"role"`
	WatchlistPeriod int      `json:"watchlist_period"`
}

// Price represents the price table
type Price struct {
	ID            int       `json:"id"`
	AdID          int64     `json:"ad_id" binding:"required"`
	FetchedAt     time.Time `json:"fetched_at" sql:"default:now()"`
	HasPrice      bool      `json:"has_price"`
	TotalPrice    *int64    `json:"total_price" sql:"check:total_price >= 0"`
	PricePerMeter *int64    `json:"price_per_meter" sql:"check:price_per_meter >= 0"`
	Mortgage      *int64    `json:"mortgage" sql:"check:mortgage >= 0"`
	NormalPrice   *int64    `json:"normal_price" sql:"check:normal_price >= 0"`
	WeekendPrice  *int64    `json:"weekend_price" sql:"check:weekend_price >= 0"`
}

// AdPicture represents the ad_picture table
type AdPicture struct {
	ID   int64  `json:"id"`
	AdID int64  `json:"ad_id" binding:"required"`
	URL  string `json:"url" binding:"required" sql:"type:varchar(255)"`
}

// FavoriteAds represents the favorite_ads table
type FavoriteAds struct {
	ID     int64  `json:"id"`
	UserID string `json:"user_id" binding:"required"`
	AdID   int64  `json:"ad_id" binding:"required"`
}
