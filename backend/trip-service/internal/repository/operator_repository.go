package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OperatorRepository interface {
	CreateOperator(ctx context.Context, operator *model.Operator) error
	GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error)
	UpdateOperator(ctx context.Context, operator *model.Operator) error
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	ListOperators(ctx context.Context, page, limit int) ([]model.OperatorSummary, int64, error)
	GetOperatorByEmail(ctx context.Context, email string) (*model.Operator, error)
}

type OperatorRepositoryImpl struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) OperatorRepository {
	return &OperatorRepositoryImpl{db: db}
}

func (r *OperatorRepositoryImpl) CreateOperator(ctx context.Context, operator *model.Operator) error {
	return r.db.WithContext(ctx).Create(operator).Error
}

func (r *OperatorRepositoryImpl) GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error) {
	var operator model.Operator
	if err := r.db.WithContext(ctx).First(&operator, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &operator, nil
}

func (r *OperatorRepositoryImpl) UpdateOperator(ctx context.Context, operator *model.Operator) error {
	return r.db.WithContext(ctx).Model(operator).Updates(operator).Error
}

func (r *OperatorRepositoryImpl) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Operator{}, "id = ?", id).Error
}

func (r *OperatorRepositoryImpl) ListOperators(ctx context.Context, page, limit int) ([]model.OperatorSummary, int64, error) {
	var results []model.OperatorSummary
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.Operator{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Table("operators o").
		Select(`
			o.id, o.name, o.contact_email, o.contact_phone, o.status, o.created_at,
			COUNT(DISTINCT r.id) as active_routes,
			COUNT(DISTINCT b.id) as active_buses
		`).
		Joins("LEFT JOIN routes r ON o.id = r.operator_id AND r.is_active = true").
		Joins("LEFT JOIN buses b ON o.id = b.operator_id AND b.is_active = true").
		Group("o.id").
		Offset(offset).Limit(limit).
		Order("o.created_at DESC")

	rows, err := query.Rows()
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close rows")
		}
	}()

	for rows.Next() {
		var result model.OperatorSummary
		err := rows.Scan(
			&result.ID, &result.Name, &result.ContactEmail, &result.ContactPhone, &result.Status, &result.CreatedAt,
			&result.ActiveRoutes, &result.ActiveBuses,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, nil
}

// Get operator by email
func (r *OperatorRepositoryImpl) GetOperatorByEmail(ctx context.Context, email string) (*model.Operator, error) {
	var operator model.Operator
	if err := r.db.WithContext(ctx).Where("contact_email = ?", email).First(&operator).Error; err != nil {
		return nil, err
	}
	return &operator, nil
}
