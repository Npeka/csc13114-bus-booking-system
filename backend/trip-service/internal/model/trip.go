package model

import (
	"bus-booking/trip-service/internal/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Trip struct {
	BaseModel
	RouteID       uuid.UUID            `gorm:"type:uuid;not null" json:"route_id" validate:"required"`
	BusID         uuid.UUID            `gorm:"type:uuid;not null" json:"bus_id" validate:"required"`
	DepartureTime time.Time            `gorm:"type:timestamptz;not null" json:"departure_time" validate:"required"`
	ArrivalTime   time.Time            `gorm:"type:timestamptz;not null" json:"arrival_time" validate:"required"`
	BasePrice     float64              `gorm:"type:decimal(10,2);not null" json:"base_price" validate:"required,min=0"`
	Status        constants.TripStatus `gorm:"type:varchar(50);not null;default:'scheduled'" json:"status" validate:"required"`
	IsActive      bool                 `gorm:"type:boolean;not null;default:true" json:"is_active"`

	Route *Route `gorm:"constraint:OnUpdate:CASCADE" json:"route,omitempty"`
	Bus   *Bus   `gorm:"constraint:OnUpdate:CASCADE" json:"bus,omitempty"`
}

func (Trip) TableName() string {
	return "trips"
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TripSearchRequest struct {
	// Basic search - now all optional
	Origin      *string `form:"origin" json:"origin,omitempty"`
	Destination *string `form:"destination" json:"destination,omitempty"`

	// Advanced filters - time range instead of date
	DepartureTimeStart *string              `form:"departure_time_start" json:"departure_time_start,omitempty"` // Format: ISO8601 or HH:MM
	DepartureTimeEnd   *string              `form:"departure_time_end" json:"departure_time_end,omitempty"`     // Format: ISO8601 or HH:MM
	ArrivalTimeStart   *string              `form:"arrival_time_start" json:"arrival_time_start,omitempty"`     // Format: ISO8601 or HH:MM
	ArrivalTimeEnd     *string              `form:"arrival_time_end" json:"arrival_time_end,omitempty"`         // Format: ISO8601 or HH:MM
	MinPrice           *float64             `form:"min_price" json:"min_price,omitempty" validate:"omitempty,min=0"`
	MaxPrice           *float64             `form:"max_price" json:"max_price,omitempty" validate:"omitempty,min=0"`
	SeatTypes          []constants.SeatType `form:"seat_types" json:"seat_types,omitempty"`
	Amenities          []constants.Amenity  `form:"amenities" json:"amenities,omitempty"`

	// Status filter (for admin)
	Status *string `form:"status" json:"status,omitempty"`

	// Sorting
	SortBy    string `form:"sort_by" json:"sort_by" validate:"omitempty,oneof=price departure_time duration"`
	SortOrder string `form:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"`

	// Pagination
	Page     int `form:"page,default=1" json:"page" validate:"min=1"`
	PageSize int `form:"page_size,default=20" json:"page_size" validate:"min=1,max=100"`
}

type GetTripByIDRequest struct {
	SeatBookingStatus bool `form:"seat_booking_status" json:"seat_booking_status"`
	PreLoadRoute      bool `form:"preload_route" json:"preload_route"`
	PreLoadRouteStop  bool `form:"preload_route_stop" json:"preload_route_stop"`
	PreloadBus        bool `form:"preload_bus" json:"preload_bus"`
	PreloadSeat       bool `form:"preload_seat" json:"preload_seat"`
}

type TripDetail struct {
	ID             uuid.UUID `json:"id"`
	RouteID        uuid.UUID `json:"route_id"`
	BusID          uuid.UUID `json:"bus_id"`
	DepartureTime  time.Time `json:"departure_time"`
	ArrivalTime    time.Time `json:"arrival_time"`
	BasePrice      float64   `json:"base_price"`
	Status         string    `json:"status"` // Raw string value
	AvailableSeats int       `json:"available_seats"`
	TotalSeats     int       `json:"total_seats"`

	Route      *RouteDetail `json:"route,omitempty"`
	Bus        *BusDetail   `json:"bus,omitempty"`
	PriceTiers []PriceTier  `json:"price_tiers,omitempty"`
}

type RouteDetail struct {
	ID              uuid.UUID `json:"id"`
	Origin          string    `json:"origin"`
	Destination     string    `json:"destination"`
	DistanceKm      int       `json:"distance_km"`
	DurationMinutes int       `json:"duration_minutes"`
}

type BusDetail struct {
	ID         uuid.UUID `json:"id"`
	Model      string    `json:"model"`
	BusType    string    `json:"bus_type"` // Raw string value
	TotalSeats int       `json:"total_seats"`
	Amenities  []string  `json:"amenities"` // Raw string values
}

type PriceTier struct {
	SeatType        string  `json:"seat_type"` // Raw string value
	BasePrice       float64 `json:"base_price"`
	PriceMultiplier float64 `json:"price_multiplier"`
	FinalPrice      float64 `json:"final_price"`
	AvailableCount  int     `json:"available_count"`
}

type ListTripsRequest struct {
	PaginationRequest
	IDStrs []string    `form:"ids[]" json:"ids"`
	IDs    []uuid.UUID `form:"-" json:"-"`
}

type TripResponse struct {
	ID            uuid.UUID      `json:"id"`
	RouteID       uuid.UUID      `json:"route_id"`
	BusID         uuid.UUID      `json:"bus_id"`
	DepartureTime time.Time      `json:"departure_time"`
	ArrivalTime   time.Time      `json:"arrival_time"`
	BasePrice     float64        `json:"base_price"`
	Status        string         `json:"status"` // Raw string value
	IsActive      bool           `json:"is_active"`
	Route         *RouteResponse `json:"route,omitempty"`
	Bus           *BusResponse   `json:"bus,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type CreateTripRequest struct {
	RouteID       uuid.UUID `json:"route_id" validate:"required"`
	BusID         uuid.UUID `json:"bus_id" validate:"required"`
	DepartureTime time.Time `json:"departure_time" validate:"required"`
	ArrivalTime   time.Time `json:"arrival_time" validate:"required"`
	BasePrice     float64   `json:"base_price" validate:"required,min=0"`
}

type UpdateTripRequest struct {
	DepartureTime *time.Time            `json:"departure_time,omitempty" validate:"omitempty"`
	ArrivalTime   *time.Time            `json:"arrival_time,omitempty" validate:"omitempty"`
	BasePrice     *float64              `json:"base_price,omitempty" validate:"omitempty,min=0"`
	Status        *constants.TripStatus `json:"status,omitempty" validate:"omitempty,oneof=scheduled in_progress completed cancelled"`
	IsActive      *bool                 `json:"is_active,omitempty"`
}
