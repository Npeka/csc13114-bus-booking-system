package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
	Transaction TransactionRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transaction: NewTransactionRepository(db),
	}
}
