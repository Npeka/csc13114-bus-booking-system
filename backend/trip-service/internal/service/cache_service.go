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
	GetSearchResults(ctx context.Context, key string) ([]model.TripDetail, int64, error)
	SetSearchResults(ctx context.Context, key string, trips []model.TripDetail, total int64, ttl time.Duration) error
	GetConstants(ctx context.Context, key string) (interface{}, error)
	SetConstants(ctx context.Context, key string, data interface{}, ttl time.Duration) error
	InvalidateTripCache(ctx context.Context, tripID string) error
	InvalidateRouteCache(ctx context.Context, routeID string) error
	InvalidateSearchCache(ctx context.Context, pattern string) error
}

type CacheServiceImpl struct {
	redis db.RedisManager
}

func NewCacheService(redis db.RedisManager) CacheService {
	return &CacheServiceImpl{redis: redis}
}

const (
	searchCachePrefix    = "trip:search:"
	tripCachePrefix      = "trip:detail:"
	routeCachePrefix     = "route:detail:"
	constantsCachePrefix = "constants:"
	searchCacheTTL       = 5 * time.Minute
	detailCacheTTL       = 1 * time.Hour
	constantsCacheTTL    = 24 * time.Hour // Constants rarely change
)

type cachedSearchResult struct {
	Trips []model.TripDetail `json:"trips"`
	Total int64              `json:"total"`
}

func (s *CacheServiceImpl) GetSearchResults(ctx context.Context, key string) ([]model.TripDetail, int64, error) {
	cacheKey := searchCachePrefix + key
	data, err := s.redis.Get(ctx, cacheKey)
	if err != nil {
		return nil, 0, err
	}

	var result cachedSearchResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal cached search results")
		return nil, 0, err
	}

	log.Debug().Str("key", key).Msg("Cache hit for search results")
	return result.Trips, result.Total, nil
}

func (s *CacheServiceImpl) SetSearchResults(ctx context.Context, key string, trips []model.TripDetail, total int64, ttl time.Duration) error {
	cacheKey := searchCachePrefix + key
	result := cachedSearchResult{
		Trips: trips,
		Total: total,
	}
	data, err := json.Marshal(result)
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

func (s *CacheServiceImpl) GetConstants(ctx context.Context, key string) (interface{}, error) {
	cacheKey := constantsCachePrefix + key
	data, err := s.redis.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal cached constants")
		return nil, err
	}

	log.Debug().Str("key", key).Msg("Cache hit for constants")
	return result, nil
}

func (s *CacheServiceImpl) SetConstants(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	cacheKey := constantsCachePrefix + key
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal constants")
		return err
	}

	if err := s.redis.Set(ctx, cacheKey, string(jsonData), ttl); err != nil {
		log.Error().Err(err).Msg("Failed to cache constants")
		return err
	}

	log.Debug().Str("key", key).Dur("ttl", ttl).Msg("Cached constants")
	return nil
}

// GenerateSearchCacheKey creates a cache key from search parameters
func GenerateSearchCacheKey(req *model.TripSearchRequest) string {
	origin := ""
	if req.Origin != nil {
		origin = *req.Origin
	}
	destination := ""
	if req.Destination != nil {
		destination = *req.Destination
	}
	departureStart := ""
	if req.DepartureTimeStart != nil {
		departureStart = *req.DepartureTimeStart
	}
	departureEnd := ""
	if req.DepartureTimeEnd != nil {
		departureEnd = *req.DepartureTimeEnd
	}
	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s:%d:%d",
		origin,
		destination,
		departureStart,
		departureEnd,
		status,
		req.Page,
		req.PageSize,
	)
}
