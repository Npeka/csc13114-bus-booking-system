package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"bus-booking/shared/ginext"
	"bus-booking/shared/storage"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusService interface {
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	ListBuses(ctx context.Context, req model.ListBusesRequest) ([]model.Bus, int64, error)

	CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error)
	UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error)
	DeleteBus(ctx context.Context, id uuid.UUID) error

	UploadImages(ctx context.Context, busID uuid.UUID, files []multipart.File, headers []*multipart.FileHeader) (*model.Bus, error)
	DeleteImage(ctx context.Context, busID uuid.UUID, imageURL string) (*model.Bus, error)
}

type BusServiceImpl struct {
	busRepo        repository.BusRepository
	seatRepo       repository.SeatRepository
	storageService storage.StorageService
}

func NewBusService(
	busRepo repository.BusRepository,
	seatRepo repository.SeatRepository,
	storageService storage.StorageService,
) BusService {
	return &BusServiceImpl{
		busRepo:        busRepo,
		seatRepo:       seatRepo,
		storageService: storageService,
	}
}

func (s *BusServiceImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusWithSeatsByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get bus")
	}
	return bus, nil
}

func (s *BusServiceImpl) ListBuses(ctx context.Context, req model.ListBusesRequest) ([]model.Bus, int64, error) {
	buses, total, err := s.busRepo.ListBuses(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("failed to list buses")
	}

	return buses, total, nil
}

func (s *BusServiceImpl) CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error) {
	existing, err := s.busRepo.GetBusByPlateNumber(ctx, req.PlateNumber)
	if err == nil && existing != nil {
		return nil, ginext.NewBadRequestError("plate number already exists")
	}

	// Calculate total seat capacity from all floors
	totalCapacity := 0
	for _, floor := range req.Floors {
		totalCapacity += len(floor.Seats)
	}

	if totalCapacity > 100 {
		return nil, ginext.NewBadRequestError("total seat capacity cannot exceed 100")
	}

	bus := &model.Bus{
		PlateNumber:  req.PlateNumber,
		Model:        req.Model,
		BusType:      req.BusType,
		SeatCapacity: totalCapacity,
		Amenities:    req.Amenities,
		IsActive:     req.IsActive,
		Seats:        s.generateSeatsFromFloorConfig(req.Floors),
	}

	if err := s.busRepo.CreateBus(ctx, bus); err != nil {
		return nil, ginext.NewInternalServerError("failed to create bus")
	}

	return bus, nil
}

// generateSeatsFromFloorConfig creates seats based on individual seat configurations
func (s *BusServiceImpl) generateSeatsFromFloorConfig(floors []model.FloorConfig) []model.Seat {
	seats := make([]model.Seat, 0)
	rowNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}

	for _, floorConfig := range floors {
		// Create map to detect duplicate seat positions
		seatPositions := make(map[string]bool)

		for _, seatConfig := range floorConfig.Seats {
			// Check for duplicate positions
			posKey := fmt.Sprintf("%d-%d", seatConfig.Row, seatConfig.Column)
			if seatPositions[posKey] {
				// Skip duplicates (or could return error)
				continue
			}
			seatPositions[posKey] = true

			// Validate seat position is within floor boundaries
			if seatConfig.Row > floorConfig.Rows || seatConfig.Column > floorConfig.Columns {
				// Skip invalid positions (or could return error)
				continue
			}

			// Generate seat number
			seatNumber := fmt.Sprintf("%s%d", rowNames[seatConfig.Row-1], seatConfig.Column)
			if len(floors) > 1 {
				seatNumber = fmt.Sprintf("F%d-%s%d", floorConfig.Floor, rowNames[seatConfig.Row-1], seatConfig.Column)
			}

			seat := model.Seat{
				SeatNumber:  seatNumber,
				SeatType:    seatConfig.SeatType,
				Row:         seatConfig.Row,
				Column:      seatConfig.Column,
				Floor:       floorConfig.Floor,
				IsAvailable: true,
			}

			// Use custom price multiplier if provided, otherwise use default from seat type
			if seatConfig.PriceMultiplier != nil {
				seat.PriceMultiplier = *seatConfig.PriceMultiplier
			}

			seats = append(seats, seat)
		}
	}

	return seats
}

