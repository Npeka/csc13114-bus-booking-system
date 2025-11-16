package repository

import (
	"bus-booking/booking-service/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatStatusRepository interface {
	GetSeatStatusByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error)
	UpdateSeatStatus(ctx context.Context, seatStatus *model.SeatStatus) error
	BulkUpdateSeatStatus(ctx context.Context, seatStatuses []*model.SeatStatus) error
	GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error)
	ReserveSeat(ctx context.Context, tripID, seatID uuid.UUID, userID uuid.UUID, reservationTime time.Duration) error
	ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error
}

type seatStatusRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatStatusRepository(db *gorm.DB) SeatStatusRepository {
	return &seatStatusRepositoryImpl{db: db}
}

func (r *seatStatusRepositoryImpl) GetSeatStatusByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error) {
	var seatStatuses []*model.SeatStatus
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("seat_number").
		Find(&seatStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get seat statuses: %w", err)
	}

	return seatStatuses, nil
}

func (r *seatStatusRepositoryImpl) UpdateSeatStatus(ctx context.Context, seatStatus *model.SeatStatus) error {
	if err := r.db.WithContext(ctx).Save(seatStatus).Error; err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}
	return nil
}

func (r *seatStatusRepositoryImpl) BulkUpdateSeatStatus(ctx context.Context, seatStatuses []*model.SeatStatus) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seatStatus := range seatStatuses {
			if err := tx.Save(seatStatus).Error; err != nil {
				return fmt.Errorf("failed to update seat status: %w", err)
			}
		}
		return nil
	})
}

func (r *seatStatusRepositoryImpl) GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error) {
	var seatStatuses []*model.SeatStatus
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND status = ?", tripID, "available").
		Order("seat_number").
		Find(&seatStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available seats: %w", err)
	}

	return seatStatuses, nil
}

func (r *seatStatusRepositoryImpl) ReserveSeat(ctx context.Context, tripID, seatID uuid.UUID, userID uuid.UUID, reservationTime time.Duration) error {
	expiresAt := time.Now().UTC().Add(reservationTime)

	err := r.db.WithContext(ctx).Model(&model.SeatStatus{}).
		Where("trip_id = ? AND seat_id = ? AND status = ?", tripID, seatID, "available").
		Updates(map[string]interface{}{
			"status":         "reserved",
			"user_id":        userID,
			"reserved_until": &expiresAt,
			"updated_at":     time.Now().UTC(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to reserve seat: %w", err)
	}

	return nil
}

func (r *seatStatusRepositoryImpl) ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&model.SeatStatus{}).
		Where("trip_id = ? AND seat_id = ?", tripID, seatID).
		Updates(map[string]interface{}{
			"status":         "available",
			"user_id":        nil,
			"booking_id":     nil,
			"reserved_until": nil,
			"updated_at":     time.Now().UTC(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to release seat: %w", err)
	}

	return nil
}
