package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"bus-booking/shared/db/mocks"
	"bus-booking/trip-service/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	assert.NotNil(t, service)
	assert.IsType(t, &CacheServiceImpl{}, service)
}

func TestGetSearchResults_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	key := "search-key"
	cacheKey := "trip:search:" + key

	expectedTrips := []model.TripDetail{{ID: uuid.New()}}
	cachedData := cachedSearchResult{
		Trips: expectedTrips,
		Total: 1,
	}
	cachedJSON, _ := json.Marshal(cachedData)

	mockRedis.EXPECT().
		Get(ctx, cacheKey).
		Return(string(cachedJSON), nil).
		Times(1)

	trips, total, err := service.GetSearchResults(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, trips, 1)
}

func TestGetSearchResults_Miss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	key := "search-key"
	cacheKey := "trip:search:" + key

	mockRedis.EXPECT().
		Get(ctx, cacheKey).
		Return("", assert.AnError). // simulate miss or error
		Times(1)

	trips, total, err := service.GetSearchResults(ctx, key)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, trips)
}

func TestSetSearchResults_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	key := "search-key"
	cacheKey := "trip:search:" + key
	trips := []model.TripDetail{{ID: uuid.New()}}
	total := int64(1)
	ttl := 5 * time.Minute

	mockRedis.EXPECT().
		Set(ctx, cacheKey, gomock.Any(), ttl).
		Return(nil).
		Times(1)

	err := service.SetSearchResults(ctx, key, trips, total, ttl)

	assert.NoError(t, err)
}

func TestInvalidateTripCache_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	tripID := uuid.New().String()
	cacheKey := "trip:detail:" + tripID

	mockRedis.EXPECT().
		Del(ctx, cacheKey).
		Return(nil).
		Times(1)

	err := service.InvalidateTripCache(ctx, tripID)

	assert.NoError(t, err)
}

func TestInvalidateRouteCache_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	routeID := uuid.New().String()
	cacheKey := "route:detail:" + routeID

	mockRedis.EXPECT().
		Del(ctx, cacheKey).
		Return(nil).
		Times(1)

	err := service.InvalidateRouteCache(ctx, routeID)

	assert.NoError(t, err)
}

func TestInvalidateSearchCache_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	service := NewCacheService(mockRedis)

	ctx := context.Background()
	pattern := "test-pattern"

	// Mock nothing needed as implementation is just logging currently
	err := service.InvalidateSearchCache(ctx, pattern)

	assert.NoError(t, err)
}
