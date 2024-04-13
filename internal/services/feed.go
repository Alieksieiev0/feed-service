package services

import (
	"context"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"gorm.io/gorm"
)

type FeedService interface {
	GetPosts(ctx context.Context, params ...Param) ([]models.Post, error)
	Subscribe(ctx context.Context, user *models.User, sub *models.User) error
	Post(ctx context.Context, user *models.User, post *models.Post) error
}

func NewFeedService(db *gorm.DB) FeedService {
	return &feedService{
		db: db,
	}
}

type feedService struct {
	db *gorm.DB
}

func (fs *feedService) GetPosts(ctx context.Context, params ...Param) ([]models.Post, error) {
	posts := []models.Post{}
	err := ApplyParams(fs.db, params...).Find(&posts).Error
	return posts, err
}

func (fs *feedService) Subscribe(ctx context.Context, user *models.User, sub *models.User) error {
	return fs.db.Model(user).Association("Subscribers").Append(sub)
}

func (fs *feedService) Post(ctx context.Context, user *models.User, post *models.Post) error {
	return fs.db.Model(user).Association("Posts").Append(post)
}
