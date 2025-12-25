package service

import (
	"context"
	"fmt"
	"testing"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewReviewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)

	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &ReviewServiceImpl{}, service)
}

func TestCreateReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()
	tripID := uuid.New()

	req := &model.CreateReviewRequest{
		BookingID: bookingID,
		Rating:    5,
		Comment:   "Great trip!",
	}

	booking := &model.Booking{
		BaseModel: model.BaseModel{
			ID: bookingID,
		},
		UserID: userID,
		TripID: tripID,
		Status: model.BookingStatusConfirmed,
	}

	// Expect GetBookingByID
	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	// Expect GetByBookingID for duplicate check
	mockReviewRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(nil, gorm.ErrRecordNotFound).
		Times(1)

	// Expect Create
	mockReviewRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, review *model.Review) error {
			assert.Equal(t, tripID, review.TripID)
			assert.Equal(t, userID, review.UserID)
			assert.Equal(t, bookingID, review.BookingID)
			assert.Equal(t, 5, review.Rating)
			assert.Equal(t, "Great trip!", review.Comment)
			assert.True(t, review.IsVerified)
			assert.Equal(t, model.ReviewStatusActive, review.Status)
			review.ID = uuid.New()
			return nil
		}).
		Times(1)

	result, err := service.CreateReview(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tripID, result.TripID)
	assert.Equal(t, 5, result.Rating)
}

func TestCreateReview_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.CreateReviewRequest{
		BookingID: bookingID,
		Rating:    5,
		Comment:   "Great trip!",
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, gorm.ErrRecordNotFound).
		Times(1)

	result, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "booking not found")
}

func TestCreateReview_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	differentUserID := uuid.New()
	bookingID := uuid.New()

	req := &model.CreateReviewRequest{
		BookingID: bookingID,
		Rating:    5,
	}

	booking := &model.Booking{
		BaseModel: model.BaseModel{
			ID: bookingID,
		},
		UserID: differentUserID, // Different user
		Status: model.BookingStatusConfirmed,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	result, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestCreateReview_InvalidBookingStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.CreateReviewRequest{
		BookingID: bookingID,
		Rating:    5,
	}

	booking := &model.Booking{
		BaseModel: model.BaseModel{
			ID: bookingID,
		},
		UserID: userID,
		Status: model.BookingStatusPending, // Not confirmed
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	result, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "can only review confirmed bookings")
}

func TestCreateReview_DuplicateReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.CreateReviewRequest{
		BookingID: bookingID,
		Rating:    5,
	}

	booking := &model.Booking{
		BaseModel: model.BaseModel{
			ID: bookingID,
		},
		UserID: userID,
		Status: model.BookingStatusConfirmed,
	}

	existingReview := &model.Review{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingID: bookingID,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockReviewRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(existingReview, nil).
		Times(1)

	result, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
}

func TestGetReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	reviewID := uuid.New()

	expectedReview := &model.Review{
		BaseModel: model.BaseModel{
			ID: reviewID,
		},
		Rating: 5,
	}

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(expectedReview, nil).
		Times(1)

	result, err := service.GetReview(ctx, reviewID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReview, result)
}

func TestGetReviewByBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	bookingID := uuid.New()

	expectedReview := &model.Review{
		BookingID: bookingID,
		Rating:    4,
	}

	mockReviewRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(expectedReview, nil).
		Times(1)

	result, err := service.GetReviewByBooking(ctx, bookingID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReview, result)
}

