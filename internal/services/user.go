package services

import (
	"context"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService interface {
	GetById(c context.Context, id string) (*types.User, error)
	GetByUsername(c context.Context, username string) (*types.User, error)
	GetUsers(c context.Context, params ...Param) ([]types.User, error)
	Save(c context.Context, user *models.User) error
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (us *userService) GetById(c context.Context, id string) (*types.User, error) {
	user := &models.User{}
	if err := us.db.Preload(clause.Associations).First(user, "Id = ?", id).Error; err != nil {
		return nil, err
	}

	return &us.mapUsers([]models.User{*user})[0], nil
}

func (us *userService) GetByUsername(c context.Context, username string) (*types.User, error) {
	user := &models.User{}
	if err := us.db.Preload(clause.Associations).First(user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &us.mapUsers([]models.User{*user})[0], nil
}

func (us *userService) GetUsers(c context.Context, params ...Param) ([]types.User, error) {
	users := []models.User{}
	if err := ApplyParams(us.db, params...).Preload(clause.Associations).Find(&users).Error; err != nil {
		return nil, err
	}

	return us.mapUsers(users), nil
}

func (us *userService) Save(c context.Context, user *models.User) error {
	return us.db.Save(user).Error
}

func (us *userService) mapUsers(users []models.User) []types.User {
	mus := []types.User{}
	for _, u := range users {
		mu := types.User{
			UserBase: types.UserBase{
				Id:       u.ID,
				Username: u.Username,
				Email:    u.Email,
			},
			Password:   u.Password,
			Subcribers: us.mapSubscribers(&u),
			Posts:      us.mapPosts(&u),
		}
		mus = append(mus, mu)
	}

	return mus
}

func (us *userService) mapSubscribers(u *models.User) []types.UserBase {
	mss := []types.UserBase{}
	for _, s := range u.Subcribers {
		ms := types.UserBase{
			Id:       s.ID,
			Username: s.Username,
			Email:    s.Email,
		}
		mss = append(mss, ms)
	}

	return mss
}

func (us *userService) mapPosts(u *models.User) []types.Post {
	mps := []types.Post{}
	for _, p := range u.Posts {
		mp := types.Post{
			Id:        p.ID,
			CreatedAt: p.CreatedAt,
			Title:     p.Title,
			Body:      p.Body,
			OwnerName: u.Username,
			OwnerId:   u.ID,
		}
		mps = append(mps, mp)
	}

	return mps
}
