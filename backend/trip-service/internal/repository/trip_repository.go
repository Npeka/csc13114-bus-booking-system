package repository

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripRepository interface {
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripByID(ctx context.Context, req *model.GetTripByIDRequest, id uuid.UUID) (*model.Trip, error)
	ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error)
	GetTripsByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Trip, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error)
	GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error)

	CreateTrip(ctx context.Context, trip *model.Trip) error
	UpdateTrip(ctx context.Context, trip *model.Trip) error
	DeleteTrip(ctx context.Context, id uuid.UUID) error
}

type TripRepositoryImpl struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &TripRepositoryImpl{db: db}
}

func (r *TripRepositoryImpl) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	var trips []model.Trip
	var total int64

	// Build base query with Preload
	query := r.db.WithContext(ctx).Model(&model.Trip{}).
		Preload("Route", "is_active = ?", true).
		Preload("Bus", "is_active = ?", true).
		Joins("JOIN routes ON routes.id = trips.route_id").
		Joins("JOIN buses ON buses.id = trips.bus_id")

	// Base filter - only active trips by default
	query = query.Where("trips.is_active = ?", true)

	// Optional filters
	if req.Origin != nil && *req.Origin != "" {
		query = query.Where("routes.origin ILIKE ?", "%"+*req.Origin+"%")
	}
	if req.Destination != nil && *req.Destination != "" {
		query = query.Where("routes.destination ILIKE ?", "%"+*req.Destination+"%")
	}

	// Status filter (for admin, default to scheduled for public)
	if req.Status != nil && *req.Status != "" {
		query = query.Where("trips.status = ?", *req.Status)
	} else {
		// Default: only show scheduled trips for public search
		query = query.Where("trips.status = ?", "scheduled")
	}

	// Time range filters
	if req.DepartureTimeStart != nil && *req.DepartureTimeStart != "" {
		// Try parsing as ISO8601 first, then HH:MM
		if t, err := time.Parse(time.RFC3339, *req.DepartureTimeStart); err == nil {
			query = query.Where("trips.departure_time >= ?", t)
		}
	}
	if req.DepartureTimeEnd != nil && *req.DepartureTimeEnd != "" {
		if t, err := time.Parse(time.RFC3339, *req.DepartureTimeEnd); err == nil {
			query = query.Where("trips.departure_time <= ?", t)
		}
	}
	if req.ArrivalTimeStart != nil && *req.ArrivalTimeStart != "" {
		if t, err := time.Parse(time.RFC3339, *req.ArrivalTimeStart); err == nil {
			query = query.Where("trips.arrival_time >= ?", t)
		}
	}
	if req.ArrivalTimeEnd != nil && *req.ArrivalTimeEnd != "" {
		if t, err := time.Parse(time.RFC3339, *req.ArrivalTimeEnd); err == nil {
			query = query.Where("trips.arrival_time <= ?", t)
		}
	}

	// Price filters
	if req.MinPrice != nil {
		query = query.Where("trips.base_price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("trips.base_price <= ?", *req.MaxPrice)
	}

	// Amenities filter
	if len(req.Amenities) > 0 {
		for _, amenity := range req.Amenities {
			query = query.Where("? = ANY(buses.amenities)", string(amenity))
		}
	}

	// Seat types filter
	if len(req.SeatTypes) > 0 {
		seatTypeStrs := make([]string, len(req.SeatTypes))
		for i, st := range req.SeatTypes {
			seatTypeStrs[i] = string(st)
		}
		query = query.Where(`EXISTS (
			SELECT 1 FROM seats 
			WHERE seats.bus_id = buses.id 
			AND seats.seat_type IN (?)
		)`, seatTypeStrs)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "trips.departure_time"
	if req.SortBy != "" {
		switch req.SortBy {
		case "price":
			sortBy = "trips.base_price"
		case "departure_time":
			sortBy = "trips.departure_time"
		case "duration":
			sortBy = "(trips.arrival_time - trips.departure_time)"
		}
	}
	sortOrder := "ASC"
	if req.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.PageSize
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Execute query
	if err := query.Find(&trips).Error; err != nil {
		return nil, 0, err
	}

	// Map to TripDetail
	results := make([]model.TripDetail, 0, len(trips))
	for _, trip := range trips {
		detail := model.TripDetail{
			ID:            trip.ID,
			RouteID:       trip.RouteID,
			BusID:         trip.BusID,
			DepartureTime: trip.DepartureTime,
			ArrivalTime:   trip.ArrivalTime,
			BasePrice:     trip.BasePrice,
			Status:        string(trip.Status),
		}

		// Map Route details
		if trip.Route != nil {
			detail.Route = &model.RouteDetail{
				ID:              trip.Route.ID,
				Origin:          trip.Route.Origin,
				Destination:     trip.Route.Destination,
				DistanceKm:      trip.Route.DistanceKm,
				DurationMinutes: trip.Route.EstimatedMinutes,
			}
		}

		// Map Bus details
		if trip.Bus != nil {
			// Count seats to get total seats
			var seatCount int64
			r.db.Model(&model.Seat{}).Where("bus_id = ?", trip.Bus.ID).Count(&seatCount)

			// Convert pq.StringArray to []string
			amenities := make([]string, len(trip.Bus.Amenities))
			copy(amenities, trip.Bus.Amenities)

			detail.Bus = &model.BusDetail{
				ID:         trip.Bus.ID,
				Model:      trip.Bus.Model,
				BusType:    "standard",
				TotalSeats: int(seatCount),
				Amenities:  amenities,
			}
			detail.TotalSeats = int(seatCount)
		}

		// TODO: Calculate available seats by checking bookings
		detail.AvailableSeats = detail.TotalSeats

		results = append(results, detail)
	}

	return results, total, nil
}

func (r *TripRepositoryImpl) GetTripByID(ctx context.Context, req *model.GetTripByIDRequest, id uuid.UUID) (*model.Trip, error) {
	var trip model.Trip
	query := r.db.WithContext(ctx)

	// Preload Route based on request
	if req.PreLoadRoute {
		if req.PreLoadRouteStop {
			query = query.Preload("Route", func(db *gorm.DB) *gorm.DB {
				return db.Preload("RouteStops", func(db *gorm.DB) *gorm.DB {
					return db.Order("stop_order ASC")
				})
			})
		} else {
			query = query.Preload("Route")
		}
	}

	// Preload Bus and Seats based on request
	if req.PreloadBus {
		if req.PreloadSeat {
			query = query.Preload("Bus", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Seats", func(db *gorm.DB) *gorm.DB {
					return db.Order("seat_number ASC")
				})
			})
		} else {
			query = query.Preload("Bus")
		}
	}

	err := query.First(&trip, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *TripRepositoryImpl) ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error) {
	var trips []model.Trip
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Trip{}).
		Preload("Route").
		Preload("Bus")

	// Count total
	countQuery := r.db.WithContext(ctx).Model(&model.Trip{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("departure_time DESC").Find(&trips).Error

	return trips, total, err
}

// GetTripsByIDs fetches trips by a list of IDs (for batch requests)
func (r *TripRepositoryImpl) GetTripsByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Bus").
		Where("id IN ?", ids).
		Find(&trips).Error
	return trips, err
}

func (r *TripRepositoryImpl) GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Where("route_id = ? AND DATE(departure_time) = DATE(?)", routeID, date).
		Where("is_active = ?", true).
		Order("departure_time ASC").
		Find(&trips).Error
	return trips, err
}

func (r *TripRepositoryImpl) GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Where("bus_id = ? AND departure_time >= ? AND departure_time <= ?", busID, startDate, endDate).
		Where("is_active = ?", true).
		Order("departure_time ASC").
		Find(&trips).Error
	return trips, err
}

func (r *TripRepositoryImpl) CreateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

func (r *TripRepositoryImpl) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Model(trip).Updates(trip).Error
}

func (r *TripRepositoryImpl) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Trip{}, "id = ?", id).Error
}
