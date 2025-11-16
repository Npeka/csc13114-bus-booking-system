package repository

import (
	"context"
	"fmt"
	"time"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripRepositoryInterface interface {
	// Trip operations
	CreateTrip(ctx context.Context, trip *model.Trip) error
	GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error)
	UpdateTrip(ctx context.Context, trip *model.Trip) error
	DeleteTrip(ctx context.Context, id uuid.UUID) error
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error)
	GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error)

	// Route operations
	CreateRoute(ctx context.Context, route *model.Route) error
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	UpdateRoute(ctx context.Context, route *model.Route) error
	DeleteRoute(ctx context.Context, id uuid.UUID) error
	ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.RouteSummary, int64, error)
	GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error)

	// Bus operations
	CreateBus(ctx context.Context, bus *model.Bus) error
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	UpdateBus(ctx context.Context, bus *model.Bus) error
	DeleteBus(ctx context.Context, id uuid.UUID) error
	ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error)
	GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error)

	// Operator operations
	CreateOperator(ctx context.Context, operator *model.Operator) error
	GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error)
	UpdateOperator(ctx context.Context, operator *model.Operator) error
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	ListOperators(ctx context.Context, page, limit int) ([]model.OperatorSummary, int64, error)
	GetOperatorByEmail(ctx context.Context, email string) (*model.Operator, error)

	// Seat operations
	CreateSeats(ctx context.Context, seats []model.Seat) error
	GetSeatsByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
	GetSeatByID(ctx context.Context, id uuid.UUID) (*model.Seat, error)
	UpdateSeat(ctx context.Context, seat *model.Seat) error
	DeleteSeat(ctx context.Context, id uuid.UUID) error
}

type TripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepositoryInterface {
	return &TripRepository{db: db}
}

// Trip operations
func (r *TripRepository) CreateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

func (r *TripRepository) GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error) {
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

func (r *TripRepository) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Model(trip).Updates(trip).Error
}

func (r *TripRepository) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Trip{}, "id = ?", id).Error
}

func (r *TripRepository) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	var results []model.TripDetail
	var total int64

	query := r.db.WithContext(ctx).Table("trips t").
		Select(`
			t.id, t.route_id, t.bus_id, t.departure_time, t.arrival_time, 
			t.base_price, t.status,
			r.origin, r.destination, r.distance_km,
			b.model as bus_model, b.plate_number as bus_plate_number, b.amenities_json,
			o.id as operator_id, o.name as operator_name,
			b.seat_capacity as total_seats
		`).
		Joins("JOIN routes r ON t.route_id = r.id").
		Joins("JOIN buses b ON t.bus_id = b.id").
		Joins("JOIN operators o ON r.operator_id = o.id").
		Where("r.origin ILIKE ? AND r.destination ILIKE ?", "%"+req.Origin+"%", "%"+req.Destination+"%").
		Where("DATE(t.departure_time) = DATE(?)", req.DepartureDate).
		Where("t.is_active = ? AND r.is_active = ? AND b.is_active = ?", true, true, true).
		Where("t.status IN (?)", []string{"scheduled"})

	// Apply filters
	if req.OperatorID != nil {
		query = query.Where("o.id = ?", *req.OperatorID)
	}
	if req.PriceMin != nil {
		query = query.Where("t.base_price >= ?", *req.PriceMin)
	}
	if req.PriceMax != nil {
		query = query.Where("t.base_price <= ?", *req.PriceMax)
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
		case "arrival_time":
			sortBy = "t.arrival_time"
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
	defer rows.Close()

	for rows.Next() {
		var result model.TripDetail
		var amenitiesJSON string
		err := rows.Scan(
			&result.ID, &result.RouteID, &result.BusID, &result.DepartureTime, &result.ArrivalTime,
			&result.BasePrice, &result.Status,
			&result.Origin, &result.Destination, &result.DistanceKm,
			&result.BusModel, &result.BusPlateNumber, &amenitiesJSON,
			&result.OperatorID, &result.OperatorName,
			&result.TotalSeats,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse amenities JSON
		if amenitiesJSON != "" {
			// This would need JSON unmarshaling
		}

		// Calculate duration
		duration := result.ArrivalTime.Sub(result.DepartureTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		result.Duration = fmt.Sprintf("%dh %dm", hours, minutes)

		// TODO: Calculate available seats by checking bookings
		result.AvailableSeats = result.TotalSeats

		results = append(results, result)
	}

	return results, total, nil
}

func (r *TripRepository) GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Where("route_id = ? AND DATE(departure_time) = DATE(?)", routeID, date).
		Where("is_active = ?", true).
		Order("departure_time ASC").
		Find(&trips).Error
	return trips, err
}

func (r *TripRepository) GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Where("bus_id = ? AND departure_time >= ? AND departure_time <= ?", busID, startDate, endDate).
		Where("is_active = ?", true).
		Order("departure_time ASC").
		Find(&trips).Error
	return trips, err
}

// Route operations
func (r *TripRepository) CreateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *TripRepository) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	var route model.Route
	err := r.db.WithContext(ctx).
		Preload("Operator").
		First(&route, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *TripRepository) UpdateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Model(route).Updates(route).Error
}

