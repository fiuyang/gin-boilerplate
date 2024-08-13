package engine

import "gorm.io/gorm"

var (
	Instance *gorm.DB // Instance engine
)

// Mock engine connection & query
func Mock(db *gorm.DB) {
	Instance = db
}
