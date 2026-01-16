package database

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// WithTx returns a context with the GORM transaction attached.
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx retrieves the GORM transaction from the context, or returns nil if not present.
func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}