func (r *TripRepository) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Route{}, "id = ?", id).Error
}

func (r *TripRepository) ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.RouteSummary, int64, error) {
	var results []model.RouteSummary
	var total int64

	query := r.db.WithContext(ctx).Table("routes r").
		Select(`
			r.id, r.origin, r.destination, r.distance_km, r.estimated_minutes,
			r.is_active, r.created_at, o.name as operator_name,
			COUNT(t.id) as active_trips
		`).
		Joins("JOIN operators o ON r.operator_id = o.id").
		Joins("LEFT JOIN trips t ON r.id = t.route_id AND t.is_active = true AND t.departure_time > NOW()").
		Group("r.id, o.name")

	if operatorID != nil {
		query = query.Where("r.operator_id = ?", *operatorID)
	}

	// Count total
	countQuery := r.db.WithContext(ctx).Model(&model.Route{})
	if operatorID != nil {
		countQuery = countQuery.Where("operator_id = ?", *operatorID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("r.created_at DESC")

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var result model.RouteSummary
		err := rows.Scan(
			&result.ID, &result.Origin, &result.Destination, &result.DistanceKm, &result.EstimatedMinutes,
			&result.IsActive, &result.CreatedAt, &result.OperatorName, &result.ActiveTrips,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, nil
}

func (r *TripRepository) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	var routes []model.Route
	err := r.db.WithContext(ctx).
		Where("origin = ? AND destination = ? AND is_active = ?", origin, destination, true).
		Preload("Operator").
		Find(&routes).Error
	return routes, err
}

// Bus operations
func (r *TripRepository) CreateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Create(bus).Error
}

func (r *TripRepository) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	var bus model.Bus
	err := r.db.WithContext(ctx).
		Preload("Operator").
		First(&bus, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *TripRepository) UpdateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Model(bus).Updates(bus).Error
}

func (r *TripRepository) DeleteBus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Bus{}, "id = ?", id).Error
}

func (r *TripRepository) ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error) {
	var buses []model.Bus
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Bus{}).Preload("Operator")
	if operatorID != nil {
		query = query.Where("operator_id = ?", *operatorID)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&buses).Error

	return buses, total, err
}

func (r *TripRepository) GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error) {
	var bus model.Bus
	err := r.db.WithContext(ctx).
		Preload("Operator").
		Where("plate_number = ?", plateNumber).
		First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

// Operator operations
func (r *TripRepository) CreateOperator(ctx context.Context, operator *model.Operator) error {
	return r.db.WithContext(ctx).Create(operator).Error
}

func (r *TripRepository) GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error) {
	var operator model.Operator
	err := r.db.WithContext(ctx).First(&operator, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &operator, nil
}

func (r *TripRepository) UpdateOperator(ctx context.Context, operator *model.Operator) error {
	return r.db.WithContext(ctx).Model(operator).Updates(operator).Error
}

func (r *TripRepository) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Operator{}, "id = ?", id).Error
}

func (r *TripRepository) ListOperators(ctx context.Context, page, limit int) ([]model.OperatorSummary, int64, error) {
	var results []model.OperatorSummary
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.Operator{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Table("operators o").
		Select(`
			o.id, o.name, o.contact_email, o.contact_phone, o.status, o.created_at,
			COUNT(DISTINCT r.id) as active_routes,
			COUNT(DISTINCT b.id) as active_buses
		`).
		Joins("LEFT JOIN routes r ON o.id = r.operator_id AND r.is_active = true").
		Joins("LEFT JOIN buses b ON o.id = b.operator_id AND b.is_active = true").
		Group("o.id").
		Offset(offset).Limit(limit).
		Order("o.created_at DESC")

	rows, err := query.Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var result model.OperatorSummary
		err := rows.Scan(
			&result.ID, &result.Name, &result.ContactEmail, &result.ContactPhone, &result.Status, &result.CreatedAt,
			&result.ActiveRoutes, &result.ActiveBuses,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, nil
}

func (r *TripRepository) GetOperatorByEmail(ctx context.Context, email string) (*model.Operator, error) {
	var operator model.Operator
	err := r.db.WithContext(ctx).Where("contact_email = ?", email).First(&operator).Error
	if err != nil {
		return nil, err
	}
	return &operator, nil
}

// Seat operations
func (r *TripRepository) CreateSeats(ctx context.Context, seats []model.Seat) error {
	return r.db.WithContext(ctx).Create(&seats).Error
}

func (r *TripRepository) GetSeatsByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.WithContext(ctx).
		Where("bus_id = ? AND is_active = ?", busID, true).
		Order("seat_code ASC").
		Find(&seats).Error
	return seats, err
}

func (r *TripRepository) GetSeatByID(ctx context.Context, id uuid.UUID) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.WithContext(ctx).First(&seat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *TripRepository) UpdateSeat(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Model(seat).Updates(seat).Error
}

func (r *TripRepository) DeleteSeat(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Seat{}, "id = ?", id).Error
}
