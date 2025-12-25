package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at"  gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-"           gorm:"index"`
}

type PaginationRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// Normalize ensures valid pagination values
func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}
