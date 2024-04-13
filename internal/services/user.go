package services

import (
	"context"
	"fmt"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService interface {
	GetById(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (us *userService) GetById(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	fmt.Println(id)
	return user, us.db.Preload(clause.Associations).First(user, "Id = ?", id).Error
}

func (us *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	return user, us.db.Preload(clause.Associations).First(user, "username = ?", username).Error
}

func (us *userService) Save(ctx context.Context, user *models.User) error {
	return us.db.Save(user).Error
}
