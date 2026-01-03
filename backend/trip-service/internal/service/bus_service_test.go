package service

import (
	"context"
	"mime/multipart"
	"testing"

	storage_mocks "bus-booking/shared/storage/mocks"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBusService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	assert.NotNil(t, service)
	assert.IsType(t, &BusServiceImpl{}, service)
}

func TestGetBusByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()

	expectedBus := &model.Bus{
		BaseModel:    model.BaseModel{ID: busID},
		PlateNumber:  "29A-12345",
		Model:        "Hyundai",
		SeatCapacity: 40,
	}

	mockBusRepo.EXPECT().
		GetBusWithSeatsByID(ctx, busID).
		Return(expectedBus, nil).
		Times(1)

	result, err := service.GetBusByID(ctx, busID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "29A-12345", result.PlateNumber)
}

func TestGetBusByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()

	mockBusRepo.EXPECT().
		GetBusWithSeatsByID(ctx, busID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetBusByID(ctx, busID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get bus")
}

func TestListBuses_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	req := model.ListBusesRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	expectedBuses := []model.Bus{
		{PlateNumber: "29A-11111"},
		{PlateNumber: "29A-22222"},
	}

	mockBusRepo.EXPECT().
		ListBuses(ctx, 1, 10).
		Return(expectedBuses, int64(2), nil).
		Times(1)

	buses, total, err := service.ListBuses(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, buses, 2)
	assert.Equal(t, int64(2), total)
}

func TestCreateBus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	req := &model.CreateBusRequest{
		PlateNumber: "29A-12345",
		Model:       "Hyundai",
		BusType:     "luxury",
		IsActive:    true,
		Floors: []model.FloorConfig{
			{
				Floor:   1,
				Rows:    5,
				Columns: 4,
				Seats: []model.SeatConfig{
					{Row: 1, Column: 1, SeatType: "standard"},
					{Row: 1, Column: 2, SeatType: "standard"},
				},
			},
		},
	}

	// No existing bus
	mockBusRepo.EXPECT().
		GetBusByPlateNumber(ctx, "29A-12345").
		Return(nil, assert.AnError).
		Times(1)

	mockBusRepo.EXPECT().
		CreateBus(ctx, gomock.Any()).
		Do(func(_ context.Context, bus *model.Bus) {
			assert.Equal(t, "29A-12345", bus.PlateNumber)
			assert.Equal(t, 2, bus.SeatCapacity)
			assert.Len(t, bus.Seats, 2)
		}).
		Return(nil).
		Times(1)

	result, err := service.CreateBus(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateBus_DuplicatePlateNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	req := &model.CreateBusRequest{
		PlateNumber: "29A-12345",
		Floors:      []model.FloorConfig{},
	}

	existingBus := &model.Bus{PlateNumber: "29A-12345"}

	mockBusRepo.EXPECT().
		GetBusByPlateNumber(ctx, "29A-12345").
		Return(existingBus, nil).
		Times(1)

	result, err := service.CreateBus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "plate number already exists")
}

func TestCreateBus_ExceedCapacity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()

	// Create 101 seats
	seats := make([]model.SeatConfig, 101)
	for i := 0; i < 101; i++ {
		seats[i] = model.SeatConfig{Row: 1, Column: i + 1, SeatType: "standard"}
	}

	req := &model.CreateBusRequest{
		PlateNumber: "29A-12345",
		Floors: []model.FloorConfig{
			{Floor: 1, Rows: 1, Columns: 101, Seats: seats},
		},
	}

	mockBusRepo.EXPECT().
		GetBusByPlateNumber(ctx, "29A-12345").
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.CreateBus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot exceed 100")
}

