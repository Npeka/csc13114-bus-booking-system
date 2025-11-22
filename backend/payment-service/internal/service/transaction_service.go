package service

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *model.CreateTransactionRequest) error
}

type TransactionServiceImpl struct {
	repositories *repository.Repositories
}

func NewTransactionService(repositories *repository.Repositories) TransactionService {
	return &TransactionServiceImpl{
		repositories: repositories,
	}
}

func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, transaction *model.CreateTransactionRequest) error {
	err := s.repositories.Transaction.CreateTransaction(ctx, &model.Transaction{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingID:     transaction.BookingID,
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		PaymentMethod: transaction.PaymentMethod,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transaction")
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}
