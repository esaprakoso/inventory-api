package utils

import (
	"errors"

	"gorm.io/gorm"
)

func IsDuplicate[T any](db *gorm.DB, field string, value any, excludeID any) (bool, error) {
	var result T
	query := db.Where(field+" = ?", value)

	if excludeID != nil {
		query = query.Where("id != ?", excludeID)
	}

	err := query.First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // not a duplicate
		}
		return false, err // other error
	}

	return true, nil // found -> duplicate
}
