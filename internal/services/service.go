package services

import (
	"context"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserFeedService interface {
	UserService
	FeedService
}

type UserService interface {
	Get(ctx context.Context, id string) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
}

type FeedService interface {
	Subscribe(ctx context.Context, user *models.User, sub *models.User) error
	Post(ctx context.Context, user *models.User, post *models.Post) error
}

func NewUserFeedService(db *gorm.DB) UserFeedService {
	return &userFeedService{
		userService: userService{db: db},
		feedService: feedService{db: db},
	}
}

type userFeedService struct {
	userService
	feedService
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (us *userService) Get(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	return user, us.db.Preload(clause.Associations).First(user, id).Error
}

func (us *userService) Save(ctx context.Context, user *models.User) error {
	return us.db.Save(user).Error
}

func NewFeedService(db *gorm.DB) FeedService {
	return &feedService{
		db: db,
	}
}

type feedService struct {
	db *gorm.DB
}

func (fs *feedService) Subscribe(ctx context.Context, user *models.User, sub *models.User) error {
	return fs.db.Model(user).Association("subscribers").Append(sub)
}

func (fs *feedService) Post(ctx context.Context, user *models.User, post *models.Post) error {
	return fs.db.Model(user).Association("posts").Append(post)
}
