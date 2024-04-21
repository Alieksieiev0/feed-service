package services

import (
	"fmt"

	"gorm.io/gorm"
)

type Param func(db *gorm.DB) *gorm.DB

func Limit(limit int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Order(column string, order string) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s  %s", column, order))
	}
}

func Filter(name string, value string, isStrict bool) Param {
	return func(db *gorm.DB) *gorm.DB {
		if isStrict {
			return db.Where(name+"= ?", value)
		}
		return db.Where(
			fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", name),
			fmt.Sprintf("%%%s%%", value),
		)
	}
}

func ApplyParams(db *gorm.DB, params ...Param) *gorm.DB {
	for _, param := range params {
		db = param(db)
	}
	return db
}
