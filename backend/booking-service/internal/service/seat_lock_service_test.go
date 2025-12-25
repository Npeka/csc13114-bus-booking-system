package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/booking-service/internal/repository/mocks"
	"bus-booking/shared/ginext"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSeatLockService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	assert.NotNil(t, service)
	impl, ok := service.(*SeatLockServiceImpl)
	assert.True(t, ok)
	assert.Equal(t, 5*time.Minute, impl.lockDuration)
}

func TestLockSeats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()
	seatIDs := []uuid.UUID{seatID1, seatID2}
	sessionID := "test-session-123"

	// Expect IsLocked calls for each seat
	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID1).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID2).
		Return(false, nil).
		Times(1)

	// Expect LockSeats call
	mockRepo.EXPECT().
		LockSeats(ctx, tripID, seatIDs, sessionID, 5*time.Minute).
		Return(nil).
		Times(1)

	expiresAt, err := service.LockSeats(ctx, tripID, seatIDs, sessionID)

	assert.NoError(t, err)
	assert.False(t, expiresAt.IsZero())
	assert.True(t, expiresAt.After(time.Now()))
}

func TestLockSeats_SeatAlreadyLocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}
	sessionID := "test-session-123"

	// First seat is already locked
	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID).
		Return(true, nil).
		Times(1)

	expiresAt, err := service.LockSeats(ctx, tripID, seatIDs, sessionID)

	assert.Error(t, err)
	assert.True(t, expiresAt.IsZero())
	assert.Contains(t, err.Error(), "already locked")
}

func TestLockSeats_IsLockedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}
	sessionID := "test-session-123"

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID).
		Return(false, expectedErr).
		Times(1)

	expiresAt, err := service.LockSeats(ctx, tripID, seatIDs, sessionID)

	assert.Error(t, err)
	assert.True(t, expiresAt.IsZero())
	assert.Equal(t, expectedErr, err)
}

func TestLockSeats_LockSeatsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}
	sessionID := "test-session-123"

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		LockSeats(ctx, tripID, seatIDs, sessionID, 5*time.Minute).
		Return(expectedErr).
		Times(1)

	expiresAt, err := service.LockSeats(ctx, tripID, seatIDs, sessionID)

	assert.Error(t, err)
	assert.True(t, expiresAt.IsZero())
	assert.Equal(t, expectedErr, err)
}

func TestUnlockSeats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	sessionID := "test-session-123"

	mockRepo.EXPECT().
		UnlockSeats(ctx, sessionID).
		Return(nil).
		Times(1)

	err := service.UnlockSeats(ctx, sessionID)

	assert.NoError(t, err)
}

func TestUnlockSeats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	sessionID := "test-session-123"

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		UnlockSeats(ctx, sessionID).
		Return(expectedErr).
		Times(1)

	err := service.UnlockSeats(ctx, sessionID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetLockedSeats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()
	expectedSeats := []uuid.UUID{seatID1, seatID2}

	mockRepo.EXPECT().
		GetLockedSeats(ctx, tripID).
		Return(expectedSeats, nil).
		Times(1)

	result, err := service.GetLockedSeats(ctx, tripID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSeats, result)
	assert.Len(t, result, 2)
}

func TestGetLockedSeats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		GetLockedSeats(ctx, tripID).
		Return(nil, expectedErr).
		Times(1)

	result, err := service.GetLockedSeats(ctx, tripID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestValidateSeatAvailability_AllAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()
	seatIDs := []uuid.UUID{seatID1, seatID2}

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID1).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID2).
		Return(false, nil).
		Times(1)

	err := service.ValidateSeatAvailability(ctx, tripID, seatIDs)

	assert.NoError(t, err)
}

func TestValidateSeatAvailability_SeatLocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID).
		Return(true, nil).
		Times(1)

	err := service.ValidateSeatAvailability(ctx, tripID, seatIDs)

	assert.Error(t, err)
	assert.IsType(t, &ginext.Error{}, err)
	assert.Contains(t, err.Error(), "not available")
}

func TestValidateSeatAvailability_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		IsLocked(ctx, tripID, seatID).
		Return(false, expectedErr).
		Times(1)

	err := service.ValidateSeatAvailability(ctx, tripID, seatIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCleanExpiredLocks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().
		CleanExpiredLocks(ctx).
		Return(nil).
		Times(1)

	err := service.CleanExpiredLocks(ctx)

	assert.NoError(t, err)
}

func TestCleanExpiredLocks_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatLockRepository(ctrl)
	service := NewSeatLockService(mockRepo)

	ctx := context.Background()

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		CleanExpiredLocks(ctx).
		Return(expectedErr).
		Times(1)

	err := service.CleanExpiredLocks(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
