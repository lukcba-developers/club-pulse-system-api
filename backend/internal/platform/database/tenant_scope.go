package database

import "gorm.io/gorm"

// TenantScope returns a GORM scope that filters queries by club_id.
// This is a security-critical helper to ensure multi-tenant data isolation.
//
// Usage:
//
//	db.Scopes(TenantScope(clubID)).Find(&entities)
//
// IMPORTANT: This scope should be applied to ALL queries that access tenant-specific data.
func TenantScope(clubID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("club_id = ?", clubID)
	}
}