// Test the important helper function
func TestGenerateSeatsFromFloorConfig_SingleFloor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage).(*BusServiceImpl)

	floors := []model.FloorConfig{
		{
			Floor:   1,
			Rows:    2,
			Columns: 2,
			Seats: []model.SeatConfig{
				{Row: 1, Column: 1, SeatType: "standard"},
				{Row: 1, Column: 2, SeatType: "vip"},
				{Row: 2, Column: 1, SeatType: "standard"},
			},
		},
	}

	seats := service.generateSeatsFromFloorConfig(floors)

	assert.Len(t, seats, 3)
	assert.Equal(t, "A1", seats[0].SeatNumber) // Single floor no prefix
	assert.Equal(t, "A2", seats[1].SeatNumber)
	assert.Equal(t, "B1", seats[2].SeatNumber)
}

func TestGenerateSeatsFromFloorConfig_MultiFloor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage).(*BusServiceImpl)

	floors := []model.FloorConfig{
		{
			Floor:   1,
			Rows:    1,
			Columns: 2,
			Seats: []model.SeatConfig{
				{Row: 1, Column: 1, SeatType: "standard"},
			},
		},
		{
			Floor:   2,
			Rows:    1,
			Columns: 2,
			Seats: []model.SeatConfig{
				{Row: 1, Column: 1, SeatType: "vip"},
			},
		},
	}

	seats := service.generateSeatsFromFloorConfig(floors)

	assert.Len(t, seats, 2)
	assert.Equal(t, "F1-A1", seats[0].SeatNumber) // Multi floor has prefix
	assert.Equal(t, "F2-A1", seats[1].SeatNumber)
}

func TestGenerateSeatsFromFloorConfig_DuplicatePositions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage).(*BusServiceImpl)

	floors := []model.FloorConfig{
		{
			Floor:   1,
			Rows:    2,
			Columns: 2,
			Seats: []model.SeatConfig{
				{Row: 1, Column: 1, SeatType: "standard"},
				{Row: 1, Column: 1, SeatType: "vip"}, // Duplicate!
			},
		},
	}

	seats := service.generateSeatsFromFloorConfig(floors)

	assert.Len(t, seats, 1) // Duplicate skipped
}

func TestGenerateSeatsFromFloorConfig_CustomPriceMultiplier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage).(*BusServiceImpl)

	customMultiplier := 1.5
	floors := []model.FloorConfig{
		{
			Floor:   1,
			Rows:    1,
			Columns: 1,
			Seats: []model.SeatConfig{
				{Row: 1, Column: 1, SeatType: "standard", PriceMultiplier: &customMultiplier},
			},
		},
	}

	seats := service.generateSeatsFromFloorConfig(floors)

	assert.Len(t, seats, 1)
	assert.Equal(t, 1.5, seats[0].PriceMultiplier)
}

func TestUpdateBus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()

	existingBus := &model.Bus{
		BaseModel:   model.BaseModel{ID: busID},
		PlateNumber: "29A-12345",
		Model:       "Old Model",
	}

	newModel := "New Model"
	req := &model.UpdateBusRequest{
		Model: &newModel,
	}

	mockBusRepo.EXPECT().
		GetBusByID(ctx, busID).
		Return(existingBus, nil).
		Times(1)

	mockBusRepo.EXPECT().
		UpdateBus(ctx, gomock.Any()).
		Do(func(_ context.Context, bus *model.Bus) {
			assert.Equal(t, "New Model", bus.Model)
		}).
		Return(nil).
		Times(1)

	result, err := service.UpdateBus(ctx, busID, req)

	assert.NoError(t, err)
	assert.Equal(t, "New Model", result.Model)
}

func TestUpdateBus_DuplicatePlateNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()
	otherBusID := uuid.New()

	existingBus := &model.Bus{
		BaseModel:   model.BaseModel{ID: busID},
		PlateNumber: "29A-12345",
	}

	otherBus := &model.Bus{
		BaseModel:   model.BaseModel{ID: otherBusID},
		PlateNumber: "29A-99999",
	}

	newPlate := "29A-99999"
	req := &model.UpdateBusRequest{
		PlateNumber: &newPlate,
	}

	mockBusRepo.EXPECT().
		GetBusByID(ctx, busID).
		Return(existingBus, nil).
		Times(1)

	mockBusRepo.EXPECT().
		GetBusByPlateNumber(ctx, "29A-99999").
		Return(otherBus, nil).
		Times(1)

	result, err := service.UpdateBus(ctx, busID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "plate number already exists")
}

