package utils

import (
	"errors"

	"gorm.io/gorm"
)

func CheckDuplicate[T any](db *gorm.DB, field string, value any, excludeID any) (bool, error) {
	var result T
	query := db.Where(field+" = ?", value)

	if excludeID != nil {
		query = query.Where("id != ?", excludeID)
	}

	err := query.First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // tidak duplikat
		}
		return false, err // error lain
	}

	return true, nil // ditemukan → duplikat
}