func TestGetUserReviews_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()

	req := &model.GetUserReviewsRequest{
		UserID: userID,
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	reviews := []*model.Review{
		{BaseModel: model.BaseModel{ID: uuid.New()}, UserID: userID, Rating: 5},
		{BaseModel: model.BaseModel{ID: uuid.New()}, UserID: userID, Rating: 4},
	}

	mockReviewRepo.EXPECT().
		GetUserReviews(ctx, userID, 1, 10).
		Return(reviews, int64(2), nil).
		Times(1)

	results, total, err := service.GetUserReviews(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
}

func TestGetTripReviews_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	tripID := uuid.New()

	req := &model.GetTripReviewsRequest{
		TripID: &tripID,
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	reviews := []*model.Review{
		{BaseModel: model.BaseModel{ID: uuid.New()}, TripID: tripID, Rating: 5},
		{BaseModel: model.BaseModel{ID: uuid.New()}, TripID: tripID, Rating: 4},
	}

	mockReviewRepo.EXPECT().
		GetTripReviews(ctx, req).
		Return(reviews, int64(2), nil).
		Times(1)

	results, total, err := service.GetTripReviews(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
}

func TestGetTripReviewSummary_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	tripID := uuid.New()

	expectedSummary := &model.TripReviewSummary{
		TripID:        tripID,
		TotalReviews:  10,
		AverageRating: 4.5,
	}

	mockReviewRepo.EXPECT().
		GetTripReviewSummary(ctx, tripID).
		Return(expectedSummary, nil).
		Times(1)

	result, err := service.GetTripReviewSummary(ctx, tripID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSummary, result)
}

func TestUpdateReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	reviewID := uuid.New()

	newRating := 4
	newComment := "Updated comment"
	req := &model.UpdateReviewRequest{
		Rating:  &newRating,
		Comment: &newComment,
	}

	existingReview := &model.Review{
		BaseModel: model.BaseModel{ID: reviewID},
		UserID:    userID,
		Rating:    5,
	}

	updatedReview := &model.Review{
		BaseModel: model.BaseModel{ID: reviewID},
		UserID:    userID,
		Rating:    4,
		Comment:   "Updated comment",
	}

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(existingReview, nil).
		Times(1)

	mockReviewRepo.EXPECT().
		Update(ctx, reviewID, gomock.Any()).
		Return(nil).
		Times(1)

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(updatedReview, nil).
		Times(1)

	result, err := service.UpdateReview(ctx, userID, reviewID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 4, result.Rating)
}

func TestUpdateReview_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	differentUserID := uuid.New()
	reviewID := uuid.New()

	newRating := 4
	req := &model.UpdateReviewRequest{
		Rating: &newRating,
	}

	existingReview := &model.Review{
		BaseModel: model.BaseModel{ID: reviewID},
		UserID:    differentUserID, // Different user
	}

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(existingReview, nil).
		Times(1)

	result, err := service.UpdateReview(ctx, userID, reviewID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestDeleteReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	reviewID := uuid.New()

	existingReview := &model.Review{
		BaseModel: model.BaseModel{ID: reviewID},
		UserID:    userID,
	}

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(existingReview, nil).
		Times(1)

	mockReviewRepo.EXPECT().
		Delete(ctx, reviewID).
		Return(nil).
		Times(1)

	err := service.DeleteReview(ctx, userID, reviewID)

	assert.NoError(t, err)
}

func TestDeleteReview_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	userID := uuid.New()
	differentUserID := uuid.New()
	reviewID := uuid.New()

	existingReview := &model.Review{
		BaseModel: model.BaseModel{ID: reviewID},
		UserID:    differentUserID,
	}

	mockReviewRepo.EXPECT().
		GetByID(ctx, reviewID).
		Return(existingReview, nil).
		Times(1)

	err := service.DeleteReview(ctx, userID, reviewID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestModerateReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	reviewID := uuid.New()
	status := model.ReviewStatusHidden
	adminNotes := "Inappropriate content"

	mockReviewRepo.EXPECT().
		Update(ctx, reviewID, map[string]interface{}{
			"status":      status,
			"admin_notes": adminNotes,
		}).
		Return(nil).
		Times(1)

	err := service.ModerateReview(ctx, reviewID, status, adminNotes)

	assert.NoError(t, err)
}

func TestModerateReview_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mocks.NewMockReviewRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	service := NewReviewService(mockReviewRepo, mockBookingRepo)

	ctx := context.Background()
	reviewID := uuid.New()
	status := model.ReviewStatusHidden
	adminNotes := "Test"

	expectedErr := fmt.Errorf("database error")

	mockReviewRepo.EXPECT().
		Update(ctx, reviewID, gomock.Any()).
		Return(expectedErr).
		Times(1)

	err := service.ModerateReview(ctx, reviewID, status, adminNotes)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to moderate review")
}
