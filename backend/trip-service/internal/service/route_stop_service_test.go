package service

import (
	"context"
	"testing"

	"bus-booking/trip-service/internal/model"
	repo_mocks "bus-booking/trip-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRouteStopService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &RouteStopServiceImpl{}, service)
}

func TestCreateRouteStop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	routeID := uuid.New()

	req := &model.CreateRouteStopRequest{
		RouteID:       routeID,
		StopOrder:     100,
		StopType:      "pickup",
		Location:      "Thanh Hoa",
		OffsetMinutes: 120,
	}

	route := &model.Route{
		BaseModel: model.BaseModel{ID: routeID},
		RouteStops: []model.RouteStop{
			{StopOrder: 100},
		},
	}

	// createdStop := &model.RouteStop{
	// 	BaseModel: model.BaseModel{ID: uuid.New()},
	// 	RouteID:   routeID,
	// 	StopOrder: 150,
	// }

	mockRouteRepo.EXPECT().
		GetRoutesWithRouteStops(ctx, routeID).
		Return(route, nil).
		Times(1)

	// Expect ReorderStops to shift existing stops (stop order 100 will shift to 101)
	mockRouteStopRepo.EXPECT().
		ReorderStops(ctx, routeID, gomock.Any()).
		Return(nil).
		Times(1)

	mockRouteStopRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	result, err := service.CreateRouteStop(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateRouteStop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()

	existingStop := &model.RouteStop{
		BaseModel: model.BaseModel{ID: stopID},
		Location:  "Old Location",
	}

	newLocation := "New Location"
	req := &model.UpdateRouteStopRequest{
		Location: &newLocation,
	}

	mockRouteStopRepo.EXPECT().
		GetByID(ctx, stopID).
		Return(existingStop, nil).
		Times(1)

	mockRouteStopRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Do(func(_ context.Context, stop *model.RouteStop) {
			assert.Equal(t, "New Location", stop.Location)
		}).
		Return(nil).
		Times(1)

	result, err := service.UpdateRouteStop(ctx, stopID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Location", result.Location)
}

func TestUpdateRouteStop_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()

	req := &model.UpdateRouteStopRequest{}

	mockRouteStopRepo.EXPECT().
		GetByID(ctx, stopID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.UpdateRouteStop(ctx, stopID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteRouteStop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()

	mockRouteStopRepo.EXPECT().
		Delete(ctx, stopID).
		Return(nil).
		Times(1)

	err := service.DeleteRouteStop(ctx, stopID)

	assert.NoError(t, err)
}

func TestDeleteRouteStop_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()

	mockRouteStopRepo.EXPECT().
		Delete(ctx, stopID).
		Return(assert.AnError).
		Times(1)

	err := service.DeleteRouteStop(ctx, stopID)

	assert.Error(t, err)
}

func TestListRouteStops_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	routeID := uuid.New()

	expectedStops := []model.RouteStop{
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 100},
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 200},
	}

	mockRouteStopRepo.EXPECT().
		ListByRouteID(ctx, routeID).
		Return(expectedStops, nil).
		Times(1)

	result, err := service.ListRouteStops(ctx, routeID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListRouteStops_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	routeID := uuid.New()

	mockRouteStopRepo.EXPECT().
		ListByRouteID(ctx, routeID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.ListRouteStops(ctx, routeID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMoveRouteStop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()
	routeID := uuid.New()

	stop := &model.RouteStop{
		BaseModel: model.BaseModel{ID: stopID},
		RouteID:   routeID,
		StopOrder: 200,
	}

	allStops := []model.RouteStop{
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 100},
		{BaseModel: model.BaseModel{ID: stopID}, StopOrder: 200},
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 300},
	}

	req := &model.MoveRouteStopRequest{
		Position:        "before",
		ReferenceStopID: &allStops[0].ID,
	}

	mockRouteStopRepo.EXPECT().GetByID(ctx, stopID).Return(stop, nil).Times(1)
	mockRouteRepo.EXPECT().GetRoutesWithRouteStops(ctx, routeID).Return(&model.Route{RouteStops: allStops}, nil).Times(1)
	mockRouteStopRepo.EXPECT().ReorderStops(ctx, routeID, gomock.Any()).Return(nil).Times(1)

	result, err := service.MoveRouteStop(ctx, stopID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReorderStops_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	routeID := uuid.New()
	stopID1 := uuid.New()
	stopID2 := uuid.New()

	stopOrders := []StopOrder{
		{StopID: stopID2, Order: 100}, // Swapped
		{StopID: stopID1, Order: 200},
	}

	mockRouteStopRepo.EXPECT().ReorderStops(ctx, routeID, gomock.Any()).Return(nil).Times(1)

	err := service.ReorderStops(ctx, routeID, stopOrders)

	assert.NoError(t, err)
}

func TestMoveRouteStop_First(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()
	routeID := uuid.New()

	stop := &model.RouteStop{BaseModel: model.BaseModel{ID: stopID}, RouteID: routeID, StopOrder: 200}
	allStops := []model.RouteStop{
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 100},
		{BaseModel: model.BaseModel{ID: stopID}, StopOrder: 200},
	}

	req := &model.MoveRouteStopRequest{Position: "first"}

	mockRouteStopRepo.EXPECT().GetByID(ctx, stopID).Return(stop, nil).Times(1)
	mockRouteRepo.EXPECT().GetRoutesWithRouteStops(ctx, routeID).Return(&model.Route{RouteStops: allStops}, nil).Times(1)

	// First order is 100, so new order should be 99
	mockRouteStopRepo.EXPECT().ReorderStops(ctx, routeID, map[uuid.UUID]int{stopID: 99}).Return(nil).Times(1)

	result, err := service.MoveRouteStop(ctx, stopID, req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 99, result.StopOrder)
}

