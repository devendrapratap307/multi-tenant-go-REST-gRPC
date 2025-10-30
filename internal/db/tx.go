package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionFunc signature
type TransactionFunc func(tx *gorm.DB) error

// WithTransaction runs fn inside a transaction with panic recovery
func WithTransaction(ctx context.Context, db *gorm.DB, fn TransactionFunc) (err error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback().Error
			err = fmt.Errorf("panic: %v", r)
			return
		}
		if err != nil {
			_ = tx.Rollback().Error
			return
		}
		err = tx.Commit().Error
	}()
	err = fn(tx)
	return
}
