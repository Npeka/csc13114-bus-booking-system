package service

import (
	"context"
	"time"

	"bus-booking/shared/ginext"

	"bus-booking/booking-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatLockService interface {
	LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string) (time.Time, error)
	UnlockSeats(ctx context.Context, sessionID string) error
	GetLockedSeats(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error)
	ValidateSeatAvailability(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) error
	CleanExpiredLocks(ctx context.Context) error
}

type SeatLockServiceImpl struct {
	lockRepo     repository.SeatLockRepository
	lockDuration time.Duration
}

func NewSeatLockService(lockRepo repository.SeatLockRepository) SeatLockService {
	return &SeatLockServiceImpl{
		lockRepo:     lockRepo,
		lockDuration: 5 * time.Minute, // 5 minutes lock duration
	}
}

func (s *SeatLockServiceImpl) LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string) (time.Time, error) {
	// Check if any seats are already locked
	for _, seatID := range seatIDs {
		locked, err := s.lockRepo.IsLocked(ctx, tripID, seatID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check seat lock status")
			return time.Time{}, err
		}
		if locked {
			return time.Time{}, ginext.NewConflictError("one or more seats are already locked")
		}
	}

	log.Info().
		Str("trip_id", tripID.String()).
		Str("session_id", sessionID).
		Int("seat_count", len(seatIDs)).
		Msg("Locking seats")

	err := s.lockRepo.LockSeats(ctx, tripID, seatIDs, sessionID, s.lockDuration)
	if err != nil {
		return time.Time{}, err
	}

	// Return expiration time for frontend countdown
	expiresAt := time.Now().UTC().Add(s.lockDuration)
	return expiresAt, nil
}

func (s *SeatLockServiceImpl) UnlockSeats(ctx context.Context, sessionID string) error {
	log.Info().Str("session_id", sessionID).Msg("Unlocking seats")
	return s.lockRepo.UnlockSeats(ctx, sessionID)
}

func (s *SeatLockServiceImpl) GetLockedSeats(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error) {
	return s.lockRepo.GetLockedSeats(ctx, tripID)
}

func (s *SeatLockServiceImpl) ValidateSeatAvailability(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) error {
	for _, seatID := range seatIDs {
		locked, err := s.lockRepo.IsLocked(ctx, tripID, seatID)
		if err != nil {
			return err
		}
		if locked {
			return ginext.NewConflictError("seat is not available")
		}
	}
	return nil
}

func (s *SeatLockServiceImpl) CleanExpiredLocks(ctx context.Context) error {
	log.Info().Msg("Cleaning expired seat locks")
	return s.lockRepo.CleanExpiredLocks(ctx)
}
