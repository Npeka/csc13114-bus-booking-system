package service

import (
	"context"
	"errors"
	"time"

	"bus-booking/booking-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatLockService interface {
	LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string) error
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
		lockDuration: 15 * time.Minute,
	}
}

func (s *SeatLockServiceImpl) LockSeats(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID, sessionID string) error {
	// Check if any seats are already locked
	for _, seatID := range seatIDs {
		locked, err := s.lockRepo.IsLocked(ctx, tripID, seatID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check seat lock status")
			return err
		}
		if locked {
			return errors.New("one or more seats are already locked")
		}
	}

	log.Info().
		Str("trip_id", tripID.String()).
		Str("session_id", sessionID).
		Int("seat_count", len(seatIDs)).
		Msg("Locking seats")

	return s.lockRepo.LockSeats(ctx, tripID, seatIDs, sessionID, s.lockDuration)
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
			return errors.New("seat is not available")
		}
	}
	return nil
}

func (s *SeatLockServiceImpl) CleanExpiredLocks(ctx context.Context) error {
	log.Info().Msg("Cleaning expired seat locks")
	return s.lockRepo.CleanExpiredLocks(ctx)
}