func TestDeleteBus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()

	mockBusRepo.EXPECT().
		DeleteBus(ctx, busID).
		Return(nil).
		Times(1)

	err := service.DeleteBus(ctx, busID)

	assert.NoError(t, err)
}

func TestDeleteBus_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorage)

	ctx := context.Background()
	busID := uuid.New()

	mockBusRepo.EXPECT().
		DeleteBus(ctx, busID).
		Return(assert.AnError).
		Times(1)

	err := service.DeleteBus(ctx, busID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete bus")
}

func TestUploadImages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorageService := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorageService)

	ctx := context.Background()
	busID := uuid.New()
	existingBus := &model.Bus{
		BaseModel: model.BaseModel{ID: busID},
		ImageURLs: []string{"http://example.com/img1.jpg"},
	}

	mockBusRepo.EXPECT().GetBusByID(ctx, busID).Return(existingBus, nil).Times(1)
	mockStorageService.EXPECT().UploadFile(ctx, gomock.Any(), gomock.Any(), "bus-images").Return("http://example.com/new.jpg", nil).Times(1)
	mockBusRepo.EXPECT().UpdateBus(ctx, gomock.Any()).Do(func(_ context.Context, b *model.Bus) {
		assert.Len(t, b.ImageURLs, 2)
		assert.Equal(t, "http://example.com/new.jpg", b.ImageURLs[1])
	}).Return(nil).Times(1)

	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     1024,
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", "image/jpeg")

	files := []multipart.File{nil}
	headers := []*multipart.FileHeader{fileHeader}

	result, err := service.UploadImages(ctx, busID, files, headers)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUploadImages_TooManyImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorageService := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorageService)

	ctx := context.Background()
	busID := uuid.New()

	existingBus := &model.Bus{
		BaseModel: model.BaseModel{ID: busID},
		ImageURLs: make([]string, 9),
	}

	mockBusRepo.EXPECT().GetBusByID(ctx, busID).Return(existingBus, nil).Times(1)

	files := make([]multipart.File, 2)
	headers := make([]*multipart.FileHeader, 2)

	_, err := service.UploadImages(ctx, busID, files, headers)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tối đa 10 ảnh")
}

func TestDeleteImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorageService := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorageService)

	ctx := context.Background()
	busID := uuid.New()
	targetURL := "http://example.com/img1.jpg"

	existingBus := &model.Bus{
		BaseModel: model.BaseModel{ID: busID},
		ImageURLs: []string{targetURL, "http://example.com/img2.jpg"},
	}

	mockBusRepo.EXPECT().GetBusByID(ctx, busID).Return(existingBus, nil).Times(1)
	mockStorageService.EXPECT().DeleteFile(ctx, targetURL).Return(nil).Times(1)
	mockBusRepo.EXPECT().UpdateBus(ctx, gomock.Any()).Do(func(_ context.Context, b *model.Bus) {
		assert.Len(t, b.ImageURLs, 1)
		assert.Equal(t, "http://example.com/img2.jpg", b.ImageURLs[0])
	}).Return(nil).Times(1)

	result, err := service.DeleteImage(ctx, busID, targetURL)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDeleteImage_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBusRepo := mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	mockStorageService := storage_mocks.NewMockStorageService(ctrl)

	service := NewBusService(mockBusRepo, mockSeatRepo, mockStorageService)

	ctx := context.Background()
	busID := uuid.New()

	existingBus := &model.Bus{
		BaseModel: model.BaseModel{ID: busID},
		ImageURLs: []string{"http://example.com/img1.jpg"},
	}

	mockBusRepo.EXPECT().GetBusByID(ctx, busID).Return(existingBus, nil).Times(1)

	_, err := service.DeleteImage(ctx, busID, "http://example.com/missing.jpg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Không tìm thấy ảnh")
}
