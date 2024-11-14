// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type AdCategory string

const (
	AdCategoryRent     AdCategory = "rent"
	AdCategoryBuy      AdCategory = "buy"
	AdCategoryMortgage AdCategory = "mortgage"
	AdCategoryOther    AdCategory = "other"
)

func (e *AdCategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AdCategory(s)
	case string:
		*e = AdCategory(s)
	default:
		return fmt.Errorf("unsupported scan type for AdCategory: %T", src)
	}
	return nil
}

type NullAdCategory struct {
	AdCategory AdCategory `json:"ad_category"`
	Valid      bool       `json:"valid"` // Valid is true if AdCategory is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAdCategory) Scan(value interface{}) error {
	if value == nil {
		ns.AdCategory, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AdCategory.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAdCategory) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AdCategory), nil
}

type HouseType string

const (
	HouseTypeApartment HouseType = "apartment"
	HouseTypeVilla     HouseType = "villa"
	HouseTypeOther     HouseType = "other"
)

func (e *HouseType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = HouseType(s)
	case string:
		*e = HouseType(s)
	default:
		return fmt.Errorf("unsupported scan type for HouseType: %T", src)
	}
	return nil
}

type NullHouseType struct {
	HouseType HouseType `json:"house_type"`
	Valid     bool      `json:"valid"` // Valid is true if HouseType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullHouseType) Scan(value interface{}) error {
	if value == nil {
		ns.HouseType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.HouseType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullHouseType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.HouseType), nil
}

type UserRole string

const (
	UserRoleSuperAdmin UserRole = "super_admin"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSimple     UserRole = "simple"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole `json:"user_role"`
	Valid    bool     `json:"valid"` // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Ad struct {
	ID             int64     `json:"id"`
	PublisherAdKey string    `json:"publisher_ad_key"`
	PublisherID    *int32    `json:"publisher_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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

type AdPicture struct {
	ID   int64   `json:"id"`
	AdID *int64  `json:"ad_id"`
	Url  *string `json:"url"`
}

type FavoriteAd struct {
	ID     int64  `json:"id"`
	UserID string `json:"user_id"`
	AdID   int64  `json:"ad_id"`
}

type Price struct {
	ID            int32     `json:"id"`
	AdID          int64     `json:"ad_id"`
	FetchedAt     time.Time `json:"fetched_at"`
	HasPrice      *bool     `json:"has_price"`
	TotalPrice    *int64    `json:"total_price"`
	PricePerMeter *int64    `json:"price_per_meter"`
	Mortgage      *int64    `json:"mortgage"`
	NormalPrice   *int64    `json:"normal_price"`
	WeekendPrice  *int64    `json:"weekend_price"`
}

type Publisher struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type User struct {
	TgID            string       `json:"tg_id"`
	Role            NullUserRole `json:"role"`
	WatchlistPeriod *int32       `json:"watchlist_period"`
}

type UserAd struct {
	ID     int64  `json:"id"`
	UserID string `json:"user_id"`
	AdID   int64  `json:"ad_id"`
}
