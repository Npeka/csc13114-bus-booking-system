package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bus-booking/booking-service/internal/model"
)

type PaymentMethodRepository interface {
	GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*model.PaymentMethod, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (*model.PaymentMethod, error)
}

type paymentMethodRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepositoryImpl{db: db}
}

func (r *paymentMethodRepositoryImpl) GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethod, error) {
	var paymentMethods []*model.PaymentMethod
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name").
		Find(&paymentMethods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	return paymentMethods, nil
}

func (r *paymentMethodRepositoryImpl) GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	if err := r.db.WithContext(ctx).First(&paymentMethod, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	return &paymentMethod, nil
}

func (r *paymentMethodRepositoryImpl) GetPaymentMethodByCode(ctx context.Context, code string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	if err := r.db.WithContext(ctx).First(&paymentMethod, "code = ?", code).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	return &paymentMethod, nil
}
