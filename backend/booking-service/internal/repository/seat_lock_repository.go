package repository

import (
	"context"
	"time"

	"bus-booking/booking-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatLockRepository interface {
	LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string, duration time.Duration) error
	UnlockSeats(ctx context.Context, sessionID string) error
	UnlockSpecificSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) error
	GetLockedSeats(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error)
	IsLocked(ctx context.Context, tripID, seatID uuid.UUID) (bool, error)
	CleanExpiredLocks(ctx context.Context) error
}

type SeatLockRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatLockRepository(db *gorm.DB) SeatLockRepository {
	return &SeatLockRepositoryImpl{db: db}
}

func (r *SeatLockRepositoryImpl) LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string, duration time.Duration) error {
	expiresAt := time.Now().Add(duration)

	locks := make([]model.SeatLock, len(seatIDs))
	for i, seatID := range seatIDs {
		locks[i] = model.SeatLock{
			TripID:    tripID,
			SeatID:    seatID,
			SessionID: sessionID,
			LockedAt:  time.Now(),
			ExpiresAt: expiresAt,
		}
	}

	return r.db.WithContext(ctx).Create(&locks).Error
}

func (r *SeatLockRepositoryImpl) UnlockSeats(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Delete(&model.SeatLock{}).Error
}

func (r *SeatLockRepositoryImpl) UnlockSpecificSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("trip_id = ? AND seat_id IN ?", tripID, seatIDs).
		Delete(&model.SeatLock{}).Error
}

func (r *SeatLockRepositoryImpl) GetLockedSeats(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error) {
	var locks []model.SeatLock
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND expires_at > ?", tripID, time.Now()).
		Find(&locks).Error

	seatIDs := make([]uuid.UUID, len(locks))
	for i, lock := range locks {
		seatIDs[i] = lock.SeatID
	}

	return seatIDs, err
}

func (r *SeatLockRepositoryImpl) IsLocked(ctx context.Context, tripID, seatID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.SeatLock{}).
		Where("trip_id = ? AND seat_id = ? AND expires_at > ?", tripID, seatID, time.Now()).
		Count(&count).Error
	return count > 0, err
}

func (r *SeatLockRepositoryImpl) CleanExpiredLocks(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.SeatLock{}).Error
}
