package repository

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripRepository interface {
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error)
	ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error)
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

	// Parse date string
	date, err := time.Parse("02/01/2006", req.Date)
	if err != nil {
		return nil, 0, err
	}

	// Build base query with Preload
	query := r.db.WithContext(ctx).Model(&model.Trip{}).
		Preload("Route", "is_active = ?", true).
		Preload("Bus", "is_active = ?", true).
		Joins("JOIN routes ON routes.id = trips.route_id").
		Joins("JOIN buses ON buses.id = trips.bus_id")

	// Basic filters
	query = query.Where("trips.is_active = ?", true).
		Where("trips.status = ?", "scheduled").
		Where("routes.origin ILIKE ?", "%"+req.Origin+"%").
		Where("routes.destination ILIKE ?", "%"+req.Destination+"%").
		Where("DATE(trips.departure_time) = DATE(?)", date)

	// Advanced filters - parse time strings
	if req.DepartureTimeStart != nil && *req.DepartureTimeStart != "" {
		startTime, err := time.Parse("02/01/2006 15:04", req.Date+" "+*req.DepartureTimeStart)
		if err == nil {
			query = query.Where("trips.departure_time >= ?", startTime)
		}
	}
	if req.DepartureTimeEnd != nil && *req.DepartureTimeEnd != "" {
		endTime, err := time.Parse("02/01/2006 15:04", req.Date+" "+*req.DepartureTimeEnd)
		if err == nil {
			query = query.Where("trips.departure_time <= ?", endTime)
		}
	}
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
		// Map status with display name
		status := constants.TripStatus(trip.Status)
		detail := model.TripDetail{
			ID:            trip.ID,
			RouteID:       trip.RouteID,
			BusID:         trip.BusID,
			DepartureTime: trip.DepartureTime,
			ArrivalTime:   trip.ArrivalTime,
			BasePrice:     trip.BasePrice,
			Status: model.ConstantDisplay{
				Value:       status.String(),
				DisplayName: status.GetDisplayName(),
			},
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

			// Convert pq.StringArray to []ConstantDisplay with display names
			amenities := make([]model.ConstantDisplay, len(trip.Bus.Amenities))
			for i, a := range trip.Bus.Amenities {
				amenity := constants.Amenity(a)
				amenities[i] = model.ConstantDisplay{
					Value:       amenity.String(),
					DisplayName: amenity.GetDisplayName(),
				}
			}

			// Map bus type with display name
			busType := constants.BusTypeStandard // Default, could be derived from model or separate field
			detail.Bus = &model.BusDetail{
				ID:    trip.Bus.ID,
				Model: trip.Bus.Model,
				BusType: model.ConstantDisplay{
					Value:       busType.String(),
					DisplayName: busType.GetDisplayName(),
				},
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

func (r *TripRepositoryImpl) GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error) {
	var trip model.Trip
	err := r.db.WithContext(ctx).
		Preload("Route", func(db *gorm.DB) *gorm.DB {
			return db.Preload("RouteStops", func(db *gorm.DB) *gorm.DB {
				return db.Order("stop_order ASC")
			})
		}).
		Preload("Bus", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Seats", func(db *gorm.DB) *gorm.DB {
				return db.Order("seat_number ASC")
			})
		}).
		First(&trip, "id = ?", id).Error
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
