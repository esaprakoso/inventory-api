package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var DB *gorm.DB

func RegisterCustomValidators(v *validator.Validate, db *gorm.DB) {
	DB = db

	v.RegisterValidation("exists", existsValidator)
}

func existsValidator(fl validator.FieldLevel) bool {
	param := fl.Param() // format: "table-column"
	parts := strings.Split(param, "-")
	if len(parts) != 2 {
		return false
	}

	tableName := parts[0]
	columnName := parts[1]
	value := fl.Field().Interface()

	var count int64
	err := DB.Table(tableName).Where(fmt.Sprintf("%s = ?", columnName), value).Count(&count).Error
	return err == nil && count > 0
}
