package repository

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error)
	UpdateTrip(ctx context.Context, trip *model.Trip) error
	DeleteTrip(ctx context.Context, id uuid.UUID) error
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error)
	GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error)
}

type TripRepositoryImpl struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &TripRepositoryImpl{db: db}
}

func (r *TripRepositoryImpl) CreateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

func (r *TripRepositoryImpl) GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error) {
	var trip model.Trip
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.Operator").
		Preload("Bus").
		Preload("Bus.Operator").
		First(&trip, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *TripRepositoryImpl) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Model(trip).Updates(trip).Error
}

func (r *TripRepositoryImpl) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Trip{}, "id = ?", id).Error
}

func (r *TripRepositoryImpl) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	var results []model.TripDetail
	var total int64

	query := r.db.WithContext(ctx).Table("trips t").
		Select(`
			t.id, t.route_id, t.bus_id, t.departure_time, t.arrival_time, 
			t.base_price, t.status,
			r.origin, r.destination, r.distance_km,
			b.model as bus_model, b.amenities_json,
			o.id as operator_id, o.name as operator_name, o.logo as operator_logo,
			COALESCE((SELECT COUNT(*) FROM seats WHERE bus_id = b.id), 0) as total_seats
		`).
		Joins("JOIN routes r ON t.route_id = r.id").
		Joins("JOIN buses b ON t.bus_id = b.id").
		Joins("JOIN operators o ON r.operator_id = o.id").
		Where("r.origin ILIKE ? AND r.destination ILIKE ?", "%"+req.Origin+"%", "%"+req.Destination+"%").
		Where("DATE(t.departure_time) = DATE(?)", req.Date).
		Where("t.is_active = ? AND r.is_active = ? AND b.is_active = ?", true, true, true).
		Where("t.status IN (?)", []string{"scheduled"})

	// Advanced filters
	if req.DepartureTimeStart != nil {
		query = query.Where("t.departure_time >= ?", *req.DepartureTimeStart)
	}
	if req.DepartureTimeEnd != nil {
		query = query.Where("t.departure_time <= ?", *req.DepartureTimeEnd)
	}
	if req.MinPrice != nil {
		query = query.Where("t.base_price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("t.base_price <= ?", *req.MaxPrice)
	}
	if len(req.OperatorIDs) > 0 {
		query = query.Where("o.id IN (?)", req.OperatorIDs)
	}

	// Amenities filter (if bus has all requested amenities)
	if len(req.Amenities) > 0 {
		for _, amenity := range req.Amenities {
			query = query.Where("b.amenities_json LIKE ?", "%"+string(amenity)+"%")
		}
	}

	// Seat types filter (check if bus has seats of requested types)
	if len(req.SeatTypes) > 0 {
		seatTypeStrs := make([]string, len(req.SeatTypes))
		for i, st := range req.SeatTypes {
			seatTypeStrs[i] = string(st)
		}
		query = query.Where(`EXISTS (
			SELECT 1 FROM seats s 
			WHERE s.bus_id = b.id 
			AND s.seat_type IN (?)
		)`, seatTypeStrs)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "t.departure_time"
	if req.SortBy != "" {
		switch req.SortBy {
		case "price":
			sortBy = "t.base_price"
		case "departure_time":
			sortBy = "t.departure_time"
		case "duration":
			sortBy = "(t.arrival_time - t.departure_time)"
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
	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	for rows.Next() {
		var result model.TripDetail
		var amenitiesJSON string
		var operatorLogo *string
		err := rows.Scan(
			&result.ID, &result.RouteID, &result.BusID, &result.DepartureTime, &result.ArrivalTime,
			&result.BasePrice, &result.Status,
			&result.Origin, &result.Destination, &result.DistanceKm,
			&result.BusModel, &amenitiesJSON,
			&result.OperatorID, &result.OperatorName, &operatorLogo,
			&result.TotalSeats,
		)
		if err != nil {
			return nil, 0, err
		}

		if operatorLogo != nil {
			result.OperatorLogo = *operatorLogo
		}

		// Calculate duration
		duration := result.ArrivalTime.Sub(result.DepartureTime)
		result.DurationMinutes = int(duration.Minutes())

		// TODO: Calculate available seats by checking bookings
		result.AvailableSeats = result.TotalSeats

		results = append(results, result)
	}

	return results, total, nil
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
