package service

import (
	"context"
	"testing"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRouteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &RouteServiceImpl{}, service)
}

func TestGetRouteByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	expectedRoute := &model.Route{
		BaseModel:        model.BaseModel{ID: routeID},
		Origin:           "Ha Noi",
		Destination:      "Da Nang",
		DistanceKm:       750,
		EstimatedMinutes: 720,
	}

	mockRepo.EXPECT().
		GetRoutesWithRouteStops(ctx, routeID).
		Return(expectedRoute, nil).
		Times(1)

	result, err := service.GetRouteByID(ctx, routeID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Ha Noi", result.Origin)
	assert.Equal(t, "Da Nang", result.Destination)
}

func TestGetRouteByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	mockRepo.EXPECT().
		GetRoutesWithRouteStops(ctx, routeID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetRouteByID(ctx, routeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get route")
}

func TestListRoutes_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	req := &model.ListRoutesRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	expectedRoutes := []model.Route{
		{BaseModel: model.BaseModel{ID: uuid.New()}, Origin: "Ha Noi"},
		{BaseModel: model.BaseModel{ID: uuid.New()}, Origin: "Ho Chi Minh"},
	}

	mockRepo.EXPECT().
		ListRoutes(ctx, req).
		Return(expectedRoutes, int64(2), nil).
		Times(1)

	routes, total, err := service.ListRoutes(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, routes, 2)
	assert.Equal(t, int64(2), total)
}

func TestListRoutes_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	req := &model.ListRoutesRequest{}

	mockRepo.EXPECT().
		ListRoutes(ctx, req).
		Return(nil, int64(0), assert.AnError).
		Times(1)

	routes, total, err := service.ListRoutes(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, int64(0), total)
}

func TestGetRoutesByOriginDestination_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	origin := "Ha Noi"
	destination := "Da Nang"

	expectedRoutes := []model.Route{
		{BaseModel: model.BaseModel{ID: uuid.New()}, Origin: origin, Destination: destination},
	}

	mockRepo.EXPECT().
		GetRoutesByOriginDestination(ctx, origin, destination).
		Return(expectedRoutes, nil).
		Times(1)

	routes, err := service.GetRoutesByOriginDestination(ctx, origin, destination)

	assert.NoError(t, err)
	assert.Len(t, routes, 1)
	assert.Equal(t, origin, routes[0].Origin)
	assert.Equal(t, destination, routes[0].Destination)
}

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	req := &model.CreateRouteRequest{
		Origin:           "Ha Noi",
		Destination:      "Da Nang",
		DistanceKm:       750,
		EstimatedMinutes: 720,
		RouteStops: []model.CreateRouteStopRequest{
			{StopOrder: 200, Location: "Vinh", StopType: "pickup", OffsetMinutes: 180},
			{StopOrder: 100, Location: "Thanh Hoa", StopType: "pickup", OffsetMinutes: 120},
		},
	}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Do(func(_ context.Context, route *model.Route) {
			assert.Equal(t, "Ha Noi", route.Origin)
			assert.Equal(t, "Da Nang", route.Destination)
			// Verify stops are sorted and normalized
			assert.Len(t, route.RouteStops, 2)
			assert.Equal(t, 100, route.RouteStops[0].StopOrder) // Normalized to 100
			assert.Equal(t, "Thanh Hoa", route.RouteStops[0].Location)
			assert.Equal(t, 200, route.RouteStops[1].StopOrder) // Normalized to 200
			assert.Equal(t, "Vinh", route.RouteStops[1].Location)
		}).
		Return(nil).
		Times(1)

	result, err := service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreate_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	req := &model.CreateRouteRequest{
		Origin:      "Ha Noi",
		Destination: "Da Nang",
	}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create route")
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	existingRoute := &model.Route{
		BaseModel:        model.BaseModel{ID: routeID},
		Origin:           "Ha Noi",
		Destination:      "Da Nang",
		DistanceKm:       750,
		EstimatedMinutes: 720,
	}

	newDistance := 800.0
	newTime := 740

	req := &model.UpdateRouteRequest{
		DistanceKm:       &newDistance,
		EstimatedMinutes: &newTime,
	}

	mockRepo.EXPECT().
		GetRouteByID(ctx, routeID).
		Return(existingRoute, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Do(func(_ context.Context, route *model.Route) {
			assert.Equal(t, 800.0, route.DistanceKm)
			assert.Equal(t, 740, route.EstimatedMinutes)
		}).
		Return(nil).
		Times(1)

	result, err := service.Update(ctx, routeID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 800.0, result.DistanceKm)
	assert.Equal(t, 740, result.EstimatedMinutes)
}

func TestUpdate_NegativeDistance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	existingRoute := &model.Route{
		BaseModel:  model.BaseModel{ID: routeID},
		DistanceKm: 750,
	}

	negativeDistance := -10.0
	req := &model.UpdateRouteRequest{
		DistanceKm: &negativeDistance,
	}

	mockRepo.EXPECT().
		GetRouteByID(ctx, routeID).
		Return(existingRoute, nil).
		Times(1)

	result, err := service.Update(ctx, routeID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "distance must be positive")
}

func TestUpdate_NegativeTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	existingRoute := &model.Route{
		BaseModel:        model.BaseModel{ID: routeID},
		EstimatedMinutes: 720,
	}

	negativeTime := -5
	req := &model.UpdateRouteRequest{
		EstimatedMinutes: &negativeTime,
	}

	mockRepo.EXPECT().
		GetRouteByID(ctx, routeID).
		Return(existingRoute, nil).
		Times(1)

	result, err := service.Update(ctx, routeID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "estimated time must be positive")
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	mockRepo.EXPECT().
		Delete(ctx, routeID).
		Return(nil).
		Times(1)

	err := service.Delete(ctx, routeID)

	assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRouteRepository(ctrl)
	service := NewRouteService(mockRepo)

	ctx := context.Background()
	routeID := uuid.New()

	mockRepo.EXPECT().
		Delete(ctx, routeID).
		Return(assert.AnError).
		Times(1)

	err := service.Delete(ctx, routeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete route")
}
