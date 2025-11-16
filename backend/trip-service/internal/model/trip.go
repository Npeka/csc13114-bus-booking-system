package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Operator represents bus operators in the system
type Operator struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	ContactEmail  string     `gorm:"type:varchar(255);unique;not null" json:"contact_email" validate:"required,email"`
	ContactPhone  string     `gorm:"type:varchar(20)" json:"contact_phone"`
	Status        string     `gorm:"type:varchar(50);not null;default:'pending'" json:"status" validate:"oneof=pending approved rejected"`
	ApprovedAt    *time.Time `gorm:"type:timestamptz" json:"approved_at,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Routes []Route `gorm:"foreignKey:OperatorID" json:"routes,omitempty"`
	Buses  []Bus   `gorm:"foreignKey:OperatorID" json:"buses,omitempty"`
}

// BeforeCreate sets UUID for new operators
func (o *Operator) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// Route represents bus routes in the system
type Route struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OperatorID       uuid.UUID `gorm:"type:uuid;not null" json:"operator_id" validate:"required"`
	Origin           string    `gorm:"type:varchar(255);not null" json:"origin" validate:"required"`
	Destination      string    `gorm:"type:varchar(255);not null" json:"destination" validate:"required"`
	DistanceKm       int       `gorm:"type:integer;not null" json:"distance_km" validate:"required,min=1"`
	EstimatedMinutes int       `gorm:"type:integer;not null" json:"estimated_minutes" validate:"required,min=1"`
	IsActive         bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Operator *Operator `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"operator,omitempty"`
	Trips    []Trip    `gorm:"foreignKey:RouteID" json:"trips,omitempty"`
}

// BeforeCreate sets UUID for new routes
func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Bus represents buses in the system
type Bus struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OperatorID     uuid.UUID `gorm:"type:uuid;not null" json:"operator_id" validate:"required"`
	PlateNumber    string    `gorm:"type:varchar(20);unique;not null" json:"plate_number" validate:"required"`
	Model          string    `gorm:"type:varchar(255);not null" json:"model" validate:"required"`
	SeatCapacity   int       `gorm:"type:integer;not null" json:"seat_capacity" validate:"required,min=1,max=100"`
	AmenitiesJSON  string    `gorm:"type:text" json:"-"`
	Amenities      []string  `gorm:"-" json:"amenities"`
	IsActive       bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Operator *Operator `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"operator,omitempty"`
	Seats    []Seat    `gorm:"foreignKey:BusID" json:"seats,omitempty"`
	Trips    []Trip    `gorm:"foreignKey:BusID" json:"trips,omitempty"`
}

// BeforeCreate sets UUID for new buses
func (b *Bus) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BeforeSave handles JSON marshaling for amenities
func (b *Bus) BeforeSave(tx *gorm.DB) error {
	if len(b.Amenities) > 0 {
		amenitiesJSON, err := json.Marshal(b.Amenities)
		if err != nil {
			return err
		}
		b.AmenitiesJSON = string(amenitiesJSON)
	}
	return nil
}

// AfterFind handles JSON unmarshaling for amenities
func (b *Bus) AfterFind(tx *gorm.DB) error {
	if b.AmenitiesJSON != "" {
		return json.Unmarshal([]byte(b.AmenitiesJSON), &b.Amenities)
	}
	return nil
}

// Trip represents individual trips in the system
type Trip struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RouteID       uuid.UUID `gorm:"type:uuid;not null" json:"route_id" validate:"required"`
	BusID         uuid.UUID `gorm:"type:uuid;not null" json:"bus_id" validate:"required"`
	DepartureTime time.Time `gorm:"type:timestamptz;not null" json:"departure_time" validate:"required"`
	ArrivalTime   time.Time `gorm:"type:timestamptz;not null" json:"arrival_time" validate:"required"`
	BasePrice     float64   `gorm:"type:decimal(10,2);not null" json:"base_price" validate:"required,min=0"`
	Status        string    `gorm:"type:varchar(50);not null;default:'scheduled'" json:"status" validate:"oneof=scheduled in_progress completed cancelled"`
	IsActive      bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Route *Route `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"route,omitempty"`
	Bus   *Bus   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bus,omitempty"`
}

// BeforeCreate sets UUID for new trips
func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// Seat represents seats in buses
type Seat struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BusID     uuid.UUID `gorm:"type:uuid;not null" json:"bus_id" validate:"required"`
	SeatCode  string    `gorm:"type:varchar(10);not null" json:"seat_code" validate:"required"`
	SeatType  string    `gorm:"type:varchar(50);not null;default:'standard'" json:"seat_type" validate:"oneof=standard premium vip"`
	IsActive  bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Bus *Bus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bus,omitempty"`
}

// BeforeCreate sets UUID for new seats
func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName overrides to ensure consistent naming
func (Operator) TableName() string { return "operators" }
func (Route) TableName() string    { return "routes" }
func (Bus) TableName() string      { return "buses" }
func (Trip) TableName() string     { return "trips" }
func (Seat) TableName() string     { return "seats" }