func TestMoveRouteStop_Last(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()
	routeID := uuid.New()

	stop := &model.RouteStop{BaseModel: model.BaseModel{ID: stopID}, RouteID: routeID, StopOrder: 100}
	allStops := []model.RouteStop{
		{BaseModel: model.BaseModel{ID: stopID}, StopOrder: 100},
		{BaseModel: model.BaseModel{ID: uuid.New()}, StopOrder: 200},
	}

	req := &model.MoveRouteStopRequest{Position: "last"}

	mockRouteStopRepo.EXPECT().GetByID(ctx, stopID).Return(stop, nil).Times(1)
	mockRouteRepo.EXPECT().GetRoutesWithRouteStops(ctx, routeID).Return(&model.Route{RouteStops: allStops}, nil).Times(1)

	// Last order is 200, so new order should be 201
	mockRouteStopRepo.EXPECT().ReorderStops(ctx, routeID, map[uuid.UUID]int{stopID: 201}).Return(nil).Times(1)

	result, err := service.MoveRouteStop(ctx, stopID, req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 201, result.StopOrder)
}

func TestMoveRouteStop_Error_MissingRef(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)

	service := NewRouteStopService(mockRouteStopRepo, mockRouteRepo)

	ctx := context.Background()
	stopID := uuid.New()
	routeID := uuid.New()

	stop := &model.RouteStop{BaseModel: model.BaseModel{ID: stopID}, RouteID: routeID}
	allStops := []model.RouteStop{
		{BaseModel: model.BaseModel{ID: stopID}},
		{BaseModel: model.BaseModel{ID: uuid.New()}}, // Need > 1 stop to trigger reorder logic
	}

	req := &model.MoveRouteStopRequest{Position: "before"} // Missing ReferenceStopID

	mockRouteStopRepo.EXPECT().GetByID(ctx, stopID).Return(stop, nil).Times(1)
	mockRouteRepo.EXPECT().GetRoutesWithRouteStops(ctx, routeID).Return(&model.Route{RouteStops: allStops}, nil).Times(1)

	_, err := service.MoveRouteStop(ctx, stopID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}
