package service

import (
	"context"
	"errors"
	"testing"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteService_CreateRoute_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	req := &model.CreateRouteRequest{
		Origin:           "Hà Nội",
		Destination:      "Hải Phòng",
		DistanceKm:       120,
		EstimatedMinutes: 150,
		RouteStops: []model.CreateRouteStopRequest{
			{StopOrder: 100, StopType: "pickup", Location: "Bến xe Giáp Bát"},
			{StopOrder: 200, StopType: "dropoff", Location: "Bến xe Niệm Nghĩa"},
		},
	}

	mockRouteRepo.On("CreateRoute", ctx, mock.AnythingOfType("*model.Route")).Return(nil)

	// Act
	result, err := service.CreateRoute(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Origin, result.Origin)
	assert.Equal(t, req.Destination, result.Destination)
	assert.Equal(t, 2, len(result.RouteStops))
	assert.True(t, result.IsActive)
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_GetRouteByID_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	routeID := uuid.New()
	expectedRoute := &model.Route{
		Origin:      "Hà Nội",
		Destination: "Hải Phòng",
		IsActive:    true,
	}
	expectedRoute.ID = routeID

	mockRouteRepo.On("GetRoutesWithRouteStops", ctx, routeID).Return(expectedRoute, nil)

	// Act
	result, err := service.GetRouteByID(ctx, routeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, routeID, result.ID)
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_ListRoutes_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	route1 := model.Route{Origin: "Hà Nội", Destination: "Hải Phòng"}
	route1.ID = uuid.New()
	route2 := model.Route{Origin: "TP.HCM", Destination: "Vũng Tàu"}
	route2.ID = uuid.New()
	routes := []model.Route{route1, route2}

	req := &model.ListRoutesRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 20,
		},
	}

	mockRouteRepo.On("ListRoutes", ctx, req).Return(routes, int64(2), nil)

	// Act
	result, total, err := service.ListRoutes(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_UpdateRoute_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	routeID := uuid.New()
	existingRoute := &model.Route{
		Origin:      "Hà Nội",
		Destination: "Hải Phòng",
		DistanceKm:  120,
	}
	existingRoute.ID = routeID

	newDistance := 130
	req := &model.UpdateRouteRequest{DistanceKm: &newDistance}

	mockRouteRepo.On("GetRouteByID", ctx, routeID).Return(existingRoute, nil)
	mockRouteRepo.On("UpdateRoute", ctx, mock.AnythingOfType("*model.Route")).Return(nil)

	// Act
	result, err := service.UpdateRoute(ctx, routeID, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newDistance, result.DistanceKm)
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_DeleteRoute_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	routeID := uuid.New()
	mockRouteRepo.On("DeleteRoute", ctx, routeID).Return(nil)

	// Act
	err := service.DeleteRoute(ctx, routeID)

	// Assert
	assert.NoError(t, err)
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_GetRoutesByOriginDestination_Success(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	origin := "Hà Nội"
	destination := "Hải Phòng"
	route := model.Route{Origin: origin, Destination: destination}
	route.ID = uuid.New()
	routes := []model.Route{route}

	mockRouteRepo.On("GetRoutesByOriginDestination", ctx, origin, destination).Return(routes, nil)

	// Act
	result, err := service.GetRoutesByOriginDestination(ctx, origin, destination)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	mockRouteRepo.AssertExpectations(t)
}

func TestRouteService_CreateRoute_RepositoryError(t *testing.T) {
	// Arrange
	mockRouteRepo := new(mocks.MockRouteRepository)
	service := NewRouteService(mockRouteRepo)
	ctx := context.Background()

	req := &model.CreateRouteRequest{
		Origin:           "Hà Nội",
		Destination:      "Hải Phòng",
		DistanceKm:       120,
		EstimatedMinutes: 150,
		RouteStops:       []model.CreateRouteStopRequest{},
	}

	mockRouteRepo.On("CreateRoute", ctx, mock.AnythingOfType("*model.Route")).Return(errors.New("database error"))

	// Act
	result, err := service.CreateRoute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRouteRepo.AssertExpectations(t)
}
