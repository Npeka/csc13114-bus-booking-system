package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConstantsService(t *testing.T) {
	service := NewConstantsService()

	assert.NotNil(t, service)
	assert.IsType(t, &ConstantsServiceImpl{}, service)
}

func TestGetBusConstants(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetBusConstants(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.SeatTypes)
	assert.NotEmpty(t, result.Amenities)
	assert.NotEmpty(t, result.BusTypes)
}

func TestGetRouteConstants(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetRouteConstants(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.StopTypes)
}

func TestGetTripConstants(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetTripConstants(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.TripStatuses)
}

func TestGetSearchFilterConstants(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetSearchFilterConstants(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.SortOptions)
	assert.NotEmpty(t, result.PriceRanges)
	assert.NotEmpty(t, result.TimeSlots)
	assert.NotEmpty(t, result.SeatTypes)
	assert.NotEmpty(t, result.Amenities)
	assert.NotEmpty(t, result.Cities)
}

func TestGetCities(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetCities(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Hà Nội")
	assert.Contains(t, result, "TP. Hồ Chí Minh")
	assert.Contains(t, result, "Đà Nẵng")
}

func TestGetAllConstants(t *testing.T) {
	service := NewConstantsService()
	ctx := context.Background()

	result, err := service.GetAllConstants(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Bus.SeatTypes)
	assert.NotEmpty(t, result.Route.StopTypes)
	assert.NotEmpty(t, result.Trip.TripStatuses)
	assert.NotEmpty(t, result.SearchFilters.Cities)
}