func (s *BusServiceImpl) UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get bus")
	}

	if req.PlateNumber != nil {
		existing, err := s.busRepo.GetBusByPlateNumber(ctx, *req.PlateNumber)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ginext.NewBadRequestError("plate number already exists")
		}
		bus.PlateNumber = *req.PlateNumber
	}

	if req.Model != nil {
		bus.Model = *req.Model
	}

	if req.BusType != nil {
		bus.BusType = *req.BusType
	}

	if req.Amenities != nil {
		bus.Amenities = *req.Amenities
	}

	if req.IsActive != nil {
		bus.IsActive = *req.IsActive
	}

	if err := s.busRepo.UpdateBus(ctx, bus); err != nil {
		return nil, ginext.NewInternalServerError("failed to update bus")
	}

	return bus, nil
}

func (s *BusServiceImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	if err := s.busRepo.DeleteBus(ctx, id); err != nil {
		return ginext.NewInternalServerError("failed to delete bus")
	}
	return nil
}

// UploadImages uploads multiple images for a bus
func (s *BusServiceImpl) UploadImages(ctx context.Context, busID uuid.UUID, files []multipart.File, headers []*multipart.FileHeader) (*model.Bus, error) {
	// Get existing bus
	bus, err := s.busRepo.GetBusByID(ctx, busID)
	if err != nil {
		return nil, err
	}

	// Validate total image count (max 10 images per bus)
	if len(bus.ImageURLs)+len(files) > 10 {
		return nil, ginext.NewBadRequestError("Mỗi xe chỉ được tối đa 10 ảnh")
	}

	uploadedURLs := make([]string, 0)

	// Upload each file
	for i, file := range files {
		header := headers[i]

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" && contentType != "image/webp" {
			// Clean up already uploaded files
			for _, url := range uploadedURLs {
				if err := s.storageService.DeleteFile(ctx, url); err != nil {
					log.Warn().Err(err).Str("url", url).Msg("Failed to delete cleanup file")
				}
			}
			return nil, ginext.NewBadRequestError("Chỉ chấp nhận file ảnh (JPEG, PNG, WebP)")
		}

		// Validate file size (max 5MB per image)
		if header.Size > 5*1024*1024 {
			// Clean up already uploaded files
			for _, url := range uploadedURLs {
				if err := s.storageService.DeleteFile(ctx, url); err != nil {
					log.Warn().Err(err).Str("url", url).Msg("Failed to delete cleanup file")
				}
			}
			return nil, ginext.NewBadRequestError("Kích thước mỗi ảnh không được vượt quá 5MB")
		}

		// Upload file
		imageURL, err := s.storageService.UploadFile(ctx, file, header, "bus-images")
		if err != nil {
			// Clean up already uploaded files
			for _, url := range uploadedURLs {
				if err := s.storageService.DeleteFile(ctx, url); err != nil {
					log.Warn().Err(err).Str("url", url).Msg("Failed to delete cleanup file")
				}
			}
			return nil, ginext.NewInternalServerError("Không thể tải ảnh lên")
		}

		uploadedURLs = append(uploadedURLs, imageURL)
	}

	// Append new URLs to existing ones
	newImageURLs := append(bus.ImageURLs, uploadedURLs...)
	bus.ImageURLs = newImageURLs

	// Update bus in database
	if err := s.busRepo.UpdateBus(ctx, bus); err != nil {
		// Clean up uploaded files
		for _, url := range uploadedURLs {
			if err := s.storageService.DeleteFile(ctx, url); err != nil {
				log.Warn().Err(err).Str("url", url).Msg("Failed to delete cleanup file")
			}
		}
		return nil, ginext.NewInternalServerError("Không thể cập nhật thông tin xe")
	}

	return bus, nil
}

// DeleteImage removes a specific image from a bus
func (s *BusServiceImpl) DeleteImage(ctx context.Context, busID uuid.UUID, imageURL string) (*model.Bus, error) {
	// Get existing bus
	bus, err := s.busRepo.GetBusByID(ctx, busID)
	if err != nil {
		return nil, err
	}

	// Find and remove the image URL
	found := false
	newImageURLs := make([]string, 0)
	for _, url := range bus.ImageURLs {
		if url == imageURL {
			found = true
			continue
		}
		newImageURLs = append(newImageURLs, url)
	}

	if !found {
		return nil, ginext.NewBadRequestError("Không tìm thấy ảnh")
	}

	// Delete from storage
	if err := s.storageService.DeleteFile(ctx, imageURL); err != nil {
		return nil, ginext.NewInternalServerError("Không thể xóa ảnh")
	}

	// Update bus's image URLs
	bus.ImageURLs = newImageURLs
	if err := s.busRepo.UpdateBus(ctx, bus); err != nil {
		return nil, ginext.NewInternalServerError("Không thể cập nhật thông tin xe")
	}

	return bus, nil
}
