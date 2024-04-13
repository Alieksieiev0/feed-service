package services

import (
	"fmt"

	"gorm.io/gorm"
)

type UserFeedService interface {
	UserService
	FeedService
}

type userFeedService struct {
	userService
	feedService
}

func NewUserFeedService(db *gorm.DB) UserFeedService {
	return &userFeedService{
		userService: userService{db: db},
		feedService: feedService{db: db},
	}
}

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

func ApplyParams(db *gorm.DB, params ...Param) *gorm.DB {
	for _, param := range params {
		db = param(db)
	}
	return db
}
