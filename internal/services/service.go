package services

import (
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
