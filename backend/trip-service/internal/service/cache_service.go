package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bus-booking/shared/db"
	"bus-booking/trip-service/internal/model"

	"github.com/rs/zerolog/log"
)

type CacheService interface {
	GetSearchResults(ctx context.Context, key string) (*model.TripSearchResponse, error)
	SetSearchResults(ctx context.Context, key string, response *model.TripSearchResponse, ttl time.Duration) error
	InvalidateTripCache(ctx context.Context, tripID string) error
	InvalidateRouteCache(ctx context.Context, routeID string) error
	InvalidateSearchCache(ctx context.Context, pattern string) error
}

type CacheServiceImpl struct {
	redis *db.RedisManager
}

func NewCacheService(redis *db.RedisManager) CacheService {
	return &CacheServiceImpl{redis: redis}
}

const (
	searchCachePrefix = "trip:search:"
	tripCachePrefix   = "trip:detail:"
	routeCachePrefix  = "route:detail:"
	searchCacheTTL    = 5 * time.Minute
	detailCacheTTL    = 1 * time.Hour
)

func (s *CacheServiceImpl) GetSearchResults(ctx context.Context, key string) (*model.TripSearchResponse, error) {
	cacheKey := searchCachePrefix + key
	data, err := s.redis.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var response model.TripSearchResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal cached search results")
		return nil, err
	}

	log.Debug().Str("key", key).Msg("Cache hit for search results")
	return &response, nil
}

func (s *CacheServiceImpl) SetSearchResults(ctx context.Context, key string, response *model.TripSearchResponse, ttl time.Duration) error {
	cacheKey := searchCachePrefix + key
	data, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal search results")
		return err
	}

	if err := s.redis.Set(ctx, cacheKey, string(data), ttl); err != nil {
		log.Error().Err(err).Msg("Failed to cache search results")
		return err
	}

	log.Debug().Str("key", key).Dur("ttl", ttl).Msg("Cached search results")
	return nil
}

func (s *CacheServiceImpl) InvalidateTripCache(ctx context.Context, tripID string) error {
	cacheKey := tripCachePrefix + tripID
	if err := s.redis.Del(ctx, cacheKey); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to invalidate trip cache")
		return err
	}

	log.Debug().Str("trip_id", tripID).Msg("Invalidated trip cache")
	return nil
}

func (s *CacheServiceImpl) InvalidateRouteCache(ctx context.Context, routeID string) error {
	cacheKey := routeCachePrefix + routeID
	if err := s.redis.Del(ctx, cacheKey); err != nil {
		log.Error().Err(err).Str("route_id", routeID).Msg("Failed to invalidate route cache")
		return err
	}

	log.Debug().Str("route_id", routeID).Msg("Invalidated route cache")
	return nil
}

func (s *CacheServiceImpl) InvalidateSearchCache(ctx context.Context, pattern string) error {
	// This would require SCAN command to find all matching keys
	// For now, we'll just log it as a TODO
	log.Warn().Str("pattern", pattern).Msg("Search cache invalidation by pattern not fully implemented")
	return nil
}

// GenerateSearchCacheKey creates a cache key from search parameters
func GenerateSearchCacheKey(req *model.TripSearchRequest) string {
	return fmt.Sprintf("%s:%s:%s:%d:%d",
		req.Origin,
		req.Destination,
		req.Date.Format("2006-01-02"),
		req.Page,
		req.Limit,
	)
}